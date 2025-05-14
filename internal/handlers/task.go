package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cmerin0/tasky/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxPaginationLimit = 30 // Maximum number of items to return in a single request

func CreateTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var task models.Task
	defer cancel()

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	newTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		UserID:      task.UserID,
	}

	taskCollection := getTaskCollection()
	result, err := taskCollection.InsertOne(ctx, newTask)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"taskId":  result.InsertedID,
	})
}

// GetTask handles the fetching of a single task
// @Summary Get a task by ID
// @Description Fetch a task from the database by ID
// @Param taskId path string true "Task ID"
func GetTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	var task models.Task
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	taskCollection := getTaskCollection()
	err := taskCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(task)
}

// GetTasks handles the fetching of all tasks
// @Summary Get all tasks
// @Description Fetch all tasks from the database no pagination
func GetTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var tasks []models.Task
	defer cancel()

	cursor, err := getTaskCollection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tasks",
		})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cursor error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// ListTasks handles the listing of tasks with pagination
// @Summary List tasks with pagination
// @Description Get a list of tasks with pagination
func ListTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	skip := (page - 1) * limit

	// Validate limit
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}

	// Get total count
	total, err := getTaskCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count tasks",
		})
	}

	// Find with pagination
	cursor, err := getTaskCollection().Find(ctx, bson.M{}, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tasks",
		})
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tasks",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func GetUserTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var tasks []models.Task
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	taskCollection := getTaskCollection()
	cursor, err := taskCollection.Find(ctx, bson.M{"userId": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	return c.Status(http.StatusOK).JSON(tasks)
}

func UpdateTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	var task models.Task
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	update := bson.M{
		"title":       task.Title,
		"description": task.Description,
		"completed":   task.Completed,
	}

	taskCollection := getTaskCollection()
	result, err := taskCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Task not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Task updated successfully"})
}

func DeleteTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	taskCollection := getTaskCollection()
	result, err := taskCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if result.DeletedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Task not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Task deleted successfully"})
}
