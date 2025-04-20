package models

import (
	"github.com/google/uuid"
	"github.com/pobruno/casa360/config"
)

type PayerGroup struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PayerGroupMember struct {
	ID           uuid.UUID `json:"id"`
	PayerGroupID uuid.UUID `json:"payer_group_id"`
	UserID       uuid.UUID `json:"user_id"`
	Percentage   float64   `json:"percentage"`
}

func (pg *PayerGroup) Create() error {
	query := `
		INSERT INTO payer_groups (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`
	return config.GetDB().QueryRow(query, uuid.New(), pg.Name).Scan(&pg.ID, &pg.Name)
}

func (pg *PayerGroup) Get() error {
	query := `
		SELECT id, name
		FROM payer_groups
		WHERE id = $1
	`
	return config.GetDB().QueryRow(query, pg.ID).Scan(&pg.ID, &pg.Name)
}

func (pg *PayerGroup) Update() error {
	query := `
		UPDATE payer_groups
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`
	return config.GetDB().QueryRow(query, pg.Name, pg.ID).Scan(&pg.ID, &pg.Name)
}

func (pg *PayerGroup) Delete() error {
	query := `
		DELETE FROM payer_groups
		WHERE id = $1
	`
	_, err := config.GetDB().Exec(query, pg.ID)
	return err
}

func ListPayerGroups() ([]PayerGroup, error) {
	query := `
		SELECT id, name
		FROM payer_groups
		ORDER BY name
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []PayerGroup
	for rows.Next() {
		var pg PayerGroup
		if err := rows.Scan(&pg.ID, &pg.Name); err != nil {
			return nil, err
		}
		groups = append(groups, pg)
	}
	return groups, nil
}

func (pgm *PayerGroupMember) Create() error {
	query := `
		INSERT INTO payer_group_members (id, payer_group_id, user_id, percentage)
		VALUES ($1, $2, $3, $4)
		RETURNING id, payer_group_id, user_id, percentage
	`
	return config.GetDB().QueryRow(query, uuid.New(), pgm.PayerGroupID, pgm.UserID, pgm.Percentage).
		Scan(&pgm.ID, &pgm.PayerGroupID, &pgm.UserID, &pgm.Percentage)
}

func (pgm *PayerGroupMember) Delete() error {
	query := `
		DELETE FROM payer_group_members
		WHERE id = $1
	`
	_, err := config.GetDB().Exec(query, pgm.ID)
	return err
}

func ListPayerGroupMembers(payerGroupID uuid.UUID) ([]PayerGroupMember, error) {
	query := `
		SELECT id, payer_group_id, user_id, percentage
		FROM payer_group_members
		WHERE payer_group_id = $1
		ORDER BY percentage DESC
	`
	rows, err := config.GetDB().Query(query, payerGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []PayerGroupMember
	for rows.Next() {
		var pgm PayerGroupMember
		if err := rows.Scan(&pgm.ID, &pgm.PayerGroupID, &pgm.UserID, &pgm.Percentage); err != nil {
			return nil, err
		}
		members = append(members, pgm)
	}
	return members, nil
} 