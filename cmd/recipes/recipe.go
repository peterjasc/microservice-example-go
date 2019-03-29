package recipes

// Recipe is the recipe info we get from a 3rd party
type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Headline    string       `json:"headline"`
	Description string       `json:"description"`
	Difficulty  int          `json:"difficulty"`
	PrepTime    string       `json:"prepTime"`
	ImageLink   string       `json:"imageLink"`
	Ingredients []ingredient `json:"ingredients"`
}

type ingredient struct {
	Name      string `json:"name"`
	ImageLink string `json:"imageLink"`
}
