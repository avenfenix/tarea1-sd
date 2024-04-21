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

type Message struct {
	Message string `json:"message"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `json:"name"`
	Last_name string             `json:"last_name"`
	Rut       string             `json:"rut"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Name      string `json:"name"`
	Last_name string `json:"last_name"`
	Rut       string `json:"rut"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Client struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Last_name string             `bson:"last_name" json:"last_name"`
	Rut       string             `bson:"rut" json:"rut"`
	Email     string             `bson:"email" json:"email"`
}

type RegisterClient struct {
	Name      string `bson:"name,omitempty" json:"name" form:"name"`
	Last_name string `bson:"last_name,omitempty" json:"last_name" form:"last_name"`
	Rut       string `bson:"rut,omitempty" json:"rut" form:"rut"`
	Email     string `bson:"email,omitempty" json:"email" form:"email"`
}

type PersonalData struct {
	ID        string 			 `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Last_name string             `json:"last_name" bson:"last_name"`
	Rut       string             `json:"rut" bson:"rut`
	Email     string             `json:"email" bson:"email"`
}

type ResPersonalData struct {
	Data PersonalData `json:"data"`
}

type ResArrayData struct {
	Data []PersonalData `json:"data"`
}

type FormGetClient struct {
	Rut string `form:"rut"`
}

func (f FormGetClient) toBsonD() bson.D {
	filter := bson.D{}
	if f.Rut != ""{
		filter = append(filter, bson.E{Key:"rut", Value: f.Rut})
	}
	return filter
}


// Estructuras de prueba

type ResClients struct {
	Data []Client `json:"data" bson:"data"`
}

type IDParam struct {
	ID primitive.ObjectID `uri:"id" binding:"required,uuid"`
}

/////////////////////////


func TransformAll(clients []Client) ResArrayData{
	var response ResArrayData
	for _, client := range clients{
		var res PersonalData
		res.ID = client.ID.Hex()
		res.Name = client.Name
		res.Last_name = client.Last_name
		res.Rut = client.Rut
		res.Email = client.Email
		response.Data = append(response.Data, res)
	}
	return response
}

func (app *App) login(c *gin.Context) {
	var atributos LoginUser
	c.ShouldBind(&atributos)

	coll := app.mongoclient.Database("tarea1").Collection("users")

	var doc User
	filter := bson.D{{Key: "email", Value: atributos.Email}}
	err := coll.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, Message{Message: "Usuario no encontrado"})
		}
		return
	}
	if atributos.Password == doc.Password {

		response := ResPersonalData{Data: PersonalData{
			ID:        doc.ID.Hex(),
			Email:     doc.Email,
			Name:      doc.Name,
			Last_name: doc.Last_name,
			Rut:       doc.Rut,
		}}

		c.JSON(200, response)
		return
	}
	c.Status(500)

}

func (app *App) register(c *gin.Context) {

	var newuser RegisterUser
	if err := c.ShouldBind(&newuser); err != nil {
		return
	}
	coll := app.mongoclient.Database("tarea1").Collection("users")
	filter := bson.D{{Key: "rut", Value: newuser.Rut}, {Key: "email", Value: newuser.Email}}
	var user User
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			result, err := coll.InsertOne(context.TODO(), newuser)
			if err != nil {
				return
			}
			oid, _ := result.InsertedID.(primitive.ObjectID)

			response := ResPersonalData{Data: PersonalData{
				ID:        oid.Hex(),
				Email:     newuser.Email,
				Name:      newuser.Name,
				Last_name: newuser.Last_name,
				Rut:       newuser.Rut,
			}}

			c.JSON(200, response)
			return
		}
	}

}

func (app *App) register_client(c *gin.Context) {
	var newclient RegisterClient
	if err := c.ShouldBind(&newclient); err != nil {
		return
	}

	coll := app.mongoclient.Database("tarea1").Collection("clients")
	filter := bson.D{{Key: "rut", Value: newclient.Rut}, {Key: "email", Value: newclient.Email}}
	var client Client
	err := coll.FindOne(context.TODO(), filter).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			result, err := coll.InsertOne(context.TODO(), newclient)
			if err != nil {
				return
			}
			oid, _ := result.InsertedID.(primitive.ObjectID)

			response := ResPersonalData{Data: PersonalData{
				ID:        oid.Hex(),
				Email:     newclient.Email,
				Name:      newclient.Name,
				Last_name: newclient.Last_name,
				Rut:       newclient.Rut,
			}}

			c.JSON(200, response)
			return
		}
	}
}

func (app *App) get_clients(c *gin.Context) {

	var form FormGetClient
	if err:= c.ShouldBind(&form); err != nil{}

	coll := app.mongoclient.Database("tarea1").Collection("clients")

	cursor, err := coll.Find(context.TODO(), form.toBsonD() )
	if err != nil{
		c.Status(500)
		return
	}

	var results []Client
	
	err = cursor.All(context.TODO(), &results)
	if err != nil{
		c.Status(500 )
		return
	}

	response := TransformAll(results)


	c.JSON(200, response)
	return
}

func (app *App) get_clients_with_id(c *gin.Context) {

	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	coll := app.mongoclient.Database("tarea1").Collection("clients")
	filter := bson.D{{Key: "_id", Value: oid}}
	var client Client
	err = coll.FindOne(context.TODO(), filter).Decode(&client)
	if err != nil {
		return
	}

	response := ResPersonalData{Data: PersonalData{
		ID:        oid.Hex(),
		Email:     client.Email,
		Name:      client.Name,
		Last_name: client.Last_name,
		Rut:       client.Rut,
	}}
	c.JSON(200, response)
	return
}

func (app *App) put_clients_with_id(c *gin.Context) {

	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)

	var changes RegisterClient
	if err := c.ShouldBind(&changes); err != nil {
		return
	}

	coll := app.mongoclient.Database("tarea1").Collection("clients")
	filter := bson.D{{Key: "_id", Value: oid}}

	var old_client Client
	err = coll.FindOne(context.TODO(), filter).Decode(&old_client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, Message{Message: "Usuario no encontrado"})
			return
		}
	}

	bsonBytes, err := bson.Marshal(changes)
	if err != nil {

	}

	var bsonDoc bson.D
	err = bson.Unmarshal(bsonBytes, &bsonDoc)
	if err != nil {

	}

	update := bson.D{{Key: "$set", Value: bsonDoc}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {

	}
	var zero int64 = 0
	if result.MatchedCount == zero {
		c.JSON(404, Message{Message: "Usuario no encontrado"})
		return
	}
	if result.ModifiedCount == 1 || result.MatchedCount == 1 {
		response := ResPersonalData{Data: PersonalData{
			ID:        oid.Hex(),
			Email:     old_client.Email,
			Name:      old_client.Name,
			Last_name: old_client.Last_name,
			Rut:       old_client.Rut,
		}}

		c.JSON(200, response)
		return
	}

}

func (app *App) del_clients_with_id(c *gin.Context) {
	id := c.Param("id")
	oid, _ := primitive.ObjectIDFromHex(id)

	coll := app.mongoclient.Database("tarea1").Collection("clients")
	filter := bson.D{{Key: "_id", Value: oid}}
	var client Client
	err := coll.FindOne(context.TODO(), filter).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, Message{Message: "Usuario no encontrado"})
			return
		}
	}
	result, _ := coll.DeleteOne(context.TODO(), filter)
	if result.DeletedCount == 1 {
		response := ResPersonalData{Data: PersonalData{
			ID:        oid.Hex(),
			Email:     client.Email,
			Name:      client.Name,
			Last_name: client.Last_name,
			Rut:       client.Rut,
		}}
		c.JSON(200, response)
	}

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
	r.POST("/api/clients", app.register_client)
	r.GET("/api/clients", app.get_clients)
	r.GET("/api/clients/:id", app.get_clients_with_id)
	r.PUT("/api/clients/:id", app.put_clients_with_id)
	r.DELETE("/api/clients/:id", app.del_clients_with_id)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
