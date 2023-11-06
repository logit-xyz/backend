package models

type SessionData struct {
    AuthData    OAuth2Response  `json:"auth"`
    UserData    UserMetadata    `json:"user"`
}

type AccessTokenRequest struct {
    Code        string      `json:"code"` 
}

type Session struct {
    User            UserMetadata    `json:"user"`
    SessionId       string          `json:"sessionId"`
}

type Response[T any] struct {
    Message     string      `json:"message"`
    Data        T           `json:"data"`
    Status      int         `json:"status"`
}
