package fitbit

import (
    // misc.
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

    // logit libs 
    "logit/models"
)

const (
    // 304 -> 1 serving unit
    DefaultMeasurementId = 304
    BaseAuthorizationUrl = "https://www.fitbit.com/oauth2/authorize"
    FitbitApiUrl = "https://api.fitbit.com"
)

func GetAccessToken(code string) (*models.OAuth2Response, error) {
    // get all environment variables
    clientId := os.Getenv("CLIENT_ID")
    redirectUrl := os.Getenv("REDIRECT_URL")
    fitbitAuthToken := os.Getenv("TOKEN")

    client := &http.Client{
        Timeout: time.Second * 10,
    }

    // construct request
    data := url.Values{}
    data.Set("clientId", clientId)
    data.Set("grant_type", "authorization_code")
    data.Set("redirect_uri", redirectUrl)
    data.Set("code", code)

    tokenEndpoint := fmt.Sprintf("%s/oauth2/token", FitbitApiUrl)
    req, _ := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(data.Encode()))
    req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fitbitAuthToken))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    // send request
    res, err := client.Do(req)
    if err != nil {
        log.Printf("[LOGIN_HANDLER] error: %+v\n", err)
        return nil, err
    }
    
    // reads the response and converts it to an OAuth2Response type
    body, _ := io.ReadAll(res.Body)
    defer res.Body.Close()

    var authRes models.OAuth2Response 
    json.Unmarshal(body, &authRes)

    return &authRes, nil 
}

func GetUser(userId, accessToken string) (*models.User, error) {
    client := &http.Client{
        Timeout: time.Second * 10,
    }
    
    // construct the user profile endpoint
    profileEndpoint := fmt.Sprintf("%s/1/user/%s/profile.json", FitbitApiUrl, userId)
    req, _ := http.NewRequest(http.MethodGet, profileEndpoint, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Add("Content-Type", "application/json")

    // send request
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    // read response
    body, _ := io.ReadAll(res.Body)
    defer res.Body.Close()

    var user models.User
    json.Unmarshal(body, &user)
    
    return &user, nil
}

func RevokeToken(accessToken string) error {
    fitbitAuthToken := os.Getenv("TOKEN")

    client := http.Client{
        Timeout: time.Second * 10,
    }

    // revoke access_token
    revokeEndpoint := fmt.Sprintf("%s/oauth2/revoke", FitbitApiUrl)
    params := url.Values{}
    params.Set("token", accessToken)
    req, _ := http.NewRequest(http.MethodPost, revokeEndpoint, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fitbitAuthToken))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    // send request
    _, err := client.Do(req)
    return err
}

func LogFood(userId, accessToken string, logReq models.FoodLogRequest) (*http.Response, error)  {
    client := &http.Client{
        Timeout: time.Second * 10,
    }

    // construct query params
    params := url.Values{}
    params.Set("foodName", logReq.Name)
    params.Set("mealTypeId", fmt.Sprintf("%d", logReq.Meal))
    params.Set("unitId", fmt.Sprintf("%d", DefaultMeasurementId))
    params.Set("amount", ConvertFloat(logReq.Amount, 2))
    params.Set("date", time.Now().Format("2006-01-02"))

    // Set nutrition information
    params.Set("calories", ConvertFloat(logReq.Nutrition.Calories, 0))
    params.Set("totalFat", ConvertFloat(logReq.Nutrition.Fat, 2))
    params.Set("transFat", ConvertFloat(logReq.Nutrition.TransFat, 2))
    params.Set("saturatedFat", ConvertFloat(logReq.Nutrition.SaturatedFat, 2))
    params.Set("cholesterol", ConvertFloat(logReq.Nutrition.Cholesterol, 2))
    params.Set("sodium", ConvertFloat(logReq.Nutrition.Sodium, 2))
    params.Set("totalCarbohydrate", ConvertFloat(logReq.Nutrition.Carbohydrates, 2))
    params.Set("dietaryFiber", ConvertFloat(logReq.Nutrition.Fiber, 2))
    params.Set("sugars", ConvertFloat(logReq.Nutrition.Sugar, 2))
    params.Set("protein", ConvertFloat(logReq.Nutrition.Protein, 2))

    logEndpoint := fmt.Sprintf(
        "%s/1/user/%s/foods/log.json",
        FitbitApiUrl,
        userId,
    )

    req, _ := http.NewRequest(http.MethodPost, logEndpoint, nil)
    req.URL.RawQuery = params.Encode()
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Accept", "application/json")

    return client.Do(req)
}

func CreateFood(userId, accessToken string, createReq models.FoodCreateRequest) (*http.Response, error) {
    // ** make request to fitbit api **
    client := &http.Client{
        Timeout: time.Second * 10,
    }

    // construct query params
    // 304 -> 1 serving unit
    params := url.Values{}
    params.Set("name", createReq.Name)
    params.Set("defaultFoodMeasurementUnitId", fmt.Sprint(DefaultMeasurementId))
    params.Set("defaultServingSize", "1")
    params.Set("formType", "DRY")
    params.Set("description", createReq.Description)

    // Set nutrition information
    params.Set("calories", ConvertFloat(createReq.Nutrition.Calories, 0))
    params.Set("totalFat", ConvertFloat(createReq.Nutrition.Fat, 2))
    params.Set("transFat", ConvertFloat(createReq.Nutrition.TransFat, 2))
    params.Set("saturatedFat", ConvertFloat(createReq.Nutrition.SaturatedFat, 2))
    params.Set("cholesterol", ConvertFloat(createReq.Nutrition.Cholesterol, 2))
    params.Set("sodium", ConvertFloat(createReq.Nutrition.Sodium, 2))
    params.Set("totalCarbohydrate", ConvertFloat(createReq.Nutrition.Carbohydrates, 2))
    params.Set("dietaryFiber", ConvertFloat(createReq.Nutrition.Fiber, 2))
    params.Set("sugars", ConvertFloat(createReq.Nutrition.Sugar, 2))
    params.Set("protein", ConvertFloat(createReq.Nutrition.Protein, 2))

    createEndpoint := fmt.Sprintf(
        "%s/1/user/%s/foods.json",
        FitbitApiUrl,
        userId,
    )

    req, _ := http.NewRequest(http.MethodPost, createEndpoint, nil)
    req.URL.RawQuery = params.Encode()
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Accept", "application/json")

    return client.Do(req)
}
