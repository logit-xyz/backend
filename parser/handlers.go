package parser

import (
	// misc.
	"encoding/json"
	"log"
	"logit/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"

	// http and web scraper libs
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func requestHandler() colly.RequestCallback {
    return func(r *colly.Request) {
        r.Headers.Set("Referer", "https://www.google.com")
        log.Printf("Visiting %s\n", r.URL)
    }
}

func scrapedHandler() colly.ScrapedCallback {
    return func(r *colly.Response) {
        log.Printf("Finished scraping %s\n", r.Request.URL)
    }
} 

func errorHandler() colly.ErrorCallback {
    return func(r *colly.Response, err error) {
        log.Printf("Scraping error: %+v\n", err)
    }
}

func htmlHandler(recipe *map[string]interface{}, id *string) colly.HTMLCallback {
    return func(h *colly.HTMLElement) {
        var rawJSON interface{}
        json.Unmarshal([]byte(h.Text), &rawJSON)
        
        data := FindRecipe(rawJSON)
        if len(data) > 0 {
            *recipe = data
            *id = uuid.New().String() 
        }

    }
}

func CalculateHandler(uagents []string) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        randomIndx := rand.Intn(len(uagents))
		uagent := uagents[randomIndx]
        
        rawRecipe := make(map[string]interface{}, 0)
        id := ""
       
        // -------- COLLY CONFIG --------
        c := colly.NewCollector(
            colly.UserAgent(uagent),
        )
        c.Limit(&colly.LimitRule{
            RandomDelay: 1 * time.Second,
        })
        
        // -------- COLLY HANDLERS --------
        c.OnRequest(requestHandler())
        c.OnError(errorHandler()) 
        c.OnScraped(scrapedHandler())
        c.OnHTML("script[type='application/ld+json']", htmlHandler(&rawRecipe, &id))
    
        // -------- COLLY START --------
        link := ctx.Query("link")
        c.Visit(link)
    
        var recipe models.Recipe
        bytes, _ := json.Marshal(rawRecipe)
        json.Unmarshal(bytes, &recipe)
        
        // Normalize nutrition, image, and main entity data
        recipe.Nutrition = NormalizeNutritionData(recipe.Nutrition)
        recipe.Image = NormalizeImageData(recipe.Image)
        recipe.MainEntity = NormalizeMainEntity(recipe.MainEntity)
    
        ctx.JSON(http.StatusOK, models.Response[models.Recipe]{
            Message: "recipe nutrition calculated!",
            Data: recipe,
            Status: http.StatusOK,
        })
    } 
}
