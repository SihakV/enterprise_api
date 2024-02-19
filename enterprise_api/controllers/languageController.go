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

func CreateLanguage(c *gin.Context) {
	name := c.PostForm("name")

	language := models.Language{
		Name: name,
	}

	err := language.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DbConnect.Create(&language)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetLanguages(c *gin.Context) {

	languages := []models.Language{}

	result := db.DbConnect.Find(&languages)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": languages})
}

func DeleteLanguage(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	result := db.DbConnect.Delete(&models.Language{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Language deleted successfully"})
}

func EditLanguage(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the language
	var language models.Language
	if err := db.DbConnect.First(&language, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Language not found"})
		return
	}

	// Bind the request body to the language
	if err := c.ShouldBind(&language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the changes
	if err := db.DbConnect.Save(&language).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update language"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": language})
}
