package datastore

import (
	"github.com/HashirMuhammad/Time-Tracker-main/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Datastore interface {
	GetUserByEmail(email string) error
	CreateUser(user model.User) error
	GetPasswordHash(email string) string
	GetUserIdByEmail(email string) string
	taskQueries
}

type Database struct {
	conn *sqlx.DB
}

func NewDatabase() (Datastore, error) {
	db, err := sqlx.Open("postgres", "postgres://postgres:qwerty123@127.0.0.1:5432/tracker?sslmode=disable")
	return Database{db}, err
}

// table name should users
func (d Database) GetUserByEmail(email string) error {
	user := model.User{}

	query := `
         SELECT 
            * 
         FROM 
            users 
         WHERE 
            email = $1
            `

	return d.conn.Get(&user, query, email)
}

func (d Database) CreateUser(user model.User) error {
	query := `INSERT INTO users (
                   first_name, last_name, email, password, image_url, created_at, updated_at
                   ) 
                   VALUES 
                          ($1 , $2 , $3 , $4 , $5 , $6 , $7)`

	_, err := d.conn.Exec(query, user.First_Name, user.Last_Name, user.Email, user.Password, user.Image_Url, user.Created_At, user.Updated_At)

	return err
}

func (d Database) GetPasswordHash(email string) string {
	var password string

	query := `SELECT 
       				password 
						FROM 
						     users 
					WHERE 
					      email = $1`

	row := d.conn.QueryRow(query, email)
	row.Scan(&password)

	return password
}

func (d Database) GetUserIdByEmail(email string) string {
	var id string
	query := `SELECT 
            id
         FROM 
            users 
         WHERE 
            email = $1`

	row := d.conn.QueryRow(query, email)
	row.Scan(&id)

	return id
}
