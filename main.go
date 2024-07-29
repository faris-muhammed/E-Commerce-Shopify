package main

import (
	"os"

	"main.go/initializer"
	"main.go/router"
)

func init() {
	initializer.Initialize()
}
func main() {
	r := router.SetupRouter()
	r.Run(os.Getenv("PORT"))
}
