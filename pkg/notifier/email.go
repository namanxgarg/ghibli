package notifier

import (
    "fmt"
)

func SendEmail(to string, subject string, body string) {
    // Simulated email send
    fmt.Printf("📧 Email to %s\nSubject: %s\n\n%s\n\n", to, subject, body)
}
