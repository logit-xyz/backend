package models

// Authorization-related Types
type OAuth2Response struct {
    AccessToken     string  `json:"access_token"`
    TokenType       string  `json:"token_type"`
    RefreshToken    string  `json:"refresh_token"`
    UserId          string  `json:"user_id"` 
}

type User struct {
    Metadata UserMetadata `json:"user"`
}

type UserMetadata struct {
    EncodedId   string `json:"encodedId"`
    Avatar      string `json:"avatar"`
    Avatar150   string `json:"avatar150"`
    Avatar640   string `json:"avatar640"`
}

// Nutrition-related Types
type MealType int
const (
	Breakfast MealType = iota + 1
	MorningSnack
	Lunch
	AfternoonSnack
	Dinner
	Anytime MealType = 7
)

type Nutrition struct {
	Calories      float64 `json:"calories"`
	Fat           float64 `json:"fatContent"`
	TransFat      float64 `json:"transFatContent"`
	SaturatedFat  float64 `json:"saturatedFatContent"`
	Cholesterol   float64 `json:"cholesterolContent"`
	Sodium        float64 `json:"sodiumContent"`
	Carbohydrates float64 `json:"carbohydrateContent"`
	Fiber         float64 `json:"fiberContent"`
	Sugar         float64 `json:"sugarContent"`
	Protein       float64 `json:"proteinContent"`
}

type FoodLogRequest struct {
	Name      string    `json:"foodName"`
	Meal      MealType  `json:"mealTypeId"`
	UnitId    int       `json:"unitId"`
	Amount    float64   `json:"amount"`
	Nutrition Nutrition `json:"nutrition"`
}

// NOTE: calories must be a whole number
type FoodCreateRequest struct {
	Name        string    `json:"foodName"`
	UnitId      int       `json:"unitID"`
	ServingSize int       `json:"servingSize"`
	Calories    float64   `json:"calories"`
	Description string    `json:"description"`
	Nutrition   Nutrition `json:"nutrition"`
}
