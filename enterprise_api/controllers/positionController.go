package controllers

import (
	"fmt"
	"net/http"

	// replace "yourpackage" with the actual package name
	"midterm/db"
	"midterm/models" // replace "yourpackage" with the actual package name

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreatePosition(c *gin.Context) {

	position := models.Position{}

	if err := c.ShouldBind(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := position.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DbConnect.Create(&position)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetPositions(c *gin.Context) {

	positions := []models.Position{}

	result := db.DbConnect.Find(&positions)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": positions})
}

func DeletePosition(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	result := db.DbConnect.Delete(&models.Position{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Position deleted successfully"})
}

func EditPosition(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the Position
	var position models.Position
	if err := db.DbConnect.First(&position, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}

	// Bind the request body to the Position
	if err := c.ShouldBind(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := position.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the changes
	if err := db.DbConnect.Save(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": position})
}
