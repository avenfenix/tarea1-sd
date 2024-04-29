package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Last_name string `json:"last_name" bson:"last_name"`
	Rut       string `json:"rut" bson:"rut"`
	Email     string `json:"email" bson:"email"`
}

type ResPersonalData struct {
	Data PersonalData `json:"data"`
}

type ResArrayData struct {
	Data []PersonalData `json:"data"`
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

				if res.StatusCode == 200 {
					fmt.Println("¡Login exitoso!")
					menu = false
					menu_usuario = true
				} else {
					fmt.Println("Ha ocurrido un error al iniciar sesion!")
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

				if re.StatusCode == 200 {
					fmt.Println("¡Registro exitoso!")
				} else {
					fmt.Println("Ha ocurrido un error al registrarse!")
				}
			}
		case "3":
			{
				menu = false
				fmt.Println("¡Vuelve pronto!")
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
							defer res.Body.Close()

							var datos ResArrayData
							json.NewDecoder(res.Body).Decode(&datos)

							if res.StatusCode == 200 {
								for _, d := range datos.Data {
									fmt.Println("---")
									fmt.Printf("ID: %s\n", d.ID)
									fmt.Printf("Nombre: %s\n", d.Name)
									fmt.Printf("Apellido: %s\n", d.Last_name)
									fmt.Printf("RUT: %s\n", d.Rut)
									fmt.Printf("Email: %s\n", d.Email)

								}
								fmt.Println("---")
							} else {
								fmt.Println("Error al obtener cliente/s!")
							}
						}
					case "2":
						{
							// Obtener un cliente por ID
							var id string
							fmt.Print("Ingrese el ID a buscar: ")
							fmt.Scanln(&id)
							fmt.Print("\n")
							url := fmt.Sprintf("http://%s:%s/api/clients/%s", os.Getenv("HOST"), os.Getenv("PORT"), id)
							req, _ := http.NewRequest("GET", url, nil)
							http_client := &http.Client{}
							res, _ := http_client.Do(req)
							var datos ResPersonalData
							json.NewDecoder(res.Body).Decode(&datos)
							defer res.Body.Close()
							if res.StatusCode == 200 {
								fmt.Println("---")
								fmt.Printf("ID: %s\n", datos.Data.ID)
								fmt.Printf("Nombre: %s\n", datos.Data.Name)
								fmt.Printf("Apellido: %s\n", datos.Data.Last_name)
								fmt.Printf("RUT: %s\n", datos.Data.Rut)
								fmt.Printf("Email: %s\n", datos.Data.Email)
								fmt.Println("---")
							} else {
								fmt.Println("Error al obtener el cliente!")
							}
						}

					case "3":
						{
							// Obtener un cliente por RUT
							var rut string
							fmt.Print("Ingrese el RUT a buscar: ")
							fmt.Scanln(&rut)
							fmt.Print("\n")
							url := fmt.Sprintf("http://%s:%s/api/clients?rut=%s", os.Getenv("HOST"), os.Getenv("PORT"), rut)
							req, _ := http.NewRequest("GET", url, nil)
							client := &http.Client{}
							res, _ := client.Do(req)
							defer res.Body.Close()

							var datos ResArrayData
							json.NewDecoder(res.Body).Decode(&datos)

							if res.StatusCode == 200 {
								for _, d := range datos.Data {
									fmt.Println("---")
									fmt.Printf("ID: %s\n", d.ID)
									fmt.Printf("Nombre: %s\n", d.Name)
									fmt.Printf("Apellido: %s\n", d.Last_name)
									fmt.Printf("RUT: %s\n", d.Rut)
									fmt.Printf("Email: %s\n", d.Email)

								}
								fmt.Println("---")
							} else {
								fmt.Println("Error al obtener el cliente!")
							}
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
								fmt.Printf("¡Cliente “%s” creado con éxito!\n", newclient.Name)
							}
						}
					case "5":
						{
							// Actualizar datos de un cliente por ID
							var id string
							fmt.Print("Ingrese el ID a modificar: ")
							fmt.Scanln(&id)
							fmt.Print("\n")

							fmt.Println("Rellene los datos. Espacio en blanco => dejar sin modificacion.")

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

							url := fmt.Sprintf("http://%s:%s/api/clients/%s", os.Getenv("HOST"), os.Getenv("PORT"), id)

							req, _ := http.NewRequest("PUT", url, data)
							req.Header.Set("Content-Type", "application/json")
							client := &http.Client{}
							res, _ := client.Do(req)
							if res.StatusCode == 200 {
								fmt.Println("Cliente modificado con exito!")
							}
						}
					case "6":
						{
							var id string
							fmt.Print("Ingrese el ID a modificar: ")
							fmt.Scanln(&id)
							fmt.Print("\n")
							// Borrar un cliente por ID
							url := fmt.Sprintf("http://%s:%s/api/clients/%s", os.Getenv("HOST"), os.Getenv("PORT"), id)
							req, _ := http.NewRequest("DELETE", url, nil)

							client := &http.Client{}
							res, _ := client.Do(req)
							if res.StatusCode == 200 {
								fmt.Println("Cliente eliminado con exito!")
							}
						}
					case "7":
						{
							menu_cliente = false
							menu_usuario = true

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
