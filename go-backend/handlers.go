package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type Project struct {
	UserID      int    `db:"user_id" json:"user_id"`
	ProjectID   int64  `db:"project_id" json:"project_id"`
	ProjectName string `db:"project_name" json:"project_name"`
}

type Expense struct {
	ExpenseID         int64  `db:"expense_id" json:"expense_id"`
	ExpenseName       string `db:"expense_name" json:"expense_name"`
	ExpenseCost       int    `db:"expense_cost" json:"expense_cost"`
	ExpenseReceiptURL string `db:"expense_receipt_url" json:"expense_receipt_url"`
	ProjectID         int    `db:"project_id" json:"project_id"`
}

type User struct {
	UserID int    `db:"user_id" json:"user_id"`
	Email  string `db:"email" json:"email"`
	Name   string `db:"name" json:"name"`
	// DateCreated time.Time `db:"date_created" json:"date_created"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func getProjects(c *gin.Context) {
	var projects []Project
	userId := c.Param("user_id")
	fmt.Println("Getting projects for user ID with String type: ", userId)
	numUserId, err := strconv.Atoi(userId)
	fmt.Println("Getting projects for user ID with Int type: ", numUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not convert string to int" + err.Error()})
		return
	}
	err = db.Select(&projects, "SELECT * FROM projects WHERE user_id = $1", numUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not access database: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
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
	fmt.Println(c.Request)
	fmt.Println(c)
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		fmt.Println("ERROR: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ProjectID = time.Now().Unix()
	fmt.Println("PROJECT_NAME: ", project.ProjectName)
	fmt.Println("USER_ID: ", project.UserID)

	query := `INSERT INTO projects (project_id, project_name, user_id) VALUES ($1, $2, $3) RETURNING project_id`
	err := db.QueryRow(query, project.ProjectID, project.ProjectName, project.UserID).Scan(&project.ProjectID)
	fmt.Println("ERROR: ", err)
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

func handleUser(c *gin.Context, userEmail string, name string) (int64, error) {
	var userId int64
	// var dbName string
	// _ = db.QueryRow("SELECT current_database()").Scan(&dbName)

	// fmt.Println("Connected database:", dbName)

	// fmt.Println("Handling user: ", userEmail)
	err := db.Get(&userId, "SELECT user_id FROM users WHERE email=$1", userEmail)
	if err != nil {
		fmt.Println("User not found, creating new user: ", userEmail)
		userId = time.Now().Unix()
		fmt.Println("Handling user: ", userEmail)
		query := `INSERT INTO users(user_id, email, name) VALUES($1,$2,$3) RETURNING user_id;`
		queryErr := db.QueryRow(query, userId, userEmail, name).Scan(&userId)
		fmt.Println("ERROR: ", queryErr)
		if queryErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr})
			return 0, queryErr
		}
		fmt.Println("User created: ", userId)
		return userId, nil
	}
	fmt.Println("User found: ", userId)
	return userId, nil
}

func createJWT(userEmail string, userName string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  userEmail,
		"name":   userName,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 1 day expiration
	})
	return token.SignedString(jwtSecret)
}

func GoogleLogin(c *gin.Context) {
	var googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000",
		// RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	// log.Println("RedirectURL:", googleOauthConfig.RedirectURL)

	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Println("failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON: " + err.Error()})
		return
	}

	token, err := googleOauthConfig.Exchange(
		context.Background(),
		req.Code,
	)
	if err != nil {
		// log.Println("failed to exchange token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token exchange failed: " + err.Error()})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		// log.Println("failed to get id_token:", idToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing id_token: " + err.Error()})
		return
	}
	// fmt.Println("idToken: ", idToken)

	payload, err := idtoken.Validate(
		context.Background(),
		idToken,
		googleOauthConfig.ClientID,
	)
	// fmt.Println("Payload: ", payload)
	if err != nil {
		// log.Println("id_token validation failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid id token: " + err.Error()})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	//check in DB if user exists, if not create user record
	userId, handleUserErr := handleUser(c, email, name)
	fmt.Printf("User ID: %d\n", userId)
	if handleUserErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleUserErr.Error()})
		return
	}
	fmt.Println("Email: ", email)
	fmt.Println("Name: ", name)

	jwtToken, err := createJWT(email, name, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create JWT"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true in prod
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": userId})
}

// Middleware to check JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		fmt.Println("Auth Token: ", tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth cookie"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		fmt.Println("Claims: ", claims)
		if email, ok := claims["email"].(string); ok {
			fmt.Println("Setting email in context:", email)
			c.Set("email", email)
		}
		if name, ok := claims["name"].(string); ok {
			fmt.Println("Setting name in context:", name)
			c.Set("name", name)
		}
		if uI, ok := claims["userId"].(float64); ok {
			fmt.Println("Setting userId in context:", uI)
			userId := int(uI)
			c.Set("userId", userId)
		}
		c.Next()
	}
}

// Protected route example
func UserInfo(c *gin.Context) {
	email := c.GetString("email")
	name := c.GetString("name")
	userId := c.GetInt("userId")
	fmt.Println("UserInfo - Email:", email, "Name:", name, "UserID:", userId)
	c.JSON(http.StatusOK, gin.H{"message": "Hello " + email + ", your name is " + name, "email": email, "name": name, "user_id": userId})
}

func Logout(c *gin.Context) {
	// Clear the auth_token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true in prod
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
