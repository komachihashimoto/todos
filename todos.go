package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./todos.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        done BOOLEAN
    )`)
    if err != nil {
        panic(err)
    }

    router := gin.Default()
    router.LoadHTMLGlob("templates/*.html")

    router.GET("/", func(ctx *gin.Context) {
        rows, err := db.Query("SELECT id, title, done FROM todos")
        if err != nil {
            return
        }
        defer rows.Close()

        var todos []Todo
        for rows.Next() {
            var todo Todo
            if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done); err != nil {
                return
            }
            todos = append(todos, todo)
        }

        ctx.HTML(200, "index.html", gin.H{
            "todos": todos,
        })
    })

    router.POST("/todo", func(ctx *gin.Context) {
        var newTodo Todo
        if err := ctx.BindJSON(&newTodo); err != nil {
            return
        }

        res, err := db.Exec("INSERT INTO todos (title, done) VALUES (?, ?)", newTodo.Title, newTodo.Done)
        if err != nil {
            return
        }

        id, err := res.LastInsertId()
        if err != nil {
            return
        }
        newTodo.ID = int(id)

        ctx.JSON(200, newTodo)
    })

    router.PUT("/todo/:index", func(ctx *gin.Context) {
        index, err := strconv.Atoi(ctx.Param("index"))
        if err != nil {
            return
        }

        var updatedTodo Todo
        if err := ctx.BindJSON(&updatedTodo); err != nil {
            return
        }

        _, err = db.Exec("UPDATE todos SET title = ?, done = ? WHERE id = ?", updatedTodo.Title, updatedTodo.Done, index)
        if err != nil {
            return
        }

        ctx.JSON(200, updatedTodo)
    })

    router.PATCH("/todo/:index/done", func(ctx *gin.Context) {
        index, err := strconv.Atoi(ctx.Param("index"))
        if err != nil {
            return
        }

        _, err = db.Exec("UPDATE todos SET done = ? WHERE id = ?", true, index)
        if err != nil {
            return
        }

        ctx.JSON(200, gin.H{"status": "done"})
    })

    router.DELETE("/todo/:index", func(ctx *gin.Context) {
        index, err := strconv.Atoi(ctx.Param("index"))
        if err != nil {
            return
        }

        _, err = db.Exec("DELETE FROM todos WHERE id = ?", index)
        if err != nil {
            return
        }

        ctx.JSON(200, gin.H{"status": "deleted"})
    })

    router.Run()
}
