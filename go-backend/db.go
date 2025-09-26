package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var s3Client *s3.Client
var bucketName = "construction-expenses"

func initDB() {
	goErr := godotenv.Load()
	if goErr != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	)
	fmt.Println(connStr)
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}
	log.Println("Connected to Amazon RDS successfully!")
}

func initS3() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
}
