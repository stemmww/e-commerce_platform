package model

type User struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"` // hashed
	Name     string `db:"name" json:"name"`
}
