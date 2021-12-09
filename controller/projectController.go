package controller

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/HashirMuhammad/Time-Tracker-main/model"
	"github.com/gorilla/mux"
)

type ProjectRequest struct {
	Id          int64  `json:"id"`
	ClientName  string `json:"client_name"`
	StartedBy   string `json:"started_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c Controller) CreateProject(w http.ResponseWriter, r *http.Request) {
	req := ProjectRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		return
	}

	_, err = ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}

	project := model.Project{ClientName: req.ClientName, StartedBy: req.StartedBy, Title: req.Title, Description: req.Description}
	err = c.Db.CreateProject(project)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}
	json.NewEncoder(w).Encode("project created successfully ")
}

func (c Controller) UpdateProject(w http.ResponseWriter, r *http.Request) {
	req := ProjectRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	_, err = ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {

		return
	}

	project := model.Project{Title: req.Title, Description: req.Description}
	err = c.Db.UpdateProject(project, int64(id))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}
	json.NewEncoder(w).Encode("project updated successfully ")
}

func (c Controller) GetProjects(w http.ResponseWriter, r *http.Request) {
	_, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}

	projects, err := c.Db.GetProjects()
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if !ok {
			return nil, err
		}
		user_Id, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			User_Id: int64(user_Id),
		}, nil
	}
	return nil, err
}
