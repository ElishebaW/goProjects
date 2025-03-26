package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Initialize the SQLite driver
)

// Create a new database connection function
func dbConnect() (*sql.DB, error) {
	dsn := "./db.sqlite3"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Create users table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating users table: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

// Define a simple data model for our API
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Task represents a single task with its name and order.
type Task struct {
	Name     string `json:"name"`
	getOrder int    `json:"order"`
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %d: %s", t.getOrder(), t.Name)
}

var tasks = []Task{
	{"Work on task 1", 0},
	{"Work on task 2", 1},
	{"Work on task 3", 2},
	{"Take a break", 3},
	{"Work on task 4", 4},
	{"Work on task 5", 5},
}

func (ts *TaskManager) TasksList() []Task {
	return tasks
}

type TaskManager struct {
	Tasks []Task `json:"tasks"`
}

var manager = &TaskManager{}

// Define routes and handlers for your API
func main() {
	router := gin.Default()

	// Handle requests at the root path ("/")
	router.GET("/", rootHandler)

	api := router.Group("/api")
	{
		api.GET("/users", GetUsers)
		api.POST("/user", AddUser)
		api.POST("/addTask", AddTask)
	}

	// Run the server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

func rootHandler(c *gin.Context) {
	c.String(200, "Welcome to my API server!")
	if c.Request.URL.Path == "/" {
		c.File("index.html")
	} else if strings.HasPrefix(c.Request.URL.Path, "/static/") {
		// Serve static files from /static/
		c.File("./static" + c.Request.URL.Path)
	}
}

func GetUsers(c *gin.Context) {
	db, err := dbConnect()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Query database for users and return as JSON response
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// Return users as JSON response
	c.JSON(200, users)
}

func AddUser(c *gin.Context) {
	db, err := dbConnect()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Assuming you're using JSON for user data (use BindJSON for structured data)
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Insert user into database and return created status
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Optionally, get the last inserted ID
	lastID, _ := result.LastInsertId()
	c.JSON(201, gin.H{"id": lastID, "message": "User created successfully"})
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasks, err := manager.Tasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var data TaskData
	err := jsonnet.NewJSONLoader().Load(r.Body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := data.Name
	order := data.Order

	// use your local LLM to determine the best task order
	tasksList := manager.TasksList()
	for i, t := range tasksList {
		if t.Name == name {
			tasksList[i].getOrder() = order
			break
		}
	}

	if err := jsonnet.NewJSONLoader().Load([]byte(`{
                "tasks": ${tasksList}
            }`), &TaskManager{Tasks: tasksList}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add a Pomodoro timer feature

	fmt.Fprint(w, "Task added successfully!")
}

type TaskData struct {
	Name     string `json:"name"`
	getOrder int    `json:"order"`
}
