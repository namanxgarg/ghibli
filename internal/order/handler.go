package order

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/namanxgarg/ghibli-backend/pkg/auth"
)

type OrderRequest struct {
    Image string `json:"image"` // the image file user picked
}

func PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := auth.GetUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    var req OrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    order := Order{
        ID:        uuid.New().String(),
        UserID:    userID,
        ImageFile: req.Image,
        PlacedAt:  time.Now(),
        Status:    "pending",
    }

    SaveOrder(order)
    SimulateModelRender(order.ID)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "order_id": order.ID,
    })
}

func ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := auth.GetUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    userOrders := GetOrdersByUser(userID)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userOrders)
}
