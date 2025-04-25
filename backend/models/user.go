package model

import "time"

type User struct {
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"`
	Salt      string    `json:"-" bson:"salt"`
	Format    string    `json:"-" bson:"format"`    // New field to track password hashing algorithm
	IsRandom  bool      `json:"-" bson:"is_random"` // Flag to identify random accounts
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type ExpiryNotification struct {
	Username  string    `json:"username"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UserWithExpiry struct is now redundant since the User struct itself has ExpiresAt field
// But keeping it for backward compatibility
type UserWithExpiry struct {
	User      `bson:",inline"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}
