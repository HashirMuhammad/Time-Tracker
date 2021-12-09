package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"

	datastore "github.com/HashirMuhammad/Time-Tracker-main/Datastore"
	"github.com/HashirMuhammad/Time-Tracker-main/model"
	"github.com/dgrijalva/jwt-go"
)

type Controller struct {
	Db datastore.Datastore
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ImageUrl  string `json:"image_url"`
}

type userlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		return
	}
	err = c.Db.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	password := req.Password
	hash, err := HashPassword(password)
	if err != nil {
		err = json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}
	user := model.User{First_Name: req.FirstName, Last_Name: req.LastName, Email: req.Email, Password: hash, Image_Url: req.ImageUrl, Created_At: time.Now(), Updated_At: time.Now()}
	err = c.Db.CreateUser(user)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}
	json.NewEncoder(w).Encode("user created successfully ")

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (c Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	req := userlogin{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	dbPassword := c.Db.GetPasswordHash(req.Email)
	if dbPassword == "" {
		json.NewEncoder(w).Encode("Incorrect email or password")

		return
	}

	CheckPasswordHash(dbPassword, hash)
	if err != nil {
		json.NewEncoder(w).Encode("Incorrect email or password")

		return
	} else {
		userid := c.Db.GetUserIdByEmail(req.Email)
		id, _ := strconv.Atoi(userid)
		token, err := CreateToken(int64(id))
		if err != nil {
			json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

			return
		}
		w.Header().Set("Content-Type", token)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(token)

	}

}

func CreateToken(userid int64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
