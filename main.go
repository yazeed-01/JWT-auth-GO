package main

import (
	"JWTauth/initializers"
	"JWTauth/routes"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}
func main() {
	r := routes.SetupRoutes()
	r.Run()
}
