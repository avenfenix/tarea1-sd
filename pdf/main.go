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

func (op *Operations) addFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("El archivo especificado no existe")
	}
	fmt.Println("Adding file:", filePath)
	url := fmt.Sprintf("https://%s/v1/upload", op.Server)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	file, _ := os.Open(filePath)
	io.Copy(part, file)
	file.Close()
	writer.Close()
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	if result["server_filename"] != nil {
		op.Files = append(op.Files, map[string]string{
			"server_filename": result["server_filename"].(string),
			"filename":        filePath,
		})
	}
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

	// Imprimir la solicitud HTTP antes de enviarla
	fmt.Println("Sending request to:", url)
	fmt.Println("Request body:", string(jsonData))
	fmt.Println("Authorization header:", req.Header)
	fmt.Println("Content-Type header:", req.Header.Get("Content-Type"))

	// Enviar la solicitud HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing task:", err)
		return
	}
	defer resp.Body.Close()

	// Imprimir la respuesta del servidor
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response headers:", resp.Header)
}

func (op *Operations) download(outputFilename string) {
	url := fmt.Sprintf("https://%s/v1/download/%s", op.Server, op.TaskID)
	fmt.Println("Download URL:", url) // Imprime la URL de descarga

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+op.Token)
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status) // Imprime cualquier mensaje de error en la respuesta
		return
	}

	out, _ := os.Create(outputFilename)
	defer out.Close()
	io.Copy(out, resp.Body)
	resp.Body.Close()

	fmt.Println("Download successful")
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
	fmt.Println("Files to be processed:", op.Files)
	op.execute(password)

	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + "_protected.pdf"
	op.download(fileName)
}
