package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterjasc/microservice-example-go/cmd/recipes"
	pkgerrs "github.com/pkg/errors"
)

func main() {
	app, err := recipes.NewApp()
	if err != nil {
		log.Println(pkgerrs.Wrap(err, "failed to start app"))
		os.Exit(1)
	}
	recipesHandler := recipes.NewRecipeHandler()
	app.Mux.HandleFunc("/recipes", recipesHandler.ServeHTTP)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go app.ListenAndServe()

	<-stop
	app.Shutdown()
}
