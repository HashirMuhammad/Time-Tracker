package model

import "time"

type User struct {
	ID         int64     `db:"id,omitempty"`
	First_Name string    `db:"first_name"`
	Last_Name  string    `db:"last_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Image_Url  string    `db:"image_url"`
	Created_At time.Time `db:"created_at"`
	Updated_At time.Time `db:"updated_at"`
}
