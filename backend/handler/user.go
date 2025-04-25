package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"time"

	model "backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// AuthHandler manages user authentication operations
type AuthHandler struct {
	userCollection *mongo.Collection
	jwtSecret      []byte
}

// NewAuthHandler creates a new authentication handler instance
func NewAuthHandler(db *mongo.Database, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userCollection: db.Collection("users"),
		jwtSecret:      []byte(jwtSecret),
	}
}

// DatabaseUser represents the MongoDB document structure
type DatabaseUser struct {
	User      model.User `bson:"user"`
	ExpiresAt time.Time  `bson:"expires_at"`
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var newUser model.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	var existing DatabaseUser
	err := h.userCollection.FindOne(context.Background(), bson.M{"user.username": newUser.Username}).Decode(&existing)
	if err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Username exists"})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	pwHash, err := hashPasswordArgon2(newUser.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Password error"})
	}

	userDoc := DatabaseUser{
		User: model.User{
			Username: newUser.Username,
			Email:    newUser.Email,
			Password: pwHash.Hash,
			Salt:     pwHash.Salt,
			Format:   pwHash.Format,
		},
	}

	if _, err := h.userCollection.InsertOne(context.Background(), userDoc); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Create user failed"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User created",
		"user":    userDoc.User.Username,
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var dbUser DatabaseUser
	err := h.userCollection.FindOne(context.Background(), bson.M{"user.username": creds.Username}).Decode(&dbUser)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	user := dbUser.User
	valid := verifyPassword(creds.Password, user.Salt, user.Password, user.Format)
	if !valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := generateToken(user.Username, h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token error"})
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"user":    user.Username,
		"token":   token,
	})
}

// RandomAcc creates a temporary anonymous account
func (h *AuthHandler) RandomAcc(c *fiber.Ctx) error {
	username := "anon_" + mustGenerateString(12)
	password := mustGenerateString(16)
	pwHash, err := hashPasswordArgon2(password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Password error"})
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	user := DatabaseUser{
		User: model.User{
			Username: username,
			Password: pwHash.Hash,
			Salt:     pwHash.Salt,
			Format:   pwHash.Format,
			IsRandom: true,
		},
		ExpiresAt: expiresAt,
	}

	if _, err := h.userCollection.InsertOne(context.Background(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Create failed"})
	}

	token, err := generateToken(username, h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token error"})
	}

	return c.Status(201).JSON(fiber.Map{
		"username":   username,
		"password":   password,
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

// passwordHash stores hashed password details
type passwordHash struct {
	Hash   string `json:"hash" bson:"hash"`
	Salt   string `json:"salt" bson:"salt"`
	Format string `json:"format" bson:"format"`
}

// hashPasswordArgon2 creates secure password hash using Argon2
func hashPasswordArgon2(password string) (passwordHash, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return passwordHash{}, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		1,       // iterations
		64*1024, // memory
		4,       // threads
		32,      // key length
	)

	return passwordHash{
		Hash:   base64.RawStdEncoding.EncodeToString(hash),
		Salt:   base64.RawStdEncoding.EncodeToString(salt),
		Format: "argon2id",
	}, nil
}

// verifyPassword checks password against stored hash
func verifyPassword(password, salt, hash, format string) bool {
	if format == "argon2id" {
		return verifyArgon2(password, salt, hash)
	}
	return hashPasswordSHA256(password, salt) == hash
}

// verifyArgon2 verifies Argon2 hashed password
func verifyArgon2(password, saltBase64, hashBase64 string) bool {
	salt, err := base64.RawStdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false
	}

	computed := argon2.IDKey(
		[]byte(password),
		salt,
		1,
		64*1024,
		4,
		32,
	)

	return base64.RawStdEncoding.EncodeToString(computed) == hashBase64
}

// hashPasswordSHA256 creates legacy SHA256 hash
func hashPasswordSHA256(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// generateToken creates JWT for authenticated user
func generateToken(username string, secret []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}).SignedString(secret)
}

// mustGenerateString creates random string or panics
func mustGenerateString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[num.Int64()]
	}
	return string(result)
}
