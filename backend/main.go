package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("PUBLIC_ALLOWED_URLS")
	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		// Trim any whitespace from the origins
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
	} else {
		// Default to development environment
		origins = []string{"http://localhost:3000", "http://localhost:5173", "http://192.168.1.7:5173"}
	}

	// Configure CORS middleware
	corsConfig := cors.Config{
		AllowOrigins: strings.Join(origins, ","),
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}

	// Only allow credentials if we have specific origins
	if len(origins) > 0 && origins[0] != "*" {
		corsConfig.AllowCredentials = true
	}

	app.Use(cors.New(corsConfig))

	// MongoDB connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("fiveChan")

	// Get JWT secret from environment variable or generate one
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// If not set in environment, generate a random one
		// Note: This will generate a new secret on every restart
		// For production, you should set it as an environment variable
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		if err != nil {
			log.Fatal("Failed to generate JWT secret:", err)
		}
		jwtSecret = base64.StdEncoding.EncodeToString(secret)
		log.Println("WARNING: Using generated JWT secret. Set JWT_SECRET environment variable for production.")
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(db, jwtSecret)
	wsHandler := handler.NewWSHandler(db)

	// Routes
	app.Post("/signup", authHandler.Signup)
	app.Post("/login", authHandler.Login)
	app.Post("/random", authHandler.RandomAcc)

	// WebSocket route
	app.Get("/ws/expiry", wsHandler.Upgrade, websocket.New(wsHandler.ExpiryCheck))

	log.Fatal(app.Listen(":3000"))
}
