package recipes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// RecipeHandler defines the HTTP handling function ServeHTTP
type RecipeHandler struct{}

// NewRecipesHandler returns empty RecipeHandler struct
func NewRecipesHandler() *RecipeHandler {
	return &RecipeHandler{}
}

func (h RecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe

	idsParam := r.URL.Query().Get("ids")
	topParam := r.URL.Query().Get("top")
	skipParam := r.URL.Query().Get("skip")

	rs := NewRecipeService()

	if len(idsParam) > 0 {
		ids := strings.Split(idsParam, ",")
		var err error
		recipes, err = rs.GetSortedRecipes(ids)

		if err != nil {
			http.Error(w, "Could not find specified recipes ", 500)
			log.Println(err.Error())
			return
		}
	} else if len(topParam) > 0 && len(skipParam) > 0 {
		skip, err := strconv.Atoi(skipParam)
		if err != nil {
			http.Error(w, "Bad top parameter ", 400)
			log.Println(err.Error())
			return
		}
		top, err := strconv.Atoi(topParam)
		if err != nil {
			http.Error(w, "Bad skip parameter ", 400)
			log.Println(err.Error())
			return
		}
		recipes, err = rs.GetRecipesForRange(skip, top)

		if err != nil {
			http.Error(w, "Could not find specified recipes ", 500)
			log.Println(err.Error())
			return
		}
	} else {
		http.Error(w, "Bad parameters  ", 400)
		return
	}

	result, err := json.Marshal(recipes)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = w.Write(result)

	if err != nil {
		log.Println(err.Error())
	}

}
