package controllers

import (
	"fmt"
	"log"
	"midterm/db"
	"midterm/models"
	"midterm/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var applicantFolder = "cv"

func CreateApplicant(c *gin.Context) {

	applicant := models.Applicant{}

	if err := c.ShouldBind(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := applicant.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobIdStr := c.Param("id")

	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		fmt.Printf("Cannot convert jobId to integer\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid jobId"})
		return
	}

	var fileName string
	file, err := c.FormFile("image") // replace "file" with the name of your form field
	if err == nil {
		fileName, err = utils.GenerateFileName(applicantFolder, file.Filename)
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

	applicant.Status = 1
	applicant.CvFile = &fileName
	applicant.JobId = jobId

	result := db.DbConnect.Create(&applicant)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetApplicants(c *gin.Context) {

	applicants := []models.Applicant{}

	result := db.DbConnect.Find(&applicants)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": applicants})
}

// func GetApplicant(c *gin.Context) {

// 	applicants := []models.Applicant{}

// 	result := db.DbConnect.Find(&applicants)
// 	if result.Error != nil {
// 		// handle error, e.g. log it or return it in the HTTP response
// 		fmt.Println(result.Error)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": applicants})
// }

func DeleteApplicant(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	var applicant models.Applicant
	if err := db.DbConnect.First(&applicant, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Applicant not found"})
		return
	}

	if applicant.CvFile != nil || *applicant.CvFile != "" {
		err := utils.DeleteFileFromSpaces(*applicant.CvFile)
		if err != nil {
			fmt.Printf("Can not delete image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
			return
		}
	}

	result := db.DbConnect.Delete(&applicant)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Applicant deleted successfully"})
}

func EditApplicant(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the applicant
	var applicant models.Applicant
	if err := db.DbConnect.First(&applicant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
		return
	}

	// Bind the request body to the applicant
	if err := c.ShouldBind(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := applicant.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := *applicant.CvFile
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

		fileName, err = utils.GenerateFileName(applicantFolder, file.Filename)
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

	applicant.CvFile = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&applicant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update applicant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": applicant})
}

func CountApplicants(c *gin.Context) {
	applicants := []models.Applicant{}

	result := db.DbConnect.Find(&applicants)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	count := len(applicants)

	c.JSON(http.StatusOK, gin.H{"total": count})
}
