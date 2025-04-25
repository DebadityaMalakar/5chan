package handler

import (
	"context"
	"time"

	model "backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WSHandler struct {
	userCollection *mongo.Collection
}

func NewWSHandler(db *mongo.Database) *WSHandler {
	return &WSHandler{
		userCollection: db.Collection("users"),
	}
}

// Upgrade to WebSocket connection
func (h *WSHandler) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// ExpiryCheck handles WebSocket connections for expiry notifications
func (h *WSHandler) ExpiryCheck(c *websocket.Conn) {
	// Create a ticker that checks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	// Channel to handle connection close
	done := make(chan struct{})

	// Start a goroutine to detect connection close
	go func() {
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				close(done)
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			// Find and delete expired random accounts
			expiredUsers, err := h.checkAndDeleteExpired()
			if err != nil {
				c.WriteJSON(fiber.Map{
					"error": "Failed to check expired accounts",
				})
				continue
			}

			// Notify frontend about deleted accounts
			for _, user := range expiredUsers {
				if err := c.WriteJSON(user); err != nil {
					return // Close connection on error
				}
			}

		case <-done:
			return
		}
	}
}

// checkAndDeleteExpired finds and deletes expired random accounts
func (h *WSHandler) checkAndDeleteExpired() ([]model.ExpiryNotification, error) {
	ctx := context.Background()
	now := time.Now()

	// Find expired random accounts
	filter := bson.M{
		"is_random": true,
		"expires_at": bson.M{
			"$lte": now,
		},
	}

	cursor, err := h.userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var expired []model.User
	if err = cursor.All(ctx, &expired); err != nil {
		return nil, err
	}

	// Delete expired accounts
	var notifications []model.ExpiryNotification
	for _, user := range expired {
		_, err := h.userCollection.DeleteOne(ctx, bson.M{"username": user.Username})
		if err != nil {
			continue // Skip if deletion fails
		}

		notifications = append(notifications, model.ExpiryNotification{
			Username:  user.Username,
			DeletedAt: now,
		})
	}

	return notifications, nil
}
