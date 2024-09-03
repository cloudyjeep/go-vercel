package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "./items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Printf("db: %v\n", db)

	// Create table if not exists
	createTable(db)

	// Define routes
	router.GET("/items", func(c *gin.Context) {
		items, err := getItems(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	router.POST("/items", func(c *gin.Context) {
		var newItem Item
		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := createItem(db, newItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newItem)
	})

	router.PUT("/items/:id", func(c *gin.Context) {
		var updatedItem Item
		if err := c.ShouldBindJSON(&updatedItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedItem.ID, _ = strconv.Atoi(c.Param("id"))
		if err := updateItem(db, updatedItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedItem)
	})

	router.DELETE("/items/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := deleteItem(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
	})
 
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"  }) 
	})
 
	// Run the server
	router.Run(":8080")
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"price" INTEGER
	);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func getItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query("SELECT id, name, price FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func createItem(db *sql.DB, item Item) error {
	insertSQL := `INSERT INTO items(name, price) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(item.Name, item.Price)
	return err
}

func updateItem(db *sql.DB, item Item) error {
	updateSQL := `UPDATE items SET name = ?, price = ? WHERE id = ?`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(item.Name, item.Price, item.ID)
	return err
}

func deleteItem(db *sql.DB, id int) error {
	deleteSQL := `DELETE FROM items WHERE id = ?`
	statement, err := db.Prepare(deleteSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	return err
}
