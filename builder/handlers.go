package builder

import (
	"log"
	"logit/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
    EMPTY_NAME = ""
)

func RecipeBuilderHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // accepts an ingredient list
        var req models.IngredientParseRequest
        if err := ctx.BindJSON(&req); err != nil {
            log.Printf("[BUILDER] malformed JSON input")
            ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response[interface{}]{
                Message: "doesn't follow expected input format",
                Data: nil,
                Status: http.StatusBadRequest,
            })
            return
        }

        // if it doesn't follow the form of the recipe, then there will be no amounts
        // so if there are no amounts, then it can't be an ingreident?
        result, err := ParseIngredients(req)
        if err != nil {
            log.Printf("[BUILDER] ingredient parser api failed")
            ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response[interface{}]{
                Message: "couldn't parse ingredients",
                Data: nil,
                Status: http.StatusBadRequest,
            })
            return
        }

        // parse could be successful vs unsucessful
        var success []models.Ingredient
        var exclude []string
        for i, item := range(result) {
            // exclude all strings (or lines) where the parse result
            // didn't find an ingredient name or an amount
            if item.Name == EMPTY_NAME || len(item.Amounts) == 0 {
                exclude = append(exclude, req.List[i])
            } else {
                success = append(success, item)
            }
        }
        
        // builds a query to my rustlang service to parse the ingredients
        // only query the DB on successful parses
        var recipeNutrition models.Nutrition
        for _, item := range(success) {
            amnt := item.Amounts[0]
            unit := amnt.Unit
            value := amnt.Value

            food := GetFood(item.Name);
            portions := GetAvailablePortions(food.FdcId)
            commonPortionIdx := FindCommonUnit(portions, unit)
            if commonPortionIdx != -1 {
                // do some calculation shit
                portion := portions[commonPortionIdx]
                servingGramWeight := (value / portion.Amount) * portion.GramWeight 
                multiplier := servingGramWeight / 100 // every food nutrient is for a 100g serving
                AddFoodNutritionalValue(&recipeNutrition, food, multiplier)
            } else {
                exclude = append(exclude, item.Name);
            }
        } 
    
        ctx.JSON(http.StatusOK, models.Response[models.RecipeBuilderResponse]{
            Message: "recipe built",
            Data: models.RecipeBuilderResponse{
                Nutrition: recipeNutrition,
                Errors: exclude,
            },
            Status: http.StatusOK,
        })

    }
}

func ImageUploadHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        file, _, err := ctx.Request.FormFile("image")
        if err != nil {
            log.Printf("[BUILDER] image upload: %+v", err)
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response[interface{}]{
                Message: err.Error(),
                Data: nil,
                Status: http.StatusInternalServerError,
            })
            return
        }
        
        if err != nil {
            log.Printf("[BUILDER] file open: %+v", err)
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response[interface{}]{
                Message: err.Error(),
                Data: nil,
                Status: http.StatusInternalServerError,
            })
            return
        }

        ingredientList, err := RunGoogleCloudOCR(file)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response[interface{}]{
                Message: err.Error(),
                Data: nil,
                Status: http.StatusInternalServerError,
            })
            return
        }
        
        // return successful response
        ctx.JSON(http.StatusOK, models.Response[[]string]{
            Message: "image parsed successfully",
            Data: ingredientList,
            Status: http.StatusOK,
        })
    }
}
