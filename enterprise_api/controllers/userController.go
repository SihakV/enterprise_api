package controllers

import (
	"fmt"
	"log"
	"net/http"

	// replace "yourpackage" with the actual package name
	"midterm/db"
	"midterm/models" // replace "yourpackage" with the actual package name
	"midterm/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var profileFolder = "profile"

func CreateUser(c *gin.Context) {

	user := models.User{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	if err = c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	user.Password = string(hashedPassword)

	var fileName string
	file, err := c.FormFile("profile") // replace "file" with the name of your form field
	if err == nil {
		fileName, err = utils.GenerateFileName(profileFolder, file.Filename)
		if err != nil {
			log.Printf("Can not generate filename\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating filename\n"})
		}

		err = utils.UploadFileToSpaces(fileName, file)
		if err != nil {
			fmt.Printf("Can not upload profile\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading profile"})
			return
		}
	}

	user.Status = 1
	user.CreatedBy = int(auth.UserId)
	user.Profile = &fileName

	result := db.DbConnect.Create(&user)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetUsers(c *gin.Context) {

	users := []models.User{}

	result := db.DbConnect.Find(&users)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func DeleteUser(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	var user models.User
	if err := db.DbConnect.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if user.Profile != nil {
		err := utils.DeleteFileFromSpaces(*user.Profile)
		if err != nil {
			fmt.Printf("Can not delete profile\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting profile"})
			return
		}
	}

	result := db.DbConnect.Delete(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "User deleted successfully"})
}

func EditUser(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the user
	var user models.User
	if err := db.DbConnect.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind the request body to the user
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := user.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	user.Password = string(hashedPassword)

	fileName := *user.Profile
	file, err := c.FormFile("profile") // replace "file" with the name of your form field
	if err == nil {

		if fileName != "" {
			err = utils.DeleteFileFromSpaces(fileName)
			if err != nil {
				fmt.Printf("Can not delete image\n")
				c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
				return
			}
		}

		fileName, err = utils.GenerateFileName(profileFolder, file.Filename)
		if err != nil {
			log.Printf("Can not generate filename\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating filename\n"})
		}

		err = utils.UploadFileToSpaces(fileName, file)
		if err != nil {
			fmt.Printf("Can not upload profile\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading profile"})
			return
		}
	}

	user.Profile = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func EditProfile(c *gin.Context) {
	// Get the ID from the URL

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	userId := int(auth.UserId)

	// Find the user
	var user models.User
	if err := db.DbConnect.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind the request body to the user
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	user.Password = string(hashedPassword)

	fileName := *user.Profile
	file, err := c.FormFile("profile") // replace "file" with the name of your form field
	if err == nil {

		if fileName != "" {
			err = utils.DeleteFileFromSpaces(fileName)
			if err != nil {
				fmt.Printf("Can not delete image\n")
				c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
				return
			}
		}

		fileName, err = utils.GenerateFileName(profileFolder, file.Filename)
		if err != nil {
			log.Printf("Can not generate filename\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error generating filename\n"})
		}

		err = utils.UploadFileToSpaces(fileName, file)
		if err != nil {
			fmt.Printf("Can not upload profile\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading profile"})
			return
		}
	}

	user.Profile = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetProfile(c *gin.Context) {

	user := []models.User{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	result := db.DbConnect.Where("user_id = ?", int(auth.UserId)).First(&user)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
