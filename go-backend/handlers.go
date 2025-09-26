package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type Project struct {
	ProjectID   int64  `db:"project_id" json:"project_id"`
	ProjectName string `db:"project_name" json:"project_name"`
}

type Expense struct {
	ExpenseID         int64  `db:"expense_id" json:"expense_id"`
	ExpenseName       string `db:"expense_name" json:"expense_name"`
	ExpenseCost       int    `db:"expense_cost" json:"expense_cost"`
	ExpenseReceiptURL string `db:"expense_receipt_url" json:"expense_receipt_url"`
	ProjectID         int64  `db:"project_id" json:"project_id"`
}

func getProjects(c *gin.Context) {
	var project []Project
	err := db.Select(&project, "SELECT * FROM projects")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func getProject(c *gin.Context) {
	projectId := c.Param("project_id")
	numProjectId, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var project []Project
	queryError := db.Select(&project, "SELECT * FROM projects WHERE project_id = $1", numProjectId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func getExpense(c *gin.Context) {
	expenseId := c.Param("expense_id")
	numExpenseId, err := strconv.Atoi(expenseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var expense []Expense
	queryError := db.Select(&expense, "SELECT * FROM expenses WHERE expense_id = $1", numExpenseId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	c.JSON(http.StatusOK, expense)
}

func getProjectExpenses(c *gin.Context) {
	projectId := c.Param("project_id")
	numProjectId, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var expenses []Expense
	queryError := db.Select(&expenses, "SELECT * FROM expenses WHERE project_id = $1", numProjectId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func addProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ProjectID = time.Now().Unix()
	query := `INSERT INTO projects (project_id, project_name) VALUES ($1, $2) RETURNING project_id`
	err := db.QueryRow(query, project.ProjectID, project.ProjectName).Scan(&project.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

func addExpenseToProject(c *gin.Context) {
	var expense Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ExpenseID = time.Now().Unix()
	query := `INSERT INTO expenses (expense_id, expense_name, expense_cost, expense_receipt_url, project_id) VALUES ($1, $2, $3, $4, $5) RETURNING expense_id`
	err := db.QueryRow(query, expense.ExpenseID, expense.ExpenseName, expense.ExpenseCost, expense.ExpenseReceiptURL, expense.ProjectID).Scan(&expense.ExpenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, expense)
}

func deleteExpense(c *gin.Context) {
	expenseId := c.Param("expense_id")
	numExpenseId, err := strconv.Atoi(expenseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, queryError := db.Exec("DELETE FROM expenses WHERE expense_id=$1", numExpenseId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	count, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": rowsErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, count)
}

func deleteProject(c *gin.Context) {
	projectId := c.Param("project_id")
	numProjectId, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res, queryError := db.Exec("DELETE FROM projects WHERE project_id=$1", numProjectId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	count, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": rowsErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, count)
}

func updateProject(c *gin.Context) {
	projectId := c.Param("project_id")
	numProjectId, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, queryError := db.Exec("UPDATE projects SET project_name = $1 WHERE project_id=$2", project.ProjectName, numProjectId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	count, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": rowsErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, count)
}

func updateExpense(c *gin.Context) {
	expenseId := c.Param("expense_id")
	numExpenseId, err := strconv.Atoi(expenseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var expense Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, queryError := db.Exec("UPDATE expenses SET expense_name = $1, expense_cost = $2, expense_receipt_url = $3, project_id = $4 WHERE expense_id = $5", expense.ExpenseName, expense.ExpenseCost, expense.ExpenseReceiptURL, expense.ProjectID, numExpenseId)
	if queryError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryError.Error()})
		return
	}
	count, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": rowsErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, count)
}

func uploadReceipt(c *gin.Context) {
	expenseId := c.Param("expense_id")
	numExpenseId, err := strconv.Atoi(expenseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bucketName = "construction-expenses-all"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// key := fmt.Sprintf("receipts/%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	key := fmt.Sprintf("receipts/%d_%v", time.Now().Unix(), numExpenseId)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   src,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to S3"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload successful",
		"url":     fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key),
	})
}
