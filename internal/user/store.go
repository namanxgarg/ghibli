package user

import (
	"fmt"
	"sync"
)

var (
	users      = make(map[string]User) // user ID → User
	usersMutex sync.RWMutex            // handles concurrent access
)

// SaveUser adds a user to the store
func SaveUser(u User) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	users[u.ID] = u
	fmt.Println("✅ Saved user:", u.ID)
}

// FindUserByEmail returns a user by email (linear scan)
func FindUserByEmail(email string) (User, bool) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	for _, u := range users {
		if u.Email == email {
			return u, true
		}
	}

	return User{}, false
}

// FindUserByID returns a user by ID
func FindUserByID(id string) (User, bool) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	u, ok := users[id]
	return u, ok
}

func GetUserByID(id string) (User, bool) {
    usersMutex.RLock()
    defer usersMutex.RUnlock()
    user, ok := users[id]
    if !ok {
        fmt.Println("❌ GetUserByID failed for:", id)
    }
    return user, ok
}

