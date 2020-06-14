package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/middlewares"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (s *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	s.DB, err = gorm.Open("mysql", DBURL)
	if err != nil {
		fmt.Println("Cannot connect to mysql database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("connected to mysql")
	}

	s.DB.Debug().AutoMigrate(
		&models.Post{},
		&models.User{},
		&models.Comment{},
		&models.Like{},
		&models.ResetPassword{},
	)
	s.Router = gin.Default()
	s.Router.Use(middlewares.CORSMiddleware())
	s.InitializeRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Println("listening to port 8070")
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
