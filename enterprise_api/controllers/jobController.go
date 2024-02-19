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

var jobFolder = "announcement"

func CreateJob(c *gin.Context) {

	job := models.Job{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	if err := c.ShouldBind(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = job.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fileName string
	file, err := c.FormFile("image") // replace "file" with the name of your form field
	if err == nil {
		fileName, err = utils.GenerateFileName(jobFolder, file.Filename)
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

	job.Status = 1
	job.CreatedBy = int(auth.UserId)
	job.AnnouncementImage = &fileName

	result := db.DbConnect.Create(&job)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetJobs(c *gin.Context) {

	jobs := []models.Job{}

	result := db.DbConnect.Find(&jobs)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

// func GetJob(c *gin.Context) {

// 	jobs := []models.Job{}

// 	result := db.DbConnect.Find(&jobs)
// 	if result.Error != nil {
// 		// handle error, e.g. log it or return it in the HTTP response
// 		fmt.Println(result.Error)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": jobs})
// }

func DeleteJob(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	var job models.Job
	if err := db.DbConnect.First(&job, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not found"})
		return
	}

	if job.AnnouncementImage != nil || *job.AnnouncementImage != "" {
		err := utils.DeleteFileFromSpaces(*job.AnnouncementImage)
		if err != nil {
			fmt.Printf("Can not delete image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
			return
		}
	}

	result := db.DbConnect.Delete(&job)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Job deleted successfully"})
}

func EditJob(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the job
	var job models.Job
	if err := db.DbConnect.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Bind the request body to the job
	if err := c.ShouldBind(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := job.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := *job.AnnouncementImage
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

		fileName, err = utils.GenerateFileName(jobFolder, file.Filename)
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

	job.AnnouncementImage = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": job})
}
