package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pobruno/casa360/config"
)

type FinanceCC struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type FinanceCurrency struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Symbol string   `json:"symbol"`
	Value float64   `json:"value"`
}

type FinanceInstallment struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Type           bool       `json:"type"` // false = receita, true = despesa
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	RecurrenceDays int        `json:"recurrence_days"`
	Amount         float64    `json:"amount"`
	UserID         uuid.UUID  `json:"user_id"`
	PayerGroupID   uuid.UUID  `json:"payer_group_id"`
	FinanceCCID    uuid.UUID  `json:"finance_cc_id"`
	CurrencyID     uuid.UUID  `json:"currency_id"`
}

type FinanceOccurrence struct {
	ID        uuid.UUID `json:"id"`
	FinanceID uuid.UUID `json:"finance_id"`
	Date      time.Time `json:"date"`
	Amount    float64   `json:"amount"`
	Status    bool      `json:"status"`
}

type Transaction struct {
	ID                 uuid.UUID `json:"id"`
	FinanceOccurrenceID uuid.UUID `json:"finance_occurrence_id"`
	Amount             float64    `json:"amount"`
	CreatedAt          time.Time  `json:"created_at"`
}

type FinanceWallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type OccurrenceDashboard struct {
	OccurrenceType   string     `json:"occurrence_type"`
	ID               uuid.UUID  `json:"id"`
	Date             time.Time  `json:"date"`
	Status           bool       `json:"status"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	FinanceType      *bool      `json:"finance_type,omitempty"`
	Amount           *float64   `json:"amount,omitempty"`
	CurrencySymbol   *string    `json:"currency_symbol,omitempty"`
	CurrencyValue    *float64   `json:"currency_value,omitempty"`
	AmountConverted  *float64   `json:"amount_converted,omitempty"`
	CostCenter       *string    `json:"cost_center,omitempty"`
	PayerGroup       string     `json:"payer_group"`
	ResponsibleUser  string     `json:"responsible_user"`
}

// FinanceCC methods
func (fc *FinanceCC) Create() error {
	query := `
		INSERT INTO finance_cc (id, name, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, parent_id
	`
	var parentID *uuid.UUID
	if fc.ParentID != nil {
		parentID = fc.ParentID
	}
	return config.GetDB().QueryRow(query, uuid.New(), fc.Name, parentID).
		Scan(&fc.ID, &fc.Name, &fc.ParentID)
}

func (fc *FinanceCC) Get() error {
	query := `
		SELECT id, name, parent_id
		FROM finance_cc
		WHERE id = $1
	`
	return config.GetDB().QueryRow(query, fc.ID).
		Scan(&fc.ID, &fc.Name, &fc.ParentID)
}

func ListFinanceCCs() ([]FinanceCC, error) {
	query := `
		SELECT id, name, parent_id
		FROM finance_cc
		ORDER BY name
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ccs []FinanceCC
	for rows.Next() {
		var fc FinanceCC
		if err := rows.Scan(&fc.ID, &fc.Name, &fc.ParentID); err != nil {
			return nil, err
		}
		ccs = append(ccs, fc)
	}
	return ccs, nil
}

// FinanceCurrency methods
func (fc *FinanceCurrency) Create() error {
	query := `
		INSERT INTO finance_currency (id, name, symbol, value)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, symbol, value
	`
	return config.GetDB().QueryRow(query, uuid.New(), fc.Name, fc.Symbol, fc.Value).
		Scan(&fc.ID, &fc.Name, &fc.Symbol, &fc.Value)
}

func (fc *FinanceCurrency) Get() error {
	query := `
		SELECT id, name, symbol, value
		FROM finance_currency
		WHERE id = $1
	`
	return config.GetDB().QueryRow(query, fc.ID).
		Scan(&fc.ID, &fc.Name, &fc.Symbol, &fc.Value)
}

func ListFinanceCurrencies() ([]FinanceCurrency, error) {
	query := `
		SELECT id, name, symbol, value
		FROM finance_currency
		ORDER BY name
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []FinanceCurrency
	for rows.Next() {
		var fc FinanceCurrency
		if err := rows.Scan(&fc.ID, &fc.Name, &fc.Symbol, &fc.Value); err != nil {
			return nil, err
		}
		currencies = append(currencies, fc)
	}
	return currencies, nil
}

// FinanceInstallment methods
func (fi *FinanceInstallment) Create() error {
	query := `
		INSERT INTO finance_installments (id, title, description, type, start_date, end_date, recurrence_days, amount, user_id, payer_group_id, finance_cc_id, currency_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, title, description, type, start_date, end_date, recurrence_days, amount, user_id, payer_group_id, finance_cc_id, currency_id
	`
	var endDate *time.Time
	if fi.EndDate != nil {
		endDate = fi.EndDate
	}
	return config.GetDB().QueryRow(query, uuid.New(), fi.Title, fi.Description, fi.Type, fi.StartDate, endDate, fi.RecurrenceDays, fi.Amount, fi.UserID, fi.PayerGroupID, fi.FinanceCCID, fi.CurrencyID).
		Scan(&fi.ID, &fi.Title, &fi.Description, &fi.Type, &fi.StartDate, &fi.EndDate, &fi.RecurrenceDays, &fi.Amount, &fi.UserID, &fi.PayerGroupID, &fi.FinanceCCID, &fi.CurrencyID)
}

func (fi *FinanceInstallment) Get() error {
	query := `
		SELECT id, title, description, type, start_date, end_date, recurrence_days, amount, user_id, payer_group_id, finance_cc_id, currency_id
		FROM finance_installments
		WHERE id = $1
	`
	return config.GetDB().QueryRow(query, fi.ID).
		Scan(&fi.ID, &fi.Title, &fi.Description, &fi.Type, &fi.StartDate, &fi.EndDate, &fi.RecurrenceDays, &fi.Amount, &fi.UserID, &fi.PayerGroupID, &fi.FinanceCCID, &fi.CurrencyID)
}

func (fi *FinanceInstallment) Update() error {
	query := `
		UPDATE finance_installments
		SET title = $1, description = $2, type = $3, start_date = $4, end_date = $5, recurrence_days = $6, amount = $7, user_id = $8, payer_group_id = $9, finance_cc_id = $10, currency_id = $11
		WHERE id = $12
		RETURNING id, title, description, type, start_date, end_date, recurrence_days, amount, user_id, payer_group_id, finance_cc_id, currency_id
	`
	var endDate *time.Time
	if fi.EndDate != nil {
		endDate = fi.EndDate
	}
	return config.GetDB().QueryRow(query, fi.Title, fi.Description, fi.Type, fi.StartDate, endDate, fi.RecurrenceDays, fi.Amount, fi.UserID, fi.PayerGroupID, fi.FinanceCCID, fi.CurrencyID, fi.ID).
		Scan(&fi.ID, &fi.Title, &fi.Description, &fi.Type, &fi.StartDate, &fi.EndDate, &fi.RecurrenceDays, &fi.Amount, &fi.UserID, &fi.PayerGroupID, &fi.FinanceCCID, &fi.CurrencyID)
}

func (fi *FinanceInstallment) Delete() error {
	query := `
		DELETE FROM finance_installments
		WHERE id = $1
	`
	_, err := config.GetDB().Exec(query, fi.ID)
	return err
}

func ListFinanceInstallments() ([]FinanceInstallment, error) {
	query := `
		SELECT id, title, description, type, start_date, end_date, recurrence_days, amount, user_id, payer_group_id, finance_cc_id, currency_id
		FROM finance_installments
		ORDER BY start_date DESC
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var installments []FinanceInstallment
	for rows.Next() {
		var fi FinanceInstallment
		if err := rows.Scan(&fi.ID, &fi.Title, &fi.Description, &fi.Type, &fi.StartDate, &fi.EndDate, &fi.RecurrenceDays, &fi.Amount, &fi.UserID, &fi.PayerGroupID, &fi.FinanceCCID, &fi.CurrencyID); err != nil {
			return nil, err
		}
		installments = append(installments, fi)
	}
	return installments, nil
}

// FinanceOccurrence methods
func (fo *FinanceOccurrence) Create() error {
	query := `
		INSERT INTO finance_occurrences (id, finance_id, date, amount, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, finance_id, date, amount, status
	`
	return config.GetDB().QueryRow(query, uuid.New(), fo.FinanceID, fo.Date, fo.Amount, fo.Status).
		Scan(&fo.ID, &fo.FinanceID, &fo.Date, &fo.Amount, &fo.Status)
}

func (fo *FinanceOccurrence) Update() error {
	query := `
		UPDATE finance_occurrences
		SET amount = $1, status = $2
		WHERE id = $3
		RETURNING id, finance_id, date, amount, status
	`
	return config.GetDB().QueryRow(query, fo.Amount, fo.Status, fo.ID).
		Scan(&fo.ID, &fo.FinanceID, &fo.Date, &fo.Amount, &fo.Status)
}

func ListFinanceOccurrences() ([]FinanceOccurrence, error) {
	query := `
		SELECT id, finance_id, date, amount, status
		FROM finance_occurrences
		ORDER BY date DESC
	`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []FinanceOccurrence
	for rows.Next() {
		var fo FinanceOccurrence
		if err := rows.Scan(&fo.ID, &fo.FinanceID, &fo.Date, &fo.Amount, &fo.Status); err != nil {
			return nil, err
		}
		occurrences = append(occurrences, fo)
	}
	return occurrences, nil
}

func ListFinanceOccurrencesByFinanceID(financeID uuid.UUID) ([]FinanceOccurrence, error) {
	query := `
		SELECT id, finance_id, date, amount, status
		FROM finance_occurrences
		WHERE finance_id = $1
		ORDER BY date DESC
	`
	rows, err := config.GetDB().Query(query, financeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []FinanceOccurrence
	for rows.Next() {
		var fo FinanceOccurrence
		if err := rows.Scan(&fo.ID, &fo.FinanceID, &fo.Date, &fo.Amount, &fo.Status); err != nil {
			return nil, err
		}
		occurrences = append(occurrences, fo)
	}
	return occurrences, nil
}

func (fo *FinanceOccurrence) Delete() error {
	query := `
		DELETE FROM finance_occurrences
		WHERE id = $1
	`
	_, err := config.GetDB().Exec(query, fo.ID)
	return err
}

// GenerateOccurrences gera ocorrências para uma finança baseada em sua recorrência
func (fi *FinanceInstallment) GenerateOccurrences() error {
	// Define o período de geração
	now := time.Now()
	var endDate time.Time

	if fi.EndDate != nil {
		endDate = *fi.EndDate
	} else {
		endDate = now.AddDate(1, 0, 0) // Se não houver data final, gera para 1 ano à frente
	}

	// Gera as datas de ocorrência
	nextDate := fi.StartDate
	for nextDate.Before(endDate) || nextDate.Equal(endDate) {
		// Cria a ocorrência
		occurrence := FinanceOccurrence{
			FinanceID: fi.ID,
			Date:      nextDate,
			Amount:    fi.Amount,
			Status:    false,
		}

		// Tenta criar a ocorrência (ignora se já existir devido à constraint UNIQUE)
		_ = occurrence.Create()

		// Calcular próxima data
		nextDate = nextDate.AddDate(0, 0, fi.RecurrenceDays)
	}

	return nil
}

// ListOccurrencesDashboard retorna todas as ocorrências do dashboard
func ListOccurrencesDashboard() ([]OccurrenceDashboard, error) {
	query := `SELECT * FROM occurrences_dashboard ORDER BY date DESC`
	rows, err := config.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []OccurrenceDashboard
	for rows.Next() {
		var o OccurrenceDashboard
		err := rows.Scan(
			&o.OccurrenceType,
			&o.ID,
			&o.Date,
			&o.Status,
			&o.Title,
			&o.Description,
			&o.FinanceType,
			&o.Amount,
			&o.CurrencySymbol,
			&o.CurrencyValue,
			&o.AmountConverted,
			&o.CostCenter,
			&o.PayerGroup,
			&o.ResponsibleUser,
		)
		if err != nil {
			return nil, err
		}
		occurrences = append(occurrences, o)
	}
	return occurrences, nil
}

// GetLastWalletByUserID retorna o último registro da carteira de um usuário
func GetLastWalletByUserID(userID uuid.UUID) (*FinanceWallet, error) {
	query := `
		SELECT id, user_id, amount, created_at
		FROM finance_wallets
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var wallet FinanceWallet
	err := config.GetDB().QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Amount,
		&wallet.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// ListTransactionsByOccurrenceID retorna todas as transações de uma ocorrência
func ListTransactionsByOccurrenceID(occurrenceID uuid.UUID) ([]Transaction, error) {
	query := `
		SELECT id, finance_occurrence_id, amount, created_at
		FROM transactions
		WHERE finance_occurrence_id = $1
		ORDER BY created_at DESC
	`
	rows, err := config.GetDB().Query(query, occurrenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.FinanceOccurrenceID,
			&t.Amount,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
} 