package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	datastore "github.com/HashirMuhammad/Time-Tracker-main/Datastore"
	"github.com/HashirMuhammad/Time-Tracker-main/model"
	"github.com/dgrijalva/jwt-go"
)

type taskController struct {
	Db datastore.Datastore
}

type AccessDetails struct {
	User_Id int64
}

type taskRequest struct {
	User_Id     int64  `db:"user_id"`
	ProjectId   int64  `db:"project_id"`
	Description string `db:"description"`
}

func (c Controller) StartTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	req := taskRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	tokenAuth, err := extractTokenMetadata(r)
	if err != nil {
		json.NewEncoder(w).Encode("unauthorized")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	task := model.Task{
		ProjectId:   int64(id),
		UserId:      tokenAuth.User_Id,
		Description: req.Description,
		StartedAt:   time.Now(),
	}
	err = c.Db.CreateTask(task)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	json.NewEncoder(w).Encode("task created successfully ")
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
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

func extractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)
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

func (c Controller) StopTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	tokenAuth, err := extractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}

	err = c.Db.StopTask(int64(id), tokenAuth.User_Id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	json.NewEncoder(w).Encode("task ended successfully ")
}

func (c Controller) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	req := taskRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	tokenAuth, err := extractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}

	err = c.Db.GetTaskByProjectID(int64(projectID))
	if err != nil {
		json.NewEncoder(w).Encode("project id is not found")

		return
	}
	task := model.Task{Description: req.Description}

	err = c.Db.UpdateTask(task, int64(projectID), tokenAuth.User_Id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	json.NewEncoder(w).Encode("task updated successfully ")
}

func (c Controller) GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}
	userID := tokenAuth.User_Id
	tasks, err := c.Db.GetTasksByUserID(userID)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	for i, task := range tasks {
		diff := task.EndedAt.Minute() - task.StartedAt.Minute()
		tasks[i].TotalDuration = diff
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func (c Controller) GetLast24HrTask(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}
	userID := tokenAuth.User_Id
	from := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	tasks, err := c.Db.GetLast24HrTask(userID, from)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	for i, task := range tasks {
		diff := task.EndedAt.Minute() - task.StartedAt.Minute()
		tasks[i].TotalDuration = diff
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (c Controller) GetLastWeekTask(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}
	userID := tokenAuth.User_Id
	//f := time.Monday - time.Now().Weekday()
	//t := time.Now().AddDate(time.Now().Year(), int(time.Now().Month()), int(f))
	//ft := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	from := time.Date(time.Now().Year(), time.Now().Month(), timecheck(time.Now().Weekday()), 0, 0, 0, 0, time.Local)
	tasks, err := c.Db.GetLastWeekTask(userID, from)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))

		return
	}

	for i, task := range tasks {
		diff := task.EndedAt.Minute() - task.StartedAt.Minute()
		tasks[i].TotalDuration = diff
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func timecheck(weekday time.Weekday) int {
	var b = 0
	switch weekday {
	case time.Monday:
		{
			b = int(time.Now().Day())
			break
		}
	case time.Tuesday:
		{
			b = int(time.Now().Day()) - 1
			break
		}
	case time.Wednesday:
		{
			b = int(time.Now().Day()) - 2
			break
		}
	case time.Thursday:
		{
			b = int(time.Now().Day()) - 3
			break
		}
	case time.Friday:
		{
			b = int(time.Now().Day()) - 4
			break
		}
	case time.Saturday:
		{
			b = int(time.Now().Day()) - 5
			break
		}
	case time.Sunday:
		{
			b = int(time.Now().Day()) - 6
			break
		}
	default:
		{
			break
		}
	}

	return b
}

func (c Controller) GetLastMonthTask(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("unauthorized")

		return
	}
	userID := tokenAuth.User_Id
	from := time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Local)
	tasks, err := c.Db.GetLastMonthTask(userID, from)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("err: %s", err.Error()))
	}

	for i, task := range tasks {
		diff := task.EndedAt.Minute() - task.StartedAt.Minute()
		tasks[i].TotalDuration = diff
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
