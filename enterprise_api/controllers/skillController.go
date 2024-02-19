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

func CreateSkill(c *gin.Context) {

	skill := models.Skill{}

	if err := c.ShouldBind(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := skill.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DbConnect.Create(&skill)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetSkills(c *gin.Context) {

	skills := []models.Skill{}

	result := db.DbConnect.Find(&skills)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skills})
}

func DeleteSkill(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	result := db.DbConnect.Delete(&models.Skill{}, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Skill deleted successfully"})
}

func EditSkill(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the Skill
	var skill models.Skill
	if err := db.DbConnect.First(&skill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	// Bind the request body to the Skill
	if err := c.ShouldBind(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := skill.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the changes
	if err := db.DbConnect.Save(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
