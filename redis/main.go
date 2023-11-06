package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
    "crypto/sha256"

	// logit libs
	"logit/models"

	// redis client lib
	goredis "github.com/go-redis/redis/v8"
)

var (
    redisCtx = context.Background()
    Client *goredis.Client
)

const Expiry = 604800 // ~1 week

func ConfigureRedis() {
    host := os.Getenv("REDIS_HOST")
    port := os.Getenv("REDIS_PORT")
    usr := os.Getenv("REDIS_USERNAME")
    pswd := os.Getenv("REDIS_PASSWORD")
    
    Client = goredis.NewClient(&goredis.Options{
        Addr: fmt.Sprintf("%s:%s", host, port),
        Username: usr,
        Password: pswd,
        DB: 0,
    })

    if err := Client.Ping(redisCtx).Err(); err != nil {
        log.Printf("[REDIS] ping error: %+v", err)
        return
    } else {
        log.Println("[REDIS] successfully connected to redis instance") 
    }
}

func SetSession(sessionId string, data models.SessionData) error {
    json, err := json.Marshal(data)
    if err != nil {
        log.Printf("[REDIS] marshal error: %+v", err)
        return err
    }
    
    // set's the session's TTL to the same as the FitBit API expiry duration 
    var d time.Duration = time.Second * Expiry
    err = Client.Set(redisCtx, sessionId, json, d).Err()
    if err != nil {
        log.Printf("[REDIS] error: %+v", err)
        return err
    }

    return nil
}

func GetSession(sessionId string) (*models.SessionData, error) {
    hashAlgo := sha256.New()
    hashedSessionId := hashAlgo.Sum([]byte(sessionId))

    sessionJSON, err := Client.Get(redisCtx, string(hashedSessionId)).Result()
    if err == goredis.Nil {
        // key doesn't exist
        log.Printf("[REDIS] session key doesn't exist")
        return nil, err
    } else if err != nil {
        log.Printf("[REDIS] session key doesn't exist")
        return nil, err
    }
    
    var sessData models.SessionData
    if err = json.Unmarshal([]byte(sessionJSON), &sessData); err != nil {
        log.Printf("[REDIS] unmarshal error: %+v", err)
        return nil, err
    }

    return &sessData, nil
}

func DeleteSession(sessionId string) error {
    hashAlgo := sha256.New()
    hashedSessionId := hashAlgo.Sum([]byte(sessionId))

    return Client.Del(redisCtx, string(hashedSessionId)).Err()
}

