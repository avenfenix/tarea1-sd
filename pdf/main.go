package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	APIEntryPoint1 = "https://api.ilovepdf.com/v1/auth"
	APIEntryPoint2 = "https://api.ilovepdf.com/v1/start"
)

type ILovePdf struct {
	PublicKey string
	Token     string
}

func NewILovePdf(publicKey string) *ILovePdf {
	resp, _ := http.PostForm(APIEntryPoint1, map[string][]string{
		"public_key": {publicKey},
	})
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &ILovePdf{PublicKey: publicKey, Token: result["token"]}
}

type Operations struct {
	*ILovePdf
	TaskID string
	Tool   string
	Server string
	Files  []map[string]string
}

func NewOperations(publicKey string) *Operations {
	return &Operations{ILovePdf: NewILovePdf(publicKey)}
}

func (op *Operations) startTask(tool string) {
	op.Tool = tool
	req, _ := http.NewRequest("GET", APIEntryPoint2+"/"+tool, nil)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	resp, _ := http.DefaultClient.Do(req)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	op.TaskID, op.Server = result["task"].(string), result["server"].(string)

}

func (op *Operations) addFile(filename string) error {
	// Verificar si el archivo existe
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("el archivo especificado no existe: %s", filename)
	}

	// Construir la URL de la solicitud
	url := fmt.Sprintf("https://%s/v1/upload", op.Server)

	// Abrir el archivo
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Crear un buffer para el cuerpo del formulario
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Agregar el archivo al formulario
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error al crear la parte del formulario: %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("error al copiar el contenido del archivo: %v", err)
	}

	// Agregar el parámetro "task" al formulario
	writer.WriteField("task", op.TaskID)

	// Cerrar el escritor multipart
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error al cerrar el escritor multipart: %v", err)
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error al crear la solicitud http: %v", err)
	}

	// Establecer el tipo de contenido en la solicitud
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Agregar el encabezado de autorización
	req.Header.Set("Authorization", "Bearer "+op.Token)

	// Realizar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar la solicitud http: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta JSON
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error al decodificar la respuesta json: %v", err)
	}

	// Verificar si el archivo se agregó correctamente
	if serverFilename, ok := response["server_filename"].(string); ok {
		op.Files = append(op.Files, map[string]string{
			"server_filename": serverFilename,
			"filename":        filename,
		})
		return nil
	}

	return fmt.Errorf("error al agregar el archivo: %v", response)
}

func (op *Operations) execute(password string) {
	url := fmt.Sprintf("https://%s/v1/process", op.Server)
	params := map[string]interface{}{
		"task":     op.TaskID,
		"tool":     op.Tool,
		"files":    op.Files,
		"password": password,
	}
	jsonData, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+op.Token)
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing task:", err)
		return
	}
	defer resp.Body.Close()

}

func (op *Operations) download(outputFilename string, inputPath string) {
	url := fmt.Sprintf("https://%s/v1/download/%s", op.Server, op.TaskID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status) // Imprime cualquier mensaje de error en la respuesta
		return
	}

	// Obtener el directorio del archivo de entrada
	outputDir := filepath.Dir(inputPath)

	// Concatenar el directorio y el nombre de archivo de salida
	outputPath := filepath.Join(outputDir, outputFilename)

	out, _ := os.Create(outputPath)
	defer out.Close()
	io.Copy(out, resp.Body)
	resp.Body.Close()

	fmt.Println("Descargado en:", outputPath)
}

func main() {
	publicKey := "project_public_db7deec963dc9219b319768d2766bfc6_9-1mScb0a712112737d004c62656bb16f2eb1"
	op := NewOperations(publicKey)
	op.startTask("protect")

	var password, path string

	fmt.Print("Escriba el ID del cliente objetivo: ")
	fmt.Scanln(&password)

	fmt.Print("Escriba la ruta donde se encuentra el archivo (incluya el nombre): ")
	fmt.Scanln(&path)

	op.addFile(path)
	op.execute(password)

	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + "_protegido.pdf"
	op.download(fileName, path)
}
