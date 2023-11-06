package fitbit

import (
	// misc.
	"fmt"
	"io"
	"log"
	"net/http"

	// logit libs
	"logit/models"
	"logit/redis"

	"github.com/gin-gonic/gin"
)


func LogFoodHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // get the current session
        sess := ctx.Request.Header.Get("Authorization")
        if sessData, err := redis.GetSession(sess); err == nil {
            var body models.FoodLogRequest
            if err := ctx.ShouldBindJSON(&body); err != nil {
                ctx.AbortWithError(500, err)
            }

            res, err := LogFood(sessData.AuthData.UserId, sessData.AuthData.AccessToken, body)
            if err != nil {
                log.Printf("[LOG HANDLER] error: %+v", err)
                ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response[interface{}]{
                    Message: fmt.Sprintf("%v", err),
                    Data: nil,
                    Status: http.StatusInternalServerError,
                })  
            } 

            if res.StatusCode == http.StatusCreated {
                ctx.JSON(http.StatusCreated, models.Response[interface{}]{
                    Message: "food log added",
                    Data: nil,
                    Status: http.StatusCreated,
                })
            } else {
                body, _ := io.ReadAll(res.Body)
                defer res.Body.Close()

                ctx.JSON(res.StatusCode, models.Response[string]{
                    Message: "failed to log food",
                    Data: string(body),
                    Status: res.StatusCode,
                })
            }
        } else {
            log.Print("[LOG HANDLER] error: not authorized to make a food log request")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response[interface{}]{
                Message: "not authorized to create a food log",
                Data: nil,
                Status: http.StatusUnauthorized,
            })  
        }
    }
}

func CreateFoodHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // get the current session
        sess := ctx.Request.Header.Get("Authorization")
        if sessData, err := redis.GetSession(sess); err == nil {
            var body models.FoodCreateRequest
            if err := ctx.ShouldBindJSON(&body); err != nil {
                ctx.AbortWithError(500, err)
            }

            res, err := CreateFood(sessData.AuthData.UserId, sessData.AuthData.AccessToken, body)
            if err != nil {
                log.Printf("[CREATE HANDLER] error: %+v", err)
                ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response[interface{}]{
                    Message: fmt.Sprintf("%v", err),
                    Data: nil,
                    Status: http.StatusInternalServerError,
                })  
            } 

            if res.StatusCode == http.StatusCreated {
                ctx.JSON(http.StatusCreated, models.Response[interface{}]{
                    Message: "food created",
                    Data: nil,
                    Status: http.StatusCreated,
                })
            } else {
                body, _ := io.ReadAll(res.Body)
                defer res.Body.Close()

                ctx.JSON(res.StatusCode, models.Response[string]{
                    Message: "failed to create food",
                    Data: string(body),
                    Status: res.StatusCode,
                })
            }
        } else {
            log.Print("[LOG HANDLER] error: not authorized to make a food log request")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response[interface{}]{
                Message: "not authorized to create a food ",
                Data: nil,
                Status: http.StatusUnauthorized,
            })  
        }
    }
}
