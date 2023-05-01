package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/people", GetPeople)
	router.GET("/people/:id", GetPerson)
	router.POST("/people", CreatePerson)
	router.PUT("/people/:id", UpdatePerson)
	router.DELETE("/people/:id", DeletePerson)

	router.Run(":8080")
}

func GetPeople(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, age FROM people")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	people := []Person{}
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		people = append(people, p)
	}

	c.JSON(http.StatusOK, people)
}

func GetPerson(c *gin.Context) {
	id := c.Param("id")

	var p Person
	err := db.QueryRow("SELECT id, name, age FROM people WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func CreatePerson(c *gin.Context) {
	var p Person
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO people(name, age) VALUES(?, ?)", p.Name, p.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.ID = int(id)
	c.JSON(http.StatusOK, p)
}

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")

	var p Person
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE people SET name = ?, age = ? WHERE id = ?", p.Name, p.Age, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM people WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Person %s deleted", id)})
}
