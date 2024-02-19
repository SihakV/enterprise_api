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

var expenseFolder = "expenses"

func CreateExpense(c *gin.Context) {

	expense := models.Expense{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	if err := c.ShouldBind(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = expense.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fileName string
	file, err := c.FormFile("image") // replace "file" with the name of your form field
	if err == nil {
		fileName, err = utils.GenerateFileName(expenseFolder, file.Filename)
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

	expense.Status = 1
	expense.ExpenseFile = &fileName
	expense.CreatedBy = int(auth.UserId)
	expense.ApprovedBy = 0

	result := db.DbConnect.Create(&expense)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetExpenses(c *gin.Context) {

	expenses := []models.Expense{}

	result := db.DbConnect.Find(&expenses)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

func GetExpenseEmployee(c *gin.Context) {

	expenses := []models.Expense{}

	auth, err := GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not get auth user"})
		return
	}

	result := db.DbConnect.Where("created_by = ?", int(auth.UserId)).Find(&expenses)
	if result.Error != nil {
		// handle error, e.g. log it or return it in the HTTP response
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

func DeleteExpense(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	var expense models.Expense
	if err := db.DbConnect.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expense not found"})
		return
	}

	if expense.ExpenseFile != nil || *expense.ExpenseFile != "" {
		err := utils.DeleteFileFromSpaces(*expense.ExpenseFile)
		if err != nil {
			fmt.Printf("Can not delete image\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error deleting image"})
			return
		}
	}

	result := db.DbConnect.Delete(&expense)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Expense deleted successfully"})
}

func EditExpense(c *gin.Context) {
	// Get the ID from the URL
	id := c.Param("id")

	// Find the expense
	var expense models.Expense
	if err := db.DbConnect.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Bind the request body to the expense
	if err := c.ShouldBind(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := expense.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := *expense.ExpenseFile
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

		fileName, err = utils.GenerateFileName(expenseFolder, file.Filename)
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

	expense.ExpenseFile = &fileName

	// Save the changes
	if err := db.DbConnect.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense})
}
