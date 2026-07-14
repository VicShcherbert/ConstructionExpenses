package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initDB()
	initS3()
	defer db.Close()

	r := gin.Default()
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type"},
	}))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/get-projects/:user_id", getProjects)
	r.GET("/get-project/:project_id", getProject)
	r.POST("/create-project", addProject)
	r.POST("/create-expense", addExpenseToProject)
	r.GET("/get-expense/:expense_id", getExpense)
	r.GET("/get-project-expenses/:project_id", getProjectExpenses)
	r.DELETE("/delete-project/:project_id", deleteProject)
	r.DELETE("/delete-expense/:expense_id", deleteExpense)
	r.PUT("/update-project/:project_id", updateProject)
	r.PUT("/update-expense/:expense_id", updateExpense)
	r.POST("/upload-receipt/:expense_id", uploadReceipt)
	r.POST("/api/auth/google", GoogleLogin)
	r.GET("/user-info", AuthMiddleware(), UserInfo)
	// r.GET("/user-id", AuthMiddleware())
	r.POST("/logout", Logout)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
