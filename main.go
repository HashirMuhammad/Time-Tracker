package main

import (
	datastore "github.com/HashirMuhammad/Time-Tracker-main/Datastore"
	"net/http"

	"github.com/HashirMuhammad/Time-Tracker-main/controller"
	"github.com/gorilla/mux"
)

func main() {
	db, err := datastore.NewDatabase()
	if err != nil {
		panic(err)
	}

	ctrl := controller.Controller{Db: db}

	router := mux.NewRouter()
	router.HandleFunc("/signup", ctrl.CreateUser).Methods("POST")
	router.HandleFunc("/login", ctrl.LoginUser).Methods("POST")
	router.HandleFunc("/start/{id}", ctrl.StartTask).Methods("POST")
	router.HandleFunc("/stop/{id}", ctrl.StopTask).Methods("PATCH")
	router.HandleFunc("/task/{id}", ctrl.UpdateTask).Methods("PATCH")
	router.HandleFunc("/tasks", ctrl.GetTasksByUserID).Methods("Get")
	router.HandleFunc("/tasks/24hrs", ctrl.GetLast24HrTask).Methods("Get")
	router.HandleFunc("/tasks/week", ctrl.GetLastWeekTask).Methods("Get")
	router.HandleFunc("/tasks/month", ctrl.GetLastMonthTask).Methods("Get")
	router.HandleFunc("/project", ctrl.CreateProject).Methods("POST")
	router.HandleFunc("/updproject/{id}", ctrl.UpdateProject).Methods("PATCH")
	router.HandleFunc("/projects", ctrl.GetProjects).Methods("Get")

	err = http.ListenAndServe(":8000", router)

	//handler := cors.New(cors.Options{
	//	AllowedMethods: []string{"GET", "POST", "PATCH"},
	//}).Handler(router)
	//
	//http.ListenAndServe(":8000", handler)

	if err != nil {
		return
	}
}
