package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Task represents a single task
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Order       int    `json:"order"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

// Secret key for JWT signing
var jwtSecret = []byte("secret")

// Database connection helper
func dbConnect() (*sql.DB, error) {
	return sql.Open("sqlite3", "./pomodoro.db")
}

// Login handler to authenticate users and issue JWT tokens
func handleLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate user credentials (replace with your database logic)
	if user.Username != "testuser" || user.Password != "password123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Middleware to authenticate requests using JWT
func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		c.Abort()
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		c.Abort()
		return
	}

	c.Set("user_id", int(userID))
	c.Next()
}

// AddTaskHandler handles adding a new task
func AddTaskHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	task.UserID = userID.(int)

	db, err := dbConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO tasks (name, `order`, description, user_id) VALUES (?, ?, ?, ?)",
		task.Name, task.Order, task.Description, task.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task"})
		return
	}

	rows, err := db.Query("SELECT id, name, `order`, description FROM tasks WHERE user_id = ?", task.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Order, &t.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tasks"})
			return
		}
		tasks = append(tasks, t)
	}

	reorganizedTasks, err := reorganizeTasksWithLLM(tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorganize tasks"})
		return
	}

	for _, t := range reorganizedTasks {
		_, err := db.Exec("UPDATE tasks SET `order` = ? WHERE id = ?", t.Order, t.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tasks"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task added and tasks reorganized successfully"})
}

// Simulate sending tasks to Llama 3.2 LLM for reorganization
func reorganizeTasksWithLLM(tasks []Task) ([]Task, error) {
	fmt.Println("Sending tasks to Llama 3.2 LLM for reorganization...")
	for i := range tasks {
		tasks[i].Order = i // Example: Reorganize tasks by their index
	}
	return tasks, nil
}

func main() {
	router := gin.Default()

	// Public route
	router.POST("/login", handleLogin)

	// Protected routes
	router.Use(authMiddleware)
	router.POST("/tasks", AddTaskHandler)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
