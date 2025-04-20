package models

import (
	"github.com/google/uuid"
	"github.com/pobruno/casa360/config"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (u *User) Create() error {
	query := `
		INSERT INTO users (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`
	return config.GetDB().QueryRow(query, uuid.New(), u.Name).Scan(&u.ID, &u.Name)
}

func (u *User) Get() error {
	query := `
		SELECT id, name
		FROM users
		WHERE id = $1
	`
	return config.GetDB().QueryRow(query, u.ID).Scan(&u.ID, &u.Name)
}

func (u *User) Update() error {
	query := `
		UPDATE users
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`
	return config.GetDB().QueryRow(query, u.Name, u.ID).Scan(&u.ID, &u.Name)
}

func (u *User) Delete() error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := config.GetDB().Exec(query, u.ID)
	return err
}

func ListUsers() ([]User, error) {
	query := `
		SELECT id, name
		FROM users
		ORDER BY name
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}