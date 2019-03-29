package recipes

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// RecipeClient handles requests for recipes to 3rd party API
type RecipeClient struct {
	Client *http.Client
	URL    string
}

// NewRecipeClient sets up the fields required to make the request for recipes to 3rd party API
func NewRecipeClient() *RecipeClient {
	return &RecipeClient{
		Client: http.DefaultClient,
		URL:    "https://s3-eu-west-1.amazonaws.com/test-golang-recipes",
	}
}

// GetRecipe returns a byte stream for a recipe id
func (c *RecipeClient) GetRecipe(id string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.URL+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(strconv.Itoa(http.StatusNotFound) + " - recipe not found")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
