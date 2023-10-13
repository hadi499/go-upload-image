package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Membuka koneksi ke database MySQL
	var err error
	db, err = sql.Open("mysql", "username:pass@tcp(localhost:3306)/go_upload")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	// Endpoint untuk mengunggah gambar dan data
	r.POST("/upload", uploadHandler)

	r.Run(":8080")
}

func uploadHandler(c *gin.Context) {
	// Retrieve the uploaded file from the form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a unique filename for the uploaded file
	filename := "uploads/" + file.Filename

	// Save the uploaded file to the server
	err = c.SaveUploadedFile(file, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save data to the database
	description := c.PostForm("description")

	_, err = db.Exec("INSERT INTO images (description, filename) VALUES (?, ?)", description, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data successfully uploaded to the database"})
}
