package models

// IngredientParseRequest
type IngredientParseRequest struct {
    List    []string    `json:"list"`    
}

type RecipeBuilderResponse struct {
    Nutrition   Nutrition    `json:"nutrition"` 
    Errors      []string     `json:"errors"`
}

// IngredientParseResponse
type Amount struct {
    Unit    string      `json:"unit"`
    Value   float32     `json:"value"`
}

type Ingredient struct {
    /** Parser Response 
        {
            name: "all-purpose flour",
            amounts: [...]
            modifier: optional
        } 
    */
    Name        string      `json:"name"`
    Amounts     []Amount    `json:"amounts"`
    Modifier    string      `json:"modifier"`
}

