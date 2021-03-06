package recipes

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/peterjasc/microservice-example-go/config"
	pkgerrs "github.com/pkg/errors"
)

// RecipeService wraps the Client that retrieves the recipes from 3rd party API
type RecipeService struct {
	Client Client
}

// Client retrieves recipes from 3rd party API
type Client interface{ GetRecipe(string) ([]byte, error) }

// PreparedRecipes is a map of recipies, with key specifying
// the preparation time in minutes
type PreparedRecipes map[int][]Recipe

type recipeWithErrors struct {
	Recipe   Recipe
	PrepTime int
	Error    error
}

func NewRecipeService() *RecipeService {
	return &RecipeService{
		Client: NewRecipeClient(),
	}
}

// GetSortedRecipes returns recipes for the specified ids
func (r *RecipeService) GetSortedRecipes(ids []string) ([]Recipe, error) {
	unsortedRecipes, err := r.getRecipesForIds(ids)
	if err != nil {
		return nil, pkgerrs.Wrap(err, "failed to get recipes")
	}

	return sortRecipes(unsortedRecipes), nil
}

// GetRecipesForRange returns  the recipes for a certain range,
// skipping the number of recipes specified with skip and getting
// either the number of messages specified by config.MaxPageSize or top
// (whichever is smaller). No sorting is taking place and results
// are ordered randomly.
func (r *RecipeService) GetRecipesForRange(skip int, top int) ([]Recipe, error) {
	rec := make([]Recipe, int(math.Min(float64(top), float64(config.MaxPageSize))))

	end := skip + top
	c := make(chan recipeWithErrors)
	for i := skip; i <= end && i-skip < len(rec); i++ {
		go func(i int) {
			recipeJSON, err := r.Client.GetRecipe(strconv.Itoa(i + 1))

			if err != nil {
				c <- recipeWithErrors{
					Recipe{},
					0,
					pkgerrs.Wrap(err, "failed to get recipe with id "+strconv.Itoa(i))}
				return
			}
			var recipe Recipe
			err = json.Unmarshal(recipeJSON, &recipe)
			if err != nil {
				c <- recipeWithErrors{
					Recipe{},
					0,
					pkgerrs.Wrap(err, "failed to unmarshal recipe with id "+strconv.Itoa(i))}
			}
			c <- recipeWithErrors{recipe, 0, nil}
		}(i)
	}

	err := getRecipesAsync(rec, c)

	if err != nil {
		return nil, pkgerrs.Wrap(err, "could not retrieve recipes asynchronously")
	}

	return rec, nil
}

func getRecipesAsync(recipes []Recipe, c chan recipeWithErrors) error {
	timeout := time.After(config.AsyncRecieveTimeout)
	var count int
	for i := 0; i < len(recipes); i++ {
		select {
		case rWE := <-c:
			if rWE.Error != nil {
				return pkgerrs.Wrap(rWE.Error, "received error via channel")
			}
			recipes[count] = rWE.Recipe
			count++
		case <-timeout:
			return errors.New("asynchronous retrieval timeout exceeded")
		}
	}
	return nil
}

func getPreparedRecipesAsync(unsortedRecipes PreparedRecipes, c chan recipeWithErrors, idsLen int) error {
	timeout := time.After(config.AsyncRecieveTimeout)
	for i := 0; i < idsLen; i++ {
		select {
		case rWE := <-c:
			if rWE.Error != nil {
				return pkgerrs.Wrap(rWE.Error, "received error via channel")
			}
			unsortedRecipes[rWE.PrepTime] = append(unsortedRecipes[rWE.PrepTime], rWE.Recipe)
		case <-timeout:
			return errors.New("asynchronous retrieval timeout exceeded")
		}
	}
	return nil
}

func (r *RecipeService) getRecipesForIds(ids []string) (PreparedRecipes, error) {
	unsortedRecipes := make(PreparedRecipes)
	c := make(chan recipeWithErrors)
	for _, id := range ids {
		go func(id string) {

			recipeJSON, err := r.Client.GetRecipe(id)

			if err != nil {
				c <- recipeWithErrors{
					Recipe{},
					0,
					pkgerrs.Wrap(err, "failed to get recipe with id "+id)}
				return
			}

			var recipe Recipe
			err = json.Unmarshal(recipeJSON, &recipe)
			if err != nil {
				c <- recipeWithErrors{
					Recipe{},
					0,
					pkgerrs.Wrap(err, "failed to unmarshal recipe with id "+id)}
				return
			}

			key, err := getPrepTimeInMinutes(recipe.PrepTime)

			if err != nil {
				c <- recipeWithErrors{
					Recipe{},
					0,
					pkgerrs.Wrap(err, "get preptime for recipe with id "+id)}
				return
			}

			c <- recipeWithErrors{recipe, key, err}
		}(id)
	}

	err := getPreparedRecipesAsync(unsortedRecipes, c, len(ids))

	if err != nil {
		return nil, pkgerrs.Wrap(err, "could not retrieve recipes asynchronously")
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
