package order

import (
    "sync"
    "time"
)

type Order struct {
    ID         string
    UserID     string
    ImageFile  string
    PlacedAt   time.Time
    Status     string // "pending", "ready"
}

var (
    orders      []Order
    ordersMutex sync.RWMutex
)

func SaveOrder(o Order) {
    ordersMutex.Lock()
    defer ordersMutex.Unlock()
    orders = append(orders, o)
}

func GetOrdersByUser(userID string) []Order {
    ordersMutex.RLock()
    defer ordersMutex.RUnlock()

    var result []Order
    for _, o := range orders {
        if o.UserID == userID {
            result = append(result, o)
        }
    }
    return result
}

func MarkOrderReady(orderID string) {
    ordersMutex.Lock()
    defer ordersMutex.Unlock()
    for i, o := range orders {
        if o.ID == orderID {
            orders[i].Status = "ready"
            break
        }
    }
}
