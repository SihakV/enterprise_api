package controllers

import (
	"fmt"
	"log"
	"midterm/db"
	"midterm/models"
	"midterm/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var leaveFolder = "leaves"

func CreateLeave(c *gin.Context) {

	leave := models.Leave{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	if err := c.ShouldBind(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = leave.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var fileName string
	file, err := c.FormFile("image") // replace "file" with the name of your form field
	if err == nil {
		fileName, err = utils.GenerateFileName(leaveFolder, file.Filename)
		if err != nil {
			log.Printf("Can not generate filename\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating filename\n"})
		}

		err = utils.UploadFileToSpaces(fileName, file)
		if err != nil {
			fmt.Printf("Can not upload image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading image"})
			return
		}
	}

	leave.Status = 1
	leave.LeaveFile = &fileName
	leave.CreatedBy = int(auth.UserId)
	leave.ApprovedBy = 0

	result := db.DbConnect.Create(&leave)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetLeaves(c *gin.Context) {

	leaves := []models.Leave{}

	result := db.DbConnect.Find(&leaves)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": leaves})
}

// func GetLeave(c *gin.Context) {

// 	leaves := []models.Leave{}

// 	result := db.DbConnect.Find(&leaves)
// 	if result.Error != nil {
// 		// handle error, e.g. log it or return it in the HTTP response
// 		fmt.Println(result.Error)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": leaves})
// }

func GetLeaveEmployee(c *gin.Context) {

	leaves := []models.Leave{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	result := db.DbConnect.Where("created_by = ?", int(auth.UserId)).Find(&leaves)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": leaves})
}

func DeleteLeave(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	var leave models.Leave
	if err := db.DbConnect.First(&leave, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave not found"})
		return
	}

	if leave.LeaveFile != nil && *leave.LeaveFile != "" {
		err := utils.DeleteFileFromSpaces(*leave.LeaveFile)
		if err != nil {
			fmt.Printf("Can not delete image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
			return
		}
	}

	result := db.DbConnect.Delete(&leave)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Leave deleted successfully"})
}

func EditLeave(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the leave
	var leave models.Leave
	if err := db.DbConnect.First(&leave, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave not found"})
		return
	}

	// Bind the request body to the leave
	if err := c.ShouldBind(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := leave.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := *leave.LeaveFile
	file, err := c.FormFile("image") // replace "file" with the name of your form field
	if err == nil {

		if fileName != "" {
			err = utils.DeleteFileFromSpaces(fileName)
			if err != nil {
				fmt.Printf("Can not delete image\n")
				c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
				return
			}
		}

		fileName, err = utils.GenerateFileName(leaveFolder, file.Filename)
		if err != nil {
			log.Printf("Can not generate filename\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating filename\n"})
		}

		err = utils.UploadFileToSpaces(fileName, file)
		if err != nil {
			fmt.Printf("Can not upload image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading image"})
			return
		}
	}

	leave.LeaveFile = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update leave"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": leave})
}
