// @title Todo API
// @version 1.0
// @description This is a sample API for managing todos
// @host localhost:3000
// @BasePath /

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"time"
	_ "todo-app/docs"
)

var session = DatabaseConnection()

type Todo struct {
	Id          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Notes       string    `json:"notes"`
	IsCompleted bool      `json:"completed"`
	Created_at  time.Time `json:"created_at"`
}

var bgCtx = context.Background()
var rdb = RedisInit()

func main() {
	config := LoadConfig()
	defer session.Close()
	session.Query(
		`CREATE TABLE IF NOT EXISTS todos (
			id uuid PRIMARY KEY,
			title TEXT,
			notes TEXT,
			isCompleted BOOLEAN,
			created_at timestamp
		)`,
	).Exec()
	fmt.Println("Table created successfully.")
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/todo", createTodo)
	router.GET("/todo", getAllTodos)
	router.GET("/todo/:id", getTodoById)
	router.PUT("/todo/:id", updateTodo)
	router.GET("/health", getHealth)

	router.Run(":" + config.PORT)
}

func getHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Server is working fine.")
}

func saveTodo(todo Todo) error {
	return session.Query(`
		INSERT INTO todos (id, title, notes, isCompleted, created_at)
		VALUES ( ?, ?, ?, ?, ?)
	`, todo.Id, todo.Title, todo.Notes, todo.IsCompleted, todo.Created_at).Exec()
}

// @Summary Create a todo
// @Description Create a new todo item
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body Todo true "Todo to create"
// @Success 201 {object} Todo
// @Router /todo [post]
func createTodo(ctx *gin.Context) {
	rdb.Del(bgCtx, "todos")
	var todo Todo
	id, _ := gocql.RandomUUID()
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.Id = id.String()
	todo.IsCompleted = false
	todo.Created_at = time.Now()
	if err := saveTodo(todo); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save todo"})
		return
	}
	ctx.JSON(http.StatusCreated, todo)
}

// @Summary Get all todos
// @Description Retrieve all todo items, cached with Redis
// @Tags todos
// @Produce  json
// @Success 200 {array} Todo
// @Router /todo [get]
func getAllTodos(ctx *gin.Context) {
	val, err := rdb.Get(bgCtx, "todos").Result()
	if err == nil {
		var cachedTodos []Todo
		if err := json.Unmarshal([]byte(val), &cachedTodos); err == nil {
			ctx.JSON(http.StatusOK, cachedTodos)
			return
		}
	}
	var todos []Todo

	iter := session.Query(`SELECT id, title, notes, isCompleted, created_at FROM todos`).Iter()

	var t Todo
	for iter.Scan(&t.Id, &t.Title, &t.Notes, &t.IsCompleted, &t.Created_at) {
		todos = append(todos, t)
		t = Todo{}
	}
	todoBytes, _ := json.Marshal(todos)
	rdb.Set(bgCtx, "todos", todoBytes, time.Hour)
	ctx.JSON(http.StatusOK, todos)
}

// @Summary Update a todo
// @Description Update a todo item by its ID
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path string true "Todo ID"
// @Param todo body Todo true "Updated todo data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todo/{id} [put]
func updateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	rdb.Del(bgCtx, "todos")
	var input Todo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := session.Query(`
		UPDATE todos 
		SET title = ?, notes = ?, isCompleted = ?
		WHERE id = ?`,
		input.Title, input.Notes, input.IsCompleted, id,
	).Exec()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

// @Summary Get a todo by ID
// @Description Retrieve a specific todo item by its ID
// @Tags todos
// @Produce  json
// @Param id path string true "Todo ID"
// @Success 200 {object} Todo
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todo/{id} [get]
func getTodoById(ctx *gin.Context) {
	id := ctx.Param("id")

	var todo Todo
	err := session.Query(`
		SELECT id, title, notes, isCompleted, created_at 
		FROM todos 
		WHERE id = ? LIMIT 1`,
		id,
	).Consistency(gocql.Quorum).Scan(
		&todo.Id, &todo.Title, &todo.Notes, &todo.IsCompleted, &todo.Created_at,
	)

	if err == gocql.ErrNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}
