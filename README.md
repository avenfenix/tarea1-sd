# Tarea 1

## Paquetes

`go mod init`

### [Gin](https://gin-gonic.com/)
`go get -u github.com/gin-gonic/gin`


### [Godotenv](https://github.com/joho/godotenv)
`go get -u github.com/joho/godotenv`


### MongoDB Driver
`go get -u go.mongodb.org/mongo-driver/mongo`

## Funcionalidades de la API

- Registrar y consultar clientes de la empresa. 
- Ingresar como usuario para obtener su ID. 
- Consultar a la API de iLovePDF para proteger los archivos subidos.


## iLovePDF
- [x] Pedir Bearer Token
- [x] Iniciar tarea de proteccion
- [x] Subir archivo PDF
- [x] Procesar (proteger) archivo PDF 
- [ ] Descargar archivo PDF protegido

## Endpoints

#### `/register`


- [ ] No filtra bien a los ya registrados (correo y rut).

## Persistencia de datos

- [x] Clientes
- [x] Usuarios
- [ ] Tokens

## Referencias

### Go

- [Diferencia entre go get/go install](https://stackoverflow.com/questions/24878737/what-is-the-difference-between-go-get-and-go-install)
- [Variables de entorno](https://blog.friendsofgo.tech/posts/trabajando-con-variables-de-entorno-en-go/)
- [Encoding and Decoding JSON with http request](https://kevin.burke.dev/kevin/golang-json-http/#:~:text=type%20User%20struct%7B%20Id%20string%20Balance%20uint64%20%7D,first%20and%20then%20copy%20that%20to%20a%20reader.)
- [How to parse JSON](https://hackajob.com/talent/blog/how-to-parse-json-from-apis-in-golang)


### Gin


- [Bind query string or post data](https://gin-gonic.com/docs/examples/bind-query-or-post/)
- [Bind Uri](https://gin-gonic.com/docs/examples/bind-uri/)
- [Rendering](https://gin-gonic.com/es/docs/examples/rendering/)

### Contenido

- [Software como servicio - Wikipedia](https://es.wikipedia.org/wiki/Software_como_servicio)
- [Metodologia "The twelve-factor app"](https://12factor.net/es/)
- [Rest - Wikipedia](https://es.wikipedia.org/wiki/Transferencia_de_Estado_Representacional)
- [HTTP - Request Methods](https://en.wikipedia.org/wiki/HTTP#Request_methods)
- [JSON](https://www.json.org/json-en.html)
- [JSON & BSON](https://www.mongodb.com/json-and-bson)


### MongoDB

- [MongoDB Go Driver - Documentation](https://www.mongodb.com/docs/drivers/go/current/quick-start/)
- [Golang and MongoDB](https://www.mongodb.com/languages/golang)
- [MongoDB Atlas Golang Sample Project](https://github.com/mongodb-university/atlas_starter_go)
- [MongoDB Go Driver reference](https://www.mongodb.com/docs/drivers/go/current/#introduction)
- [Get _id from document](https://dev.to/yasaricli/getting-mongodb-id-for-go-4e05)


### Operaciones CRUD

- [Modificar documento](https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/modify/)
- [UpdateOne](https://www.mongodb.com/docs/drivers/go/current/usage-examples/updateOne/)
- [Interpretar UpdateResult](https://stackoverflow.com/questions/76232471/how-can-i-read-data-from-mongo-updateresult-type-in-golang-updateone-addtose)