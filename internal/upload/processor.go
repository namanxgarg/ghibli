package upload

import (
    "fmt"
    "time"
)

// SimulateRender simulates rendering after a delay
func SimulateRender(filename string) {
    go func() {
        fmt.Printf("🎨 Starting render for %s\n", filename)
        time.Sleep(3 * time.Second) // simulate rendering delay

        uploadsMutex.Lock()
        for i, u := range uploads {
            if u.Filename == filename {
                uploads[i].Status = "rendered"
                break
            }
        }
        uploadsMutex.Unlock()
        fmt.Printf("✅ Done rendering %s\n", filename)
    }()
}
