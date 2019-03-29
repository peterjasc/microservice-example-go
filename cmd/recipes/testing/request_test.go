package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peterjasc/recipes/cmd/recipes"

	"github.com/stretchr/testify/assert"
)

func getMockServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))
}

func getMockNotFound() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", 404)
	}))
}
func TestGetRecipe(t *testing.T) {
	mockServer := getMockServer(mockResponse)

	client := recipes.NewRecipeClient()
	client.URL = mockServer.URL

	body, err := client.GetRecipe("2")

	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.Equal(t, 2683, len(body))
}

func TestNotFound(t *testing.T) {
	mockServer := getMockNotFound()

	client := recipes.NewRecipeClient()
	client.URL = mockServer.URL

	_, err := client.GetRecipe("113")

	assert.NotNil(t, err)
	assert.Equal(t, "404 - recipe not found", err.Error())
}
