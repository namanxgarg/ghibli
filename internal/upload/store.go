package upload

import (
    "sync"
    "time"
)

type Upload struct {
    Filename   string
    UserID     string
    UploadedAt time.Time
}

var (
    uploads      []Upload
    uploadsMutex sync.RWMutex
)

func SaveUpload(u Upload) {
    uploadsMutex.Lock()
    defer uploadsMutex.Unlock()
    uploads = append(uploads, u)
}

func GetUploadsByUser(userID string) []Upload {
    uploadsMutex.RLock()
    defer uploadsMutex.RUnlock()

    var result []Upload
    for _, u := range uploads {
        if u.UserID == userID {
            result = append(result, u)
        }
    }
    return result
}
