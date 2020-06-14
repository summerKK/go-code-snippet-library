package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
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

	s.DB.Debug().AutoMigrate(&models.Post{}, &models.User{})
	s.Router = mux.NewRouter()
}

func (s *Server) Run(addr string) {
	fmt.Println("listening to port 8070")
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
