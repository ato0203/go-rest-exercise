package recipes

// represents a recipe
type Recipe struct {
  Name        string       `json:"name" binding:"required"`
	Ingredients []Ingredient `json:"ingredients"`
}

// represents individual ingredients
type Ingredient struct {
  Name string `json:"name" binding:"required"`
}
