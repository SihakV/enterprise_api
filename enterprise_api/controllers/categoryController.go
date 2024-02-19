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

func CreateCategory(c *gin.Context) {

	category := models.Category{}

	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := category.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DbConnect.Create(&category)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetCategorys(c *gin.Context) {

	categorys := []models.Category{}

	result := db.DbConnect.Find(&categorys)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categorys})
}

func DeleteCategory(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	result := db.DbConnect.Delete(&models.Category{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Position deleted successfully"})
}

func EditCategory(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the Position
	var category models.Category
	if err := db.DbConnect.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}

	// Bind the request body to the Position
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the changes
	if err := db.DbConnect.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}
