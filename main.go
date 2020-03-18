package main

import (
	"log"

	"go-app-template/cmd/serve"
	_ "go-app-template/docs"

	"github.com/spf13/cobra"
)

// @title go-app API
// @version 0.1
// @description go-app API doc

// @contact.name Qi
// @contact.email navono007@gmail.com

// @host localhost:8080
// @BasePath

func main() {
	root := &cobra.Command{
		Use:     "go-app.exe",
		Short:   "Example application",
		Long:    "An example application to show how to use the service-framework",
		Version: "0.2.0",
	}

	root.AddCommand(serve.NewCommand())
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
