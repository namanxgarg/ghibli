package order

import (
    "fmt"
    "time"

    "github.com/namanxgarg/ghibli-backend/internal/user"
    "github.com/namanxgarg/ghibli-backend/pkg/notifier"
)

func SimulateModelRender(orderID string) {
    go func() {
        fmt.Printf("🛠️ Starting 3D render for order %s\n", orderID)
        time.Sleep(5 * time.Second) // simulate rendering delay
        MarkOrderReady(orderID)
        fmt.Printf("✅ Order %s is ready\n", orderID)

        // Find order to get userID and filename
        var order Order
        var found bool

        ordersMutex.RLock()
        for _, o := range orders {
            if o.ID == orderID {
                order = o
                found = true
                break
            }
        }
        ordersMutex.RUnlock()

        if !found {
            fmt.Printf("⚠️ Could not find order %s to send email\n", orderID)
            return
        }

        u, ok := user.GetUserByID(order.UserID)
        if !ok {
            fmt.Printf("⚠️ Could not find user %s for email\n", order.UserID)
            return
        }

        subject := "🎉 Your 3D model is ready!"
        body := fmt.Sprintf(
            "Hi there,\n\nYour model for image [%s] is ready to download or view.\n\nThanks!",
            order.ImageFile,
        )

        notifier.SendEmail(u.Email, subject, body)
    }()
}
