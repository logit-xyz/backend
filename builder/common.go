package builder

import (
	// misc.
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// image manipulation & ocr
	"image"
	"image/jpeg"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
	"github.com/otiai10/gosseract/v2"

	"logit/models"
)

func FindCommonUnit(portions []Portion, unit string) int  {
    for i, portion := range(portions) {
        if (IsSameUnit(portion, unit)) {
            return i
        }
    } 
    return -1
}

func IsSameUnit(portion Portion, unit string) bool {
    return portion.UnitName == unit || portion.AbbrUnitName == unit;
}

func ComputeNutrientValue(value float32, multiplier float32) float64 {
    return float64(value * multiplier)
}

func AddFoodNutritionalValue(nutrition *models.Nutrition, food Foods, multiplier float32) {
    nutrition.Calories += ComputeNutrientValue(food.Calories, multiplier)
    nutrition.Fat += ComputeNutrientValue(food.TotalFat, multiplier)
    nutrition.TransFat += ComputeNutrientValue(food.TransFat, multiplier) 
    nutrition.SaturatedFat += ComputeNutrientValue(food.SaturatedFat, multiplier)
    nutrition.Cholesterol += ComputeNutrientValue(food.Cholesterol, multiplier)
    nutrition.Sodium += ComputeNutrientValue(food.Sodium, multiplier)
    nutrition.Carbohydrates += ComputeNutrientValue(food.TotalCarbs, multiplier)
    nutrition.Fiber += ComputeNutrientValue(food.DietaryFiber, multiplier)
    nutrition.Sugar += ComputeNutrientValue(food.Sugars, multiplier)
    nutrition.Protein += ComputeNutrientValue(food.Protein, multiplier)
}


func ConstructQueryString(list []string) string {
    params := url.Values{}
    for _, item := range(list) {
        params.Add("ing", item) 
    }
    return params.Encode()
}

func ParseIngredients(r models.IngredientParseRequest) ([]models.Ingredient, error) {
    // does the http request    
    client := http.Client {
        Timeout: 10 * time.Second,
    }
    
    PARSER_API_URL := os.Getenv("PARSER_API_URL")
    parseEndpoint := fmt.Sprintf("%s/parse", PARSER_API_URL)
    req, _ := http.NewRequest(http.MethodGet, parseEndpoint, nil)
    req.URL.RawQuery = ConstructQueryString(r.List)
    
    res, err := client.Do(req)
    if err != nil {
        fmt.Printf("[BUILDER] http error: %+v", err)
        return nil, err
    }

    var apiResponse models.Response[[]models.Ingredient]
    body, _ := io.ReadAll(res.Body)
    defer res.Body.Close()

    json.Unmarshal(body, &apiResponse)

    // set some security so only this Golang service can access the Parser API 
    // https://security.stackexchange.com/questions/255762/is-this-a-right-technique-to-create-and-validate-session-tokens
    
    return apiResponse.Data, nil
}

func RunGoogleCloudOCR(f io.Reader) ([]string, error) {
    visionCtx := context.Background()

    client, err := vision.NewImageAnnotatorClient(visionCtx)
    defer client.Close()
    if err != nil {
        log.Printf("[BUILDER] AnnotatorClient error: %+v", err)
        return nil, errors.New("failed to parse image")
    }

    image, err := vision.NewImageFromReader(f)
    annotation, err := client.DetectDocumentText(visionCtx, image, nil)

    var parsed []string
    if annotation == nil {
        log.Println("[BUILDER] Vision API found no text. Using TesseractOCR...")
        parsed, err = RunTesseractOCR(f)
        if err != nil {
            log.Printf("[BUILDER] TesseractOCR error: %+v", err)
            return nil, errors.New("failed to parse image")
        } 
    } else {
        parsed = strings.Split(annotation.Text, "\n")
    }
    
    return parsed, nil
}


func RunTesseractOCR(f io.Reader) ([]string, error) {
    imgObj, _, err := image.Decode(f)
    if err != nil {
        return nil, err
    }
    
    // Convert the image into bytes after a series of preprocessing
    imgBuf := PreProcessRecipeImage(imgObj)

    client := gosseract.NewClient();
    defer client.Close();

    client.SetImageFromBytes(imgBuf.Bytes())
    client.SetWhitelist(" -:/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    // Get text result from OCR
    text, _ := client.Text()
    return strings.Split(text, "\n"), nil
}

func PreProcessRecipeImage(img image.Image) *bytes.Buffer {
    // Runs a set of image proprocessing tasks:
    // 1. Grayscale image
    // 2. Change threshold segment with value 128
    grayscale := effect.Grayscale(img)
    threshold := segment.Threshold(grayscale, 128)
    buf := new(bytes.Buffer)
    jpeg.Encode(buf, threshold, nil)
    return buf
}
