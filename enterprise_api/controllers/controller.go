package controllers

import (
	"fmt"
	"net/http"

	"midterm/db"     // replace "yourpackage" with the actual package name
	"midterm/models" // replace "yourpackage" with the actual package name

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreateRole(c *gin.Context) {
	name := c.PostForm("name")

	role := models.Category{
		Name: name,
	}
	
	result := db.DbConnect.Create(&role)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
