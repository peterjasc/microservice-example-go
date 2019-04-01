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

func (h RecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
			http.Error(w, "Could not find specified recipes ", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	} else if len(topParam) > 0 && len(skipParam) > 0 {
		skip, err := strconv.Atoi(skipParam)
		if err != nil {
			http.Error(w, "Bad top parameter ", http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		top, err := strconv.Atoi(topParam)
		if err != nil {
			http.Error(w, "Bad skip parameter ", http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		recipes, err = rs.GetRecipesForRange(skip, top)

		if err != nil {
			http.Error(w, "Could not find specified recipes ", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	} else {
		http.Error(w, "Bad parameters  ", http.StatusBadRequest)
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
