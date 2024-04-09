package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type App struct {
	mongoclient *mongo.Client
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Message string `json:"message"`
}

type UserDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `json:"name"`
	Last_name string             `json:"last_name"`
	Rut       string             `json:"rut"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
}

type UserData struct {
	ID        primitive.ObjectID `json:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `json:"name"`
	Last_name string             `json:"last_name"`
	Rut       string             `json:"rut"`
	Email     string             `json:"email"`
}

type ResUserData struct {
	Data UserData `json:"data"`
}

func (app *App) login(c *gin.Context) {
	var atributos LoginData
	c.ShouldBind(&atributos)

	// Esta registrado?

	coll := app.mongoclient.Database("tarea1").Collection("users")

	var doc UserDocument
	filter := bson.D{{Key: "email", Value: atributos.Email}}
	err := coll.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var msg Message
			msg.Message = "Usuario no encontrado"
			c.JSON(404, msg)
			return
		}
	}
	if atributos.Password == doc.Password {
		var result ResUserData
		result.Data.ID = doc.ID
		result.Data.Email = doc.Email
		result.Data.Last_name = doc.Last_name
		result.Data.Name = doc.Name
		result.Data.Rut = doc.Rut

		c.JSON(200, result)
		return
	}

}

func (app *App) register(c *gin.Context) {

	// Parsing params
	var newclient UserDocument
	if err := c.ShouldBind(&newclient); err != nil {
		return
	}
	coll := app.mongoclient.Database("tarea1").Collection("users")
	filter := bson.D{{Key: "rut", Value: newclient.Rut}}
	var doc UserDocument
	err := coll.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			/* Si no existe documento es porque no esta registrado,
			   entonces insertamos el nuevo cliente
			*/

			result, err := coll.InsertOne(context.TODO(), newclient)
			if err != nil {
				log.Fatal(err)
			}
			oid, _ := result.InsertedID.(primitive.ObjectID)

			var response ResUserData
			response.Data.ID = oid
			response.Data.Email = newclient.Email
			response.Data.Last_name = newclient.Last_name
			response.Data.Name = newclient.Name
			response.Data.Rut = newclient.Rut

			c.JSON(200, response)
			return
		}
	}

}

func (app *App) clients(c *gin.Context) {

}

func Protect(c *gin.Context) {

}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	// MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	mongoclient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongoclient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := App{
		mongoclient: mongoclient,
	}

	// Gin
	r := gin.Default()
	r.POST("/login", app.login)
	r.POST("/register", app.register)
	r.POST("/api/clients", app.clients)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
