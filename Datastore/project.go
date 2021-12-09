package datastore

import "github.com/HashirMuhammad/Time-Tracker-main/model"

type project struct {
	Id          int64  `db:"id"`
	ClientName  string `db:"client_name"`
	StartedBy   string `db:"started_by"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

type projectQueries interface {
	CreateProject(project model.Project) error
	UpdateProject(project model.Project, id int64) error
	GetProjects() ([]model.Project, error)
}

func (d Database) CreateProject(project model.Project) error {
	query := `INSERT INTO projects (
                   client_name, started_by, title, description
                   ) 
                   VALUES 
                          ($1 , $2 , $3 , $4)`

	_, err := d.conn.Exec(query, project.ClientName, project.StartedBy, project.Title, project.Description)

	return err
}

func (d Database) UpdateProject(project model.Project, id int64) error {
	query := `UPDATE  
    				projects 
			  SET 
				    title = $1, description = $2
			  WHERE 
					id = $3`

	_, err := d.conn.Exec(query, project.Title, project.Description, id)

	return err
}

func (d Database) GetProjects() ([]model.Project, error) {
	var projects []project
	query := `SELECT 
				id, client_name, started_by, title, description
				FROM 
				     projects`

	err := d.conn.Select(&projects, query)

	var resp []model.Project

	for i, _ := range projects {
		resp = append(resp, model.Project{
			Id:          projects[i].Id,
			ClientName:  projects[i].ClientName,
			StartedBy:   projects[i].StartedBy,
			Title:       projects[i].Title,
			Description: projects[i].Description,
		})
	}

	return resp, err

}
