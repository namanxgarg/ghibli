package user

import (
    "encoding/json"
    "net/http"
    "github.com/google/uuid"
    "github.com/namanxgarg/ghibli-backend/pkg/auth"
)


type SignupRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    var req SignupRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    _, exists := FindUserByEmail(req.Email)
    if exists {
        http.Error(w, "email already in use", http.StatusConflict)
        return
    }

    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        http.Error(w, "failed to hash password", http.StatusInternalServerError)
        return
    }

    u := User{
        ID:       uuid.New().String(),
        Email:    req.Email,
        Password: hashedPassword,
    }

    SaveUser(u)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "id": u.ID,
    })
}


type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    u, found := FindUserByEmail(req.Email)
    if !found || !auth.CheckPasswordHash(req.Password, u.Password) {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateJWT(u.ID)
    if err != nil {
        http.Error(w, "could not generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}
