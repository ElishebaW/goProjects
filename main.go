package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Create a new database connection function
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println(err)
		panic("Unable to connect to the database.")
	}
	return db
}

// Define a simple data model for our API
type User struct {
	ID    int
	Name  string
	Email string
}

// Define routes and handlers for your API
func main() {
	router := gin.Default()

	// Handle requests at the root path ("/")
	router.GET("/", rootHandler)

	api := router.Group("/api")
	{
		api.GET("/users", GetUsers)
		api.POST("/user", AddUser)
		api.GET("/")
	}

	// Run the server on port 8080
	router.Run(":8080")
}

func rootHandler(c *gin.Context) {
	c.String(200, "Welcome to my API server!")
}

func GetUsers(c *gin.Context) {
	db := dbConnect()
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
		user := User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(200, users)
}

func AddUser(c *gin.Context) {
	db := dbConnect()
	defer db.Close()

	name := c.PostForm("name")
	email := c.PostForm("email")

	// Insert user into database and return created status
	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(201)
}
