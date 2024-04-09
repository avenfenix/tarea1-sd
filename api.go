package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// iLovePDF

type ResAuthPDF struct {
	Token string `json:"token"`
}

type ReqAuthPDF struct {
	Public_key string `json:"public_key"`
}

type ResStartPDF struct {
	Server string `json:"server"`
	Task   string `json:"task"`
}

type ResUploadPDF struct {
	Server_filename string `json:"server_filename"`
}

type ReqUploadPDF struct {
	Task       string `json:"task"`
	Cloud_file string `json:"cloud_file"`
}

type ReqProcessPDF struct {
	Task  string    `json:"task"`
	Tool  string    `json:"tool"`
	Files []FilePDF `json:"files"`
}

type FilePDF struct {
	Server_filename string `json:"server_filename"`
	Filename        string `json:"filename"`
}

type ResProcessPDF struct {
	Download_filename string `json:"download_filename"`
	Filesize          int    `json:"filesize"`
	Output_filesize   int    `json:"output_filesize"`
	Output_filenumber int    `json:"output_filenumber"`
	Output_extensions string `json:"output_extensions"`
	Timer             string `json:"timer"`
	Status            string `json:"status"`
}

func pdfAuth(public_key string) (string, error) {

	rd := &ReqAuthPDF{Public_key: public_key}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(rd)

	request, err := http.NewRequest("POST", "https://api.ilovepdf.com/v1/auth", data)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	re, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer re.Body.Close()

	var response ResAuthPDF
	json.NewDecoder(re.Body).Decode(&response)
	return response.Token, nil
}

func pdfStart(token string) (*ResStartPDF, error) {

	req, err := http.NewRequest("GET", "https://api.ilovepdf.com/v1/start/protect", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()

	var response ResStartPDF

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func pdfUpload(server string, task string, token string) (*ResUploadPDF, error) {

	rd := &ReqUploadPDF{Task: task, Cloud_file: ""}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(rd)
	url := fmt.Sprintf("https://%s/v1/upload", server, data)

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()

	var response ResUploadPDF

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func pdfProcess(server string, server_filename string, task string, token string) (*ResProcessPDF, error) {
	files := []FilePDF{FilePDF{
		Server_filename: server_filename,
		Filename:        "file.pdf",
	}}

	rd := &ReqProcessPDF{Task: task, Tool: "protect", Files: files}
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(rd)
	url := fmt.Sprintf("https://%s/v1/process", server)

	req, err := http.NewRequest("POST", url, data)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response ResProcessPDF

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func pdfDownload() {

}

// API

type LoginDataAPI struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResUserDataAPI struct {
}

type Message struct {
	Message string `json:"message"`
}

func login(c *gin.Context) {
	var atributos LoginDataAPI
	c.ShouldBind(&atributos)

}

func register(c *gin.Context) {

}

func apiClients(c *gin.Context) {

}

func apiProtect(c *gin.Context) {

}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}
}

func main() {
	fmt.Println("Tarea 1 Sistemas Distribuidos")

	// Variables de entorno
	loadEnv()

	/* // Get Token PDF
	token, err := pdfAuth(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		log.Fatal("Error al obtener autorizacion de la API")
	}
	// Iniciar herramienta protect
	task_data, err := pdfStart(token)
	if err != nil {
		log.Fatal("Error al inicializar la herramienta")
	}
	upload_data, err := pdfUpload(task_data.Server, task_data.Task, token)
	if err != nil {
		log.Fatal("Error al subir el archivo")
	}

	process_data, err := pdfProcess(task_data.Server, upload_data.Server_filename, task_data.Task, token)
	fmt.Println(process_data) */

	// Gin y Endpoints
	r := gin.Default()
	r.POST("/login", login)
	r.POST("/register", register)

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	conn_url := fmt.Sprintf("%s:%s", ip, port)
	r.Run(conn_url)
}
