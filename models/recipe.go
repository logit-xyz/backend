package models

type Thing struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	MainEntity  interface{} `json:"mainEntityOfPage"`
	Image       interface{} `json:"image"`
}

// schema.org/recipe
type Recipe struct {
	CookTime    string      `json:"cookTime"`
	PrepTime    string      `json:"prepTime"`
	TotalTime   string      `json:"totalTime"`
	Nutrition   interface{} `json:"nutrition"`
	Ingredients interface{} `json:"recipeIngredient"`
	Thing
}

