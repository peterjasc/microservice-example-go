package testing

import (
	"errors"
	"testing"

	"github.com/peterjasc/recipes/cmd/recipes"
	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (mc MockClient) GetRecipe(id string) ([]byte, error) {
	if id == "2" {
		return []byte(mockResponse), nil
	}

	if id == "3" {
		return []byte(mockResponse2), nil
	}
	return []byte{}, errors.New("unknown id")
}

func TestGetSortedRecipes(t *testing.T) {
	client := MockClient{}

	recService := &recipes.RecipeService{
		Client: client,
	}

	ids := []string{"2", "3"}
	actual, err := recService.GetSortedRecipes(ids)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "3", actual[0].ID)
	assert.Equal(t, "2", actual[1].ID)
}

func TestGetRecipesForRange(t *testing.T) {
	client := MockClient{}

	recService := &recipes.RecipeService{
		Client: client,
	}

	actual, err := recService.GetRecipesForRange(1, 2)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(actual))
	assert.True(t, (actual[0].ID == "3" && actual[1].ID == "2") || (actual[0].ID == "2" && actual[1].ID == "3"))
}
