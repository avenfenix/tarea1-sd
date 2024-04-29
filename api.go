package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	mongoclient *mongo.Client
}

// Usuario

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" `
	Last_name string             `json:"last_name" bson:"last_name"`
	Rut       string             `json:"rut" bson:"rut"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
}

func (app *App) LoginUsuario(c *gin.Context) {
	var credenciales User
	c.ShouldBind(&credenciales)

	coll := app.mongoclient.Database("tarea1").Collection("users")

	var doc User

	// Filtro para encontrar usuario.
	filter := bson.D{{Key: "email", Value: credenciales.Email}}

	// Buscamos en la coleccion si existe y lo recuperamos
	err := coll.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, map[string]interface{}{"message": "Usuario no encontrado"})
		}
		return
	}

	// Revisar credenciales
	if credenciales.Password == doc.Password {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"_id":       doc.ID.Hex(),
				"email":     doc.Email,
				"name":      doc.Name,
				"last_name": doc.Last_name,
				"rut":       doc.Rut,
			},
		}

		c.JSON(200, response)
		return
	}
	c.Status(500)

}

func (app *App) RegistrarUsuario(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Error al registrar cliente"}

	// Recuperamos payload usuario
	var nuevoUsuario User
	if err := c.ShouldBind(&nuevoUsuario); err != nil {
		return
	}

	// Filtro para revisar si usuario esta registrado.
	// No puede repetirse ni el rut ni el correo.
	filter := bson.D{{Key: "rut", Value: nuevoUsuario.Rut}, {Key: "email", Value: nuevoUsuario.Email}}

	coll := app.mongoclient.Database("tarea1").Collection("users")

	// Revisamos si ya esta registrado
	err := coll.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Si no esta registrado lo insertamos en la base de datos
			result, err := coll.InsertOne(context.TODO(), nuevoUsuario)
			if err != nil {
				c.JSON(400, response)
				return
			}
			oid, _ := result.InsertedID.(primitive.ObjectID)

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"_id":       oid.Hex(),
					"email":     nuevoUsuario.Email,
					"name":      nuevoUsuario.Name,
					"last_name": nuevoUsuario.Last_name,
					"rut":       nuevoUsuario.Rut,
				},
			}

			c.JSON(200, response)
			return
		}
	}

}

// Clientes

type Client struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Last_name string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Rut       string             `bson:"rut,omitempty" json:"rut,omitempty"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
}

func (app *App) RegistarCliente(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Error al registrar cliente"}

	// Recuperamos payload cliente
	var nuevoCliente Client
	if err := c.ShouldBind(&nuevoCliente); err != nil {
		c.JSON(400, response)
		return
	}

	// Buscamos en la base de datos si el cliente existe
	// No puede repetirse ni el rut ni el correo.
	filter := bson.D{{Key: "rut", Value: nuevoCliente.Rut}, {Key: "email", Value: nuevoCliente.Email}}
	coll := app.mongoclient.Database("tarea1").Collection("clients")
	err := coll.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Si no existe podemos registrar el cliente
			result, err := coll.InsertOne(context.TODO(), nuevoCliente)
			if err != nil {
				c.JSON(400, response)
				return
			}
			oid, _ := result.InsertedID.(primitive.ObjectID)

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"_id":       oid.Hex(),
					"email":     nuevoCliente.Email,
					"name":      nuevoCliente.Name,
					"last_name": nuevoCliente.Last_name,
					"rut":       nuevoCliente.Rut,
				},
			}

			c.JSON(200, response)
			return
		}
	}

	c.JSON(400, response)
}

type FormGetClient struct {
	Rut string `form:"rut"`
}

func (f FormGetClient) toFilter() bson.D {
	filter := bson.D{}
	// Aqui podemos personalizar los filtros
	if f.Rut != "" {
		filter = append(filter, bson.E{Key: "rut", Value: f.Rut})
	}
	return filter
}

func (app *App) getClients(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Error al obtener los clientes"}

	var form FormGetClient
	if err := c.ShouldBind(&form); err != nil {
	}

	// Obtenemos de la base de datos todos los clientes o lo que coincida con el filtro.
	// El filtro depende de lo parseado del url.
	coll := app.mongoclient.Database("tarea1").Collection("clients")
	cursor, err := coll.Find(context.TODO(), form.toFilter())
	if err != nil {
		c.JSON(400, response)
		return
	}

	// Aqui recuperamos los clientes
	var results []Client
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		c.JSON(400, response)
		return
	}

	// Preparamos la respuesta
	data := []map[string]interface{}{}
	for _, client := range results {
		cliente := map[string]interface{}{
			"_id":       client.ID.Hex(),
			"email":     client.Email,
			"name":      client.Name,
			"last_name": client.Last_name,
			"rut":       client.Rut,
		}
		data = append(data, cliente)
	}
	response = map[string]interface{}{"data": data}

	c.JSON(200, response)
	return
}

func (app *App) getClientByID(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Cliente no encontrado"}

	// Parseamos el id
	id := c.Param("id")

	// Convertimos a ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// El id proporcionado no es correcto
		response["message"] = "El ID proporcionado no es correcto."
		c.JSON(400, response)
		return
	}

	coll := app.mongoclient.Database("tarea1").Collection("clients")
	filter := bson.D{{Key: "_id", Value: oid}}
	var cliente Client
	err = coll.FindOne(context.TODO(), filter).Decode(&cliente)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Cliente no encontrado
			c.JSON(404, response)
		}
		return
	}

	// Preparamos la respuesta
	response = map[string]interface{}{
		"data": map[string]interface{}{
			"_id":       id,
			"email":     cliente.Email,
			"name":      cliente.Name,
			"last_name": cliente.Last_name,
			"rut":       cliente.Rut,
		},
	}

	c.JSON(200, response)
	return
}

func (app *App) putClientByID(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Cliente no encontrado"}

	// Parseamos y convertimos el ID
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)

	// Parseamos el payload con los cambios a realizar al cliente
	var cambios Client
	if err := c.ShouldBind(&cambios); err != nil {
		response["message"] = "Formato incorrecto"
		c.JSON(400, response)
	}

	coll := app.mongoclient.Database("tarea1").Collection("clients")

	// Verificamos si existe el cliente y lo obtenemos para luego mostrarlo sin los cambios
	var cliente Client
	filter := bson.D{{Key: "_id", Value: oid}}
	err = coll.FindOne(context.TODO(), filter).Decode(&cliente)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, response)
		}
		return
	}

	update := bson.M{"$set": cambios}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		response["message"] = "Error al actualizar el cliente"
	}
	var zero int64 = 0
	if result.MatchedCount == zero {
		c.JSON(404, response)
		return
	}
	if result.ModifiedCount == 1 || result.MatchedCount == 1 {

		response = map[string]interface{}{
			"data": map[string]interface{}{
				"_id":       id,
				"email":     cliente.Email,
				"name":      cliente.Name,
				"last_name": cliente.Last_name,
				"rut":       cliente.Rut,
			},
		}

		c.JSON(200, response)
		return
	}
	response["message"] = "Error al actualizar el cliente"
	c.JSON(400, response)
}

func (app *App) delClientByID(c *gin.Context) {
	// Respuesta predeterminada
	response := map[string]interface{}{"message": "Cliente no encontrado"}

	// Parseamos el id y lo convertimos en ObjectID
	id := c.Param("id")
	oid, _ := primitive.ObjectIDFromHex(id)

	coll := app.mongoclient.Database("tarea1").Collection("clients")

	// Buscamos en la db al cliente y lo guardamos para luego mostrarlo
	filter := bson.D{{Key: "_id", Value: oid}}
	var cliente Client
	err := coll.FindOne(context.TODO(), filter).Decode(&cliente)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, response)
			return
		}
	}
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		response["message"] = "Error al eliminar el cliente"
		c.JSON(400, response)
		return
	}
	if result.DeletedCount == 1 {
		response = map[string]interface{}{
			"data": map[string]interface{}{
				"_id":       id,
				"email":     cliente.Email,
				"name":      cliente.Name,
				"last_name": cliente.Last_name,
				"rut":       cliente.Rut,
			},
		}

		c.JSON(200, response)
		return
	}

}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	// Conexion con MongoDB
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
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
	r.Static("/uploads", "./uploads")
	r.MaxMultipartMemory = 8 << 20
	//r.POST("/api/protect", app.Protect)

	r.POST("/login", app.LoginUsuario)
	r.POST("/register", app.RegistrarUsuario)
	r.POST("/api/clients", app.RegistarCliente)
	r.GET("/api/clients", app.getClients)
	r.GET("/api/clients/:id", app.getClientByID)
	r.PUT("/api/clients/:id", app.putClientByID)
	r.DELETE("/api/clients/:id", app.delClientByID)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
