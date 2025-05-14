package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cmerin0/tasky/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsers handles the fetching of all users
// @Summary Get all users
// @Description Fetch all users from the database
// @Success 200 {array} models.UserResponse
// @Failure 500 Internal Server Error
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.UserResponse
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		log.Error("Error fetching users: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.UserResponse
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Error("Cursor error: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
		})
	}

	log.Info("All users fetched successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
		"count": len(users),
	})
}

// CreateUser handles the creation of a new user
// @Summary Create a new user
// @Description Create a new user in the database
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 Bad Request
// @Failure 500 Internal server error
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		log.Error("Error parsing user data: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	collection := getUserCollection()
	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		log.Error("Error inserting user: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	log.Info("User created successfully")
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"userId":  result.InsertedID,
	})
}

// GetUser handles the fetching of a single user
// @Summary Get a user by ID
// @Description Fetch a user from the database by ID
// @Param userId path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 Status Not Found
// @Router /users/{userId} [get]
func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.UserResponse
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	userCollection := getUserCollection()
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		log.Error("Error fetching user: ", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	log.Info("User fetched successfully")
	return c.Status(http.StatusOK).JSON(user)
}

// UpdateUser handles the updating of a user
// @Summary Update a user by ID
// @Description Update a user in the database by ID
// @Param userId path string true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} models.UserResponse
// @Failure 400 Bad Request
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @Router /users/{userId} [put]
func UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	if err := c.BodyParser(&user); err != nil {
		log.Error("Error parsing user data: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	update := bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}

	userCollection := getUserCollection()
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		log.Error("Error updating user: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if result.MatchedCount == 0 {
		log.Error("User not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	log.Info("User updated successfully")
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

// DeleteUser handles the deletion of a user
// @Summary Delete a user by ID
// @Description Delete a user from the database by ID
// @Param userId path string true "User ID"
// @Success 200 {object} fiber.Map
// @Failure 404 Not Found
// @Failure 500 Internal Server Error
// @Router /users/{userId} [delete]
func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	userCollection := getUserCollection()
	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		log.Error("Error deleting user: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if result.DeletedCount == 0 {
		log.Error("User not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	log.Info("User deleted successfully")
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
