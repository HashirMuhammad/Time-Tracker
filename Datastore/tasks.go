package datastore

import (
	"github.com/HashirMuhammad/Time-Tracker-main/model"
	"time"
)

type task struct {
	Id          int64     `db:"id"`
	ProjectId   int64     `db:"project_id"`
	UserId      int64     `db:"user_id"`
	Description string    `db:"description"`
	StartedAt   time.Time `db:"started_at"`
	EndedAt     time.Time `db:"ended_at"`
}

type taskQueries interface {
	CreateTask(task model.Task) error
	StopTask(taskID, userID int64) error
	GetTaskByProjectID(projectID int64) error
	UpdateTask(task model.Task, projectID, userID int64) error
	GetTasksByUserID(userID int64) ([]model.Task, error)
	GetLast24HrTask(userID int64, from time.Time) ([]model.Task, error)
	GetLastWeekTask(userID int64, from time.Time) ([]model.Task, error)
	GetLastMonthTask(userID int64, from time.Time) ([]model.Task, error)
}

func (d Database) CreateTask(task model.Task) error {
	query := `INSERT INTO tasks 
                ( project_id, user_id, description, started_at ) 
              VALUES 
                ($1 , $2 , $3, $4)`

	_, err := d.conn.Exec(query, task.ProjectId, task.UserId, task.Description, task.StartedAt)

	return err
}

func (d Database) StopTask(taskID, userID int64) error {
	query := `UPDATE  
    				tasks 
			  SET 
				    ended_at = $1 
			  WHERE 
					id = $2 AND user_id = $3`

	_, err := d.conn.Exec(query, time.Now(), taskID, userID)

	return err
}

func (d Database) GetTaskByProjectID(projectID int64) error {
	query := `SELECT 
				project_id
			FROM
				tasks
			WHERE
				project_id = $1`

	err := d.conn.Get(&task{}, query, projectID)

	return err
}

func (d Database) UpdateTask(task model.Task, projectID, userID int64) error {
	query := `UPDATE  
    				tasks 
			  SET 
				    description = $1 
			  WHERE 
					project_id = $2 AND user_id = $3`

	_, err := d.conn.Exec(query, task.Description, projectID, userID)

	return err
}

func (d Database) GetTasksByUserID(userID int64) ([]model.Task, error) {
	var tasks []task
	query := `SELECT 
				id, project_id, description, started_at, ended_at
				FROM 
				     tasks
				WHERE 
				      user_id = $1`

	// it return []tasks if there is not an error
	// if there is a error it will return error and []tasks will be empty.
	err := d.conn.Select(&tasks, query, userID)

	var resp []model.Task

	for i, _ := range tasks {
		resp = append(resp, model.Task{
			Id:          tasks[i].Id,
			ProjectId:   tasks[i].ProjectId,
			Description: tasks[i].Description,
			StartedAt:   tasks[i].StartedAt,
			EndedAt:     tasks[i].EndedAt,
		})
	}

	return resp, err
}

func (d Database) GetLast24HrTask(userID int64, from time.Time) ([]model.Task, error) {
	var tasks []task
	query := `SELECT 
				id, project_id, description, started_at, ended_at
				FROM 
				     tasks
				WHERE 
				      user_id = $1 AND started_at >= $2`

	err := d.conn.Select(&tasks, query, userID, from)

	var resp []model.Task

	for i, _ := range tasks {
		resp = append(resp, model.Task{
			Id:          tasks[i].Id,
			ProjectId:   tasks[i].ProjectId,
			Description: tasks[i].Description,
			StartedAt:   tasks[i].StartedAt,
			EndedAt:     tasks[i].EndedAt,
		})
	}

	return resp, err
}

func (d Database) GetLastWeekTask(userID int64, from time.Time) ([]model.Task, error) {
	var tasks []task
	query := `SELECT 
				id, project_id, description, started_at, ended_at
				FROM 
				     tasks
				WHERE 
				      user_id = $1 AND started_at >= $2`

	err := d.conn.Select(&tasks, query, userID, from)

	var resp []model.Task

	for i, _ := range tasks {
		resp = append(resp, model.Task{
			Id:          tasks[i].Id,
			ProjectId:   tasks[i].ProjectId,
			Description: tasks[i].Description,
			StartedAt:   tasks[i].StartedAt,
			EndedAt:     tasks[i].EndedAt,
		})
	}

	return resp, err
}

func (d Database) GetLastMonthTask(userID int64, from time.Time) ([]model.Task, error) {
	var tasks []task
	query := `SELECT 
				id, project_id, description, started_at, ended_at
				FROM 
				     tasks
				WHERE 
				      user_id = $1 AND started_at >= &2`

	// started_at >= NOW() - '1 month'::INTERVAL` IS to take records for 1 month
	//i.e: today is 6 dec it will take records from 6 nov
	err := d.conn.Select(&tasks, query, userID, from)

	var resp []model.Task

	for i, _ := range tasks {
		resp = append(resp, model.Task{
			Id:          tasks[i].Id,
			ProjectId:   tasks[i].ProjectId,
			Description: tasks[i].Description,
			StartedAt:   tasks[i].StartedAt,
			EndedAt:     tasks[i].EndedAt,
		})
	}

	return resp, err
}
