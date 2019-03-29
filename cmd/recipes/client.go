package recipes

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	maxPageSize = 10
)

// Client retrieves recipes from 3rd party API
type Client interface{ GetRecipe(string) ([]byte, error) }

// RecipeService
type RecipeService struct {
	Client Client
}

// PreparedRecipes is a map of recipies, with key specifying
// the preparation time in minutes
type PreparedRecipes map[int][]Recipe

func NewRecipeService() *RecipeService {
	return &RecipeService{
		Client: NewRecipeClient(),
	}
}

// GetSortedRecipes returns recipes for the specified ids
func (r *RecipeService) GetSortedRecipes(ids []string) ([]Recipe, error) {
	unsortedRecipes, err := r.getRecipesForIds(ids)
	if err != nil {
		return nil, err
	}

	return sortRecipes(unsortedRecipes), nil
}

// GetRecipesForRange returns  the recipes for a certain range,
// skipping the number of recipes specified with skip and getting
// either number of messages specified by maxPageSize or top (whichever is smaller)
func (r *RecipeService) GetRecipesForRange(skip int, top int) ([]Recipe, error) {
	rec := make([]Recipe, int(math.Min(float64(top), float64(maxPageSize))))

	end := skip + top
	counter := 0
	for i := skip; i <= end && i-skip < len(rec); i++ {
		var recipe Recipe
		recipeJSON, err := r.Client.GetRecipe(strconv.Itoa(i + 1))

		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(recipeJSON), &recipe)
		if err != nil {
			return nil, err
		}

		rec[counter] = recipe
		counter++
	}
	return rec, nil
}

func (r *RecipeService) getRecipesForIds(ids []string) (PreparedRecipes, error) {
	unsortedRecipes := make(PreparedRecipes)
	for _, id := range ids {
		var recipe Recipe

		recipeJSON, err := r.Client.GetRecipe(id)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(recipeJSON, &recipe)
		if err != nil {
			return nil, err
		}

		key, err := getPrepTimeInMinutes(recipe.PrepTime)

		if err != nil {
			return nil, err
		}
		unsortedRecipes[key] = append(unsortedRecipes[key], recipe)
	}
	return unsortedRecipes, nil
}

func sortRecipes(recipes PreparedRecipes) []Recipe {
	var prepTimes []int
	for k := range recipes {
		prepTimes = append(prepTimes, k)
	}
	sort.Ints(prepTimes)

	var orderedRecipes []Recipe
	for _, pt := range prepTimes {
		for _, recipe := range recipes[pt] {
			orderedRecipes = append(orderedRecipes, recipe)
		}
	}
	return orderedRecipes
}

func getPrepTimeInMinutes(preptime string) (int, error) {
	preptime = strings.Replace(preptime, "PT", "", 1)
	preptime = strings.Replace(preptime, "M", "", 1)
	return strconv.Atoi(preptime)

}
