package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pobruno/casa360/config"
	"github.com/robfig/cron"
)

type TaskInstallment struct {
	ID             uuid.UUID       `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	StartDate      time.Time      `json:"start_date"`
	RecurrenceCron string         `json:"recurrence_cron"`
	Subtasks       json.RawMessage `json:"subtasks"`
	UserID         uuid.UUID      `json:"user_id"`
	PayerGroupID   uuid.UUID      `json:"payer_group_id"`
}

type TaskOccurrence struct {
	ID           uuid.UUID       `json:"id"`
	TaskID       uuid.UUID      `json:"task_id"`
	Date         time.Time      `json:"date"`
	Status       bool           `json:"status"`
	UserID       uuid.UUID      `json:"user_id"`
	PayerGroupID uuid.UUID      `json:"payer_group_id"`
	Subtasks     json.RawMessage `json:"subtasks"`
}

// Create insere uma nova tarefa no banco de dados
func (t *TaskInstallment) Create() error {
	query := `
		INSERT INTO task_installments (id, title, description, start_date, recurrence_cron, subtasks, user_id, payer_group_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, start_date, recurrence_cron, subtasks, user_id, payer_group_id`
	return config.GetDB().QueryRow(query, uuid.New(), t.Title, t.Description, t.StartDate, t.RecurrenceCron, t.Subtasks, t.UserID, t.PayerGroupID).
		Scan(&t.ID, &t.Title, &t.Description, &t.StartDate, &t.RecurrenceCron, &t.Subtasks, &t.UserID, &t.PayerGroupID)
}

// Get busca uma tarefa pelo ID
func (t *TaskInstallment) Get() error {
	query := `
		SELECT id, title, description, start_date, recurrence_cron, subtasks, user_id, payer_group_id
		FROM task_installments
		WHERE id = $1`
	return config.GetDB().QueryRow(query, t.ID).
		Scan(&t.ID, &t.Title, &t.Description, &t.StartDate, &t.RecurrenceCron, &t.Subtasks, &t.UserID, &t.PayerGroupID)
}

// Update atualiza os dados de uma tarefa
func (t *TaskInstallment) Update() error {
	query := `
		UPDATE task_installments
		SET title = $1, description = $2, start_date = $3, recurrence_cron = $4, subtasks = $5, user_id = $6, payer_group_id = $7
		WHERE id = $8
		RETURNING id, title, description, start_date, recurrence_cron, subtasks, user_id, payer_group_id`
	return config.GetDB().QueryRow(query, t.Title, t.Description, t.StartDate, t.RecurrenceCron, t.Subtasks, t.UserID, t.PayerGroupID, t.ID).
		Scan(&t.ID, &t.Title, &t.Description, &t.StartDate, &t.RecurrenceCron, &t.Subtasks, &t.UserID, &t.PayerGroupID)
}

// Delete remove uma tarefa do banco de dados
func (t *TaskInstallment) Delete() error {
	query := `
		DELETE FROM task_installments
		WHERE id = $1`
	_, err := config.GetDB().Exec(query, t.ID)
	return err
}

// ListTasks retorna todas as tarefas
func ListTasks() ([]TaskInstallment, error) {
	query := `
		SELECT id, title, description, start_date, recurrence_cron, subtasks, user_id, payer_group_id
		FROM task_installments`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskInstallment
	for rows.Next() {
		var t TaskInstallment
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.StartDate, &t.RecurrenceCron, &t.Subtasks, &t.UserID, &t.PayerGroupID); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Create insere uma nova ocorrência de tarefa no banco de dados
func (to *TaskOccurrence) Create() error {
	query := `
		INSERT INTO task_occurrences (id, task_id, date, status, user_id, payer_group_id, subtasks)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, task_id, date, status, user_id, payer_group_id, subtasks`
	return config.GetDB().QueryRow(query, uuid.New(), to.TaskID, to.Date, to.Status, to.UserID, to.PayerGroupID, to.Subtasks).
		Scan(&to.ID, &to.TaskID, &to.Date, &to.Status, &to.UserID, &to.PayerGroupID, &to.Subtasks)
}

// Get busca uma ocorrência de tarefa pelo ID
func (to *TaskOccurrence) Get() error {
	query := `
		SELECT id, task_id, date, status, user_id, payer_group_id, subtasks
		FROM task_occurrences
		WHERE id = $1`
	return config.GetDB().QueryRow(query, to.ID).
		Scan(&to.ID, &to.TaskID, &to.Date, &to.Status, &to.UserID, &to.PayerGroupID, &to.Subtasks)
}

// Update atualiza os dados de uma ocorrência de tarefa
func (to *TaskOccurrence) Update() error {
	query := `
		UPDATE task_occurrences
		SET status = $1, user_id = $2, payer_group_id = $3, subtasks = $4
		WHERE id = $5
		RETURNING id, task_id, date, status, user_id, payer_group_id, subtasks`
	return config.GetDB().QueryRow(query, to.Status, to.UserID, to.PayerGroupID, to.Subtasks, to.ID).
		Scan(&to.ID, &to.TaskID, &to.Date, &to.Status, &to.UserID, &to.PayerGroupID, &to.Subtasks)
}

// Delete remove uma ocorrência de tarefa do banco de dados
func (to *TaskOccurrence) Delete() error {
	query := `
		DELETE FROM task_occurrences
		WHERE id = $1`
	_, err := config.GetDB().Exec(query, to.ID)
	return err
}

// ListTaskOccurrences retorna todas as ocorrências de tarefas
func ListTaskOccurrences() ([]TaskOccurrence, error) {
	query := `
		SELECT id, task_id, date, status, user_id, payer_group_id, subtasks
		FROM task_occurrences`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []TaskOccurrence
	for rows.Next() {
		var o TaskOccurrence
		if err := rows.Scan(&o.ID, &o.TaskID, &o.Date, &o.Status, &o.UserID, &o.PayerGroupID, &o.Subtasks); err != nil {
			return nil, err
		}
		occurrences = append(occurrences, o)
	}
	return occurrences, nil
}

// ListTaskOccurrencesByTaskID retorna todas as ocorrências de uma tarefa específica
func ListTaskOccurrencesByTaskID(taskID uuid.UUID) ([]TaskOccurrence, error) {
	query := `
		SELECT id, task_id, date, status, user_id, payer_group_id, subtasks
		FROM task_occurrences
		WHERE task_id = $1
		ORDER BY date`
	rows, err := config.GetDB().Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []TaskOccurrence
	for rows.Next() {
		var o TaskOccurrence
		if err := rows.Scan(&o.ID, &o.TaskID, &o.Date, &o.Status, &o.UserID, &o.PayerGroupID, &o.Subtasks); err != nil {
			return nil, err
		}
		occurrences = append(occurrences, o)
	}
	return occurrences, nil
}

// GenerateOccurrences gera ocorrências para uma tarefa baseada em seu cronograma CRON
func (t *TaskInstallment) GenerateOccurrences() error {
	// Parseia a expressão CRON
	schedule, err := cron.ParseStandard(t.RecurrenceCron)
	if err != nil {
		return fmt.Errorf("erro ao parsear expressão CRON: %v", err)
	}

	// Define o período de geração
	now := time.Now()
	endDate := now.AddDate(1, 0, 0) // Gera ocorrências para 1 ano à frente por padrão

	// Gera as datas de ocorrência
	nextTime := t.StartDate
	for nextTime.Before(endDate) {
		// Cria a ocorrência
		occurrence := TaskOccurrence{
			TaskID:       t.ID,
			Date:         nextTime,
			Status:       false,
			UserID:       t.UserID,
			PayerGroupID: t.PayerGroupID,
			Subtasks:     t.Subtasks,
		}

		// Tenta criar a ocorrência (ignora se já existir devido à constraint UNIQUE)
		_ = occurrence.Create()

		// Calcula a próxima ocorrência
		nextTime = schedule.Next(nextTime)
	}

	return nil
} 