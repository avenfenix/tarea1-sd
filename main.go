package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al leer el archivo .env")
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	conn_url := fmt.Sprintf("%s:%s", ip, port)

	fmt.Println("Tarea 1 Sistemas Distribuidos")

	r := gin.Default()

	r.Run(conn_url)
}
