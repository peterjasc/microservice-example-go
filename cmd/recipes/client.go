package recipes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/peterjasc/microservice-example-go/config"
)

// RecipeClient handles requests for recipes to 3rd party API
type RecipeClient struct {
	Client *http.Client
	URL    string
}

// NewRecipeClient sets up the fields required to make the request for recipes to 3rd party API
func NewRecipeClient() *RecipeClient {
	return &RecipeClient{
		Client: newHTTPClient(),
		URL:    config.RecipesClientAPIURL,
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: config.HTTPClientTimeout,
	}
}

// GetRecipe returns a byte stream for a recipe id
func (c *RecipeClient) GetRecipe(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.URL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode) + " - recipe not found")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
