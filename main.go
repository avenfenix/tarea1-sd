package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Name      string `bson:"name" json:"name"`
	Last_name string `bson:"last_name" json:"last_name"`
	Rut       string `bson:"rut" json:"rut"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
}

type RegisterClient struct {
	Name      string `bson:"name,omitempty" json:"name" form:"name"`
	Last_name string `bson:"last_name,omitempty" json:"last_name" form:"last_name"`
	Rut       string `bson:"rut,omitempty" json:"rut" form:"rut"`
	Email     string `bson:"email,omitempty" json:"email" form:"email"`
}

type PersonalData struct {
	ID        primitive.ObjectID `json:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `json:"name"`
	Last_name string             `json:"last_name"`
	Rut       string             `json:"rut"`
	Email     string             `json:"email"`
}

type ResPersonalData struct {
	Data PersonalData `json:"data"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	fmt.Println("\nBienvenido al sistema de protección de archivos de DiSis.")
	fmt.Println("Para utilizar la aplicación seleccione los números correspondientes al menú.")

	menu := true
	menu_usuario := false
	menu_cliente := false
	var line string

	for menu {
		fmt.Println("\nIngrese o regístrese")
		fmt.Println("  1) Ingreso")
		fmt.Println("  2) Registro")
		fmt.Println("  3) Salir")
		fmt.Print("\n> ")

		fmt.Scanln(&line)
		switch line {
		case "1":
			{
				// Ingreso
				var credentials LoginUser
				fmt.Print("Ingrese el correo de su cuenta: ")
				fmt.Scanln(&credentials.Email)
				fmt.Print("Ingrese su contraseña: ")
				fmt.Scanln(&credentials.Password)

				url := fmt.Sprintf("http://%s:%s/login", os.Getenv("HOST"), os.Getenv("PORT"))

				data := new(bytes.Buffer)
				json.NewEncoder(data).Encode(credentials)

				req, _ := http.NewRequest("POST", url, data)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				res, _ := client.Do(req)
				defer res.Body.Close()

				if res.Header.Get("Content-Type") == "application/json; charset=utf-8" {
					fmt.Println("¡Login exitoso!")
					menu = false
					menu_usuario = true
				}

			}
		case "2":
			{
				// Registro
				var newUser RegisterUser

				fmt.Print("Ingrese su nombre: ")
				fmt.Scanln(&newUser.Name)
				fmt.Print("Ingrese su apellido: ")
				fmt.Scanln(&newUser.Last_name)
				fmt.Print("Ingrese su RUT: ")
				fmt.Scanln(&newUser.Rut)
				fmt.Print("Ingrese su correo: ")
				fmt.Scanln(&newUser.Email)
				fmt.Print("Ingrese su contraseña: ")
				fmt.Scanln(&newUser.Password)

				data := new(bytes.Buffer)
				json.NewEncoder(data).Encode(newUser)
				url := fmt.Sprintf("http://%s:%s/register", os.Getenv("HOST"), os.Getenv("PORT"))
				req, err := http.NewRequest("POST", url, data)
				req.Header.Set("Content-Type", "application/json")

				if err != nil {

				}
				client := &http.Client{}
				re, err := client.Do(req)
				if err != nil {

				}
				defer re.Body.Close()

				if re.Header.Get("Content-Type") == "application/json; charset=utf-8" {
					fmt.Println("¡Registro exitoso!")
				}
			}
		case "3":
			{
				menu = false
			}
		}
	}
	for menu_usuario {
		fmt.Println("\nMenu principal")
		fmt.Println("  1) Clientes")
		fmt.Println("  2) Protección")
		fmt.Println("  3) Salir")
		fmt.Print("\n> ")

		fmt.Scanln(&line)
		switch line {
		case "1":
			{
				menu_cliente = true
				for menu_cliente {
					fmt.Println("\nMenu clientes")
					fmt.Println("  1) Listar los clientes registrados")
					fmt.Println("  2) Obtener un cliente por ID")
					fmt.Println("  3) Obtener un cliente por RUT")
					fmt.Println("  4) Registrar un nuevo cliente")
					fmt.Println("  5) Actualizar datos de un cliente")
					fmt.Println("  6) Borrar un cliente por ID")
					fmt.Println("  7) Volver")
					fmt.Print("\n> ")

					fmt.Scanln(&line)

					switch line {
					case "1":
						{
							// Listar los clientes registrados
							url := fmt.Sprintf("http://%s:%s/api/clients/", os.Getenv("HOST"), os.Getenv("PORT"))
							req, _ := http.NewRequest("GET", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							fmt.Print(res)
						}
					case "2":
						{
							// Obtener un cliente por ID
							id := ""
							url := fmt.Sprintf("http://%s:%s/api/clients/%s", os.Getenv("HOST"), os.Getenv("PORT"), id)
							req, _ := http.NewRequest("GET", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							fmt.Print(res)
						}
					case "3":
						{
							// Obtener un cliente por RUT
							url := fmt.Sprintf("http://%s:%s/api/clients", os.Getenv("HOST"), os.Getenv("PORT"))
							req, _ := http.NewRequest("GET", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							fmt.Print(res)
						}
					case "4":
						{
							// Registrar un nuevo cliente
							var newclient RegisterClient
							fmt.Print("Ingrese nombre del cliente: ")
							fmt.Scanln(&newclient.Name)
							fmt.Print("Ingrese apellido: ")
							fmt.Scanln(&newclient.Last_name)
							fmt.Print("Ingrese RUT del cliente: ")
							fmt.Scanln(&newclient.Rut)
							fmt.Print("Ingrese el correo del cliente: ")
							fmt.Scanln(&newclient.Email)

							data := new(bytes.Buffer)
							json.NewEncoder(data).Encode(newclient)

							url := fmt.Sprintf("http://%s:%s/api/clients", os.Getenv("HOST"), os.Getenv("PORT"))

							req, _ := http.NewRequest("POST", url, data)
							req.Header.Set("Content-Type", "application/json")

							client := &http.Client{}
							re, _ := client.Do(req)
							if re.Header.Get("Content-Type") == "application/json; charset=utf-8" {
								fmt.Printf("¡Cliente “%s” creado con éxito\n", newclient.Name)
							}
						}
					case "5":
						{
							// Actualizar datos de un cliente
							url := fmt.Sprintf("http://%s:%s/api/clients", os.Getenv("HOST"), os.Getenv("PORT"))
							req, _ := http.NewRequest("PUT", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							fmt.Print(res)
						}
					case "6":
						{
							// Borrar un cliente por ID
							url := fmt.Sprintf("http://%s:%s/api/clients", os.Getenv("HOST"), os.Getenv("PORT"))
							req, _ := http.NewRequest("GET", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							fmt.Print(res)
						}
					case "7":
						{
							menu_cliente = false
						}
					}
					line = ""
				}

			}

		case "2":
			{
				// Proteccion

			}

		case "3":
			{
				fmt.Println("¡Vuelve pronto!")
				return
			}
		}

	}
}
