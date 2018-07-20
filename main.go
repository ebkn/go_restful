package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// open ad db connection
	var err error
	db, err = gorm.Open("mysql", "root:@/go_restful?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// migrate the schema
	db.AutoMigrate(&todoModel{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.PATCH("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run()
}

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

// add a new todo
func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":     http.StatusCreated,
			"message":    "Todo item created successfully!",
			"resourceId": todo.ID,
		},
	)
}

// fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			},
		)
	}

	// transform todos
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{
			ID:        item.ID,
			Title:     item.Title,
			Completed: completed,
		})
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _todos,
		},
	)
}

// fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			},
		)
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _todo,
		},
	)
}

// update a todo
func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No too found!",
			},
		)
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Todo updated successfully!",
		},
	)
}

// remove a todo
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No todo found!",
			},
		)
		return
	}

	db.Delete(&todo)
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Todo deleted successfully!",
		},
	)
}
