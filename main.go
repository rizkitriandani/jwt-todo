package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

func Login(c *gin.Context){
	var u User

	//check the json request
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	//check the credentials from json request with the data we provided as the mock data in db.
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	//create token
	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	//return 200 with its token
	c.JSON(http.StatusOK,token)
}


func CreateToken(userid uint64) (string, error) {
	var err error
	
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute*5).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256,atClaims)
	token,err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "",err
	}

	return token, nil
	
}


func main() {
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}
