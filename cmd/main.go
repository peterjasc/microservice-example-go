package main

import (
	"log"

	"github.com/peterjasc/recipes/cmd/recipes"
)

func main() {
	app, err := recipes.NewApp()
	if err != nil {
		log.Fatalln(err)
	}
	recipesHandler := recipes.NewRecipesHandler()
	app.Mux.HandleFunc("/recipes", recipesHandler.ServeHTTP)
	app.ListenAndServe()
}
