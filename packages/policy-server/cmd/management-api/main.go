package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	postgres_driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)
import "github.com/odahu/odahu-flow/packages/policy-server/pkg/handler"
import "github.com/odahu/odahu-flow/packages/policy-server/pkg/store/postgres"

func main() {

	connString := flag.String("conn-string", "", "Connection string to database")
	flag.Parse()
	if *connString == "" {
		panic("-conn-string flag must be set")
	}

	db, err := gorm.Open(postgres_driver.Open("host=localhost user=postgres password=example dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	store := postgres.Store{
		DB: db,
	}

	if err := store.AutoMigrate(); err != nil {
		panic("failed to migrate")
	}


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	h := handler.PolicyHandler{PolicyStore:    &store}
	r.POST("/policy", h.PostPolicyHandler)

	_ = r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
