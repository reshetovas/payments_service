package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"payments_service/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type PaymentStorage struct {
	dbProperty *sql.DB
}

// function for creating object (ekzemplyar)
// output PaymentStorage which implements the interface
func NewPaymentStorage(db *sql.DB) *PaymentStorage {
	return &PaymentStorage{
		dbProperty: db,
	}
}

type PaymentStorageActions interface {
	GetPayments() ([]models.Payment, error)
	CreatePayment(payment models.Payment) (int, error)
	UpdatePayment(payment models.Payment) error
	PartialUpdatePayment(id int, updates map[string]interface{}) error
	DeletePayment(id int) error
	GetPaymentByID(id int) (models.Payment, error)
	CreateItem(item models.Item) error
	GetItemsByPaymentID(paymentID int) ([]models.Item, error)
	GetPendingPayments() ([]models.Payment, error)
}

// method get
func (s *PaymentStorage) GetPayments() ([]models.Payment, error) {
	//query to db
	log.Info().Msg("GetPayments called in storage")
	rows, err := s.dbProperty.Query("SELECT id, amount, description, created_at, state, attempts FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //to distonnect the connection to the db

	//read each database entry and fills the object
	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		err := rows.Scan(&p.ID, &p.Amount, &p.Description, &p.CreatedAt, &p.State, &p.Attempts)
		if err != nil {
			return nil, err
		}

		items, err := s.GetItemsByPaymentID(p.ID)
		if err != nil {
			return nil, err
		}
		p.Items = items

		payments = append(payments, p)
	}
	return payments, nil
}

// method create
func (s *PaymentStorage) CreatePayment(payment models.Payment) (int, error) {
	log.Info().Msg("CreatePayment called in storage")

	query := `INSERT INTO payments (amount, description, 
	created_at, currency, shop_id, address) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := s.dbProperty.Exec(query, payment.Amount, payment.Description, payment.CreatedAt, payment.Currency, payment.ShopID, payment.Address)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: payments.id" {
			return 0, errors.New("payment with this id already exists")
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("Payment created, id: %d", id)
	return int(id), nil
}

// method update put
func (s *PaymentStorage) UpdatePayment(payment models.Payment) error {
	log.Info().Msg("UpdatePayment called in storage")
	result, err := s.dbProperty.Exec(
		"UPDATE payments SET amount = ?, description = ?, created_at = ?, state = ? WHERE id = ?",
		payment.Amount, payment.Description, payment.CreatedAt, payment.State, payment.ID,
	)
	if err != nil {
		return err
	}

	// checking the update result (for cases where id not found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (s *PaymentStorage) PartialUpdatePayment(id int, updates map[string]interface{}) error {
	log.Info().Msg("PartialUpdatePayment called in storage")
	query := "Update payments SET "
	args := []interface{}{}
	i := 1

	for column, value := range updates {
		if i > 1 {
			query += ", "
		}
		if column == "id" {
			continue
		}
		query += fmt.Sprintf("%s = ?", column)
		if value == "" {
			args = append(args, "NULL")
			i++
			continue
		}
		args = append(args, value)
		i++
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.dbProperty.Exec(query, args...)
	return err
}

// method delete
func (s *PaymentStorage) DeletePayment(id int) error {
	log.Info().Msg("DeletePayment called in storage")
	result, err := s.dbProperty.Exec("DELETE FROM payments WHERE id = ?", id)
	if err != nil {
		return err
	}

	// checking the update result (for cases where id not found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}
	return nil
}

// method get by id
func (s *PaymentStorage) GetPaymentByID(id int) (models.Payment, error) {
	log.Info().Msg("GetPaymentsByID called in storage")

	//query to db
	var p models.Payment
	row := s.dbProperty.QueryRow("SELECT id, amount, description, created_at FROM payments where id = ?", id)

	err := row.Scan(&p.ID, &p.Amount, &p.Description, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Payment{}, errors.New("payments not found")
		}
		return models.Payment{}, err
	}

	return p, nil
}

func (s *PaymentStorage) CreateItem(item models.Item) error {
	log.Info().Msg("CreateItem called in storage")
	_, err := s.dbProperty.Exec(`
        INSERT INTO items (payment_id, name, price, quantity)
        VALUES (?, ?, ?, ?)`,
		item.PaymentID, item.Name, item.Price, item.Quantity)
	return err
}

func (s *PaymentStorage) GetItemsByPaymentID(paymentID int) ([]models.Item, error) {
	log.Info().Msg("GetItemsByPaymentID called in storage")
	rows, err := s.dbProperty.Query(`
        SELECT id, name, payment_id, price, quantity FROM items where payment_id = ?`,
		paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var i models.Item
		err := rows.Scan(&i.ID, &i.Name, &i.PaymentID, &i.Price, &i.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

// pyments NOT IN ('CLOSED', 'FAILED')
func (s *PaymentStorage) GetPendingPayments() ([]models.Payment, error) {
	log.Info().Msg("GetPendingPayments called in storage")
	query := `SELECT id, amount, state, attempts, created_at  FROM payments WHERE state NOT IN ('CLOSED', 'FAILED')`
	rows, err := s.dbProperty.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment

	for rows.Next() {
		var p models.Payment
		err := rows.Scan(&p.ID, &p.Amount, &p.State, &p.Attempts, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		items, err := s.GetItemsByPaymentID(p.ID)
		if err != nil {
			return nil, err
		}
		p.Items = items

		payments = append(payments, p)
	}

	return payments, nil
}

// package storage

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	_ "github.com/mattn/go-sqlite3"
// )

// type Payment struct {
// 	ID          int
// 	Amount      float32
// 	Description string
// 	CreatedAt   time.Time
// 	Currency    string
// 	ShopID      int
// 	Address     string
// }

// type Item struct {
// 	ID        int
// 	PaymentID int
// 	Name      string
// 	Price     float32
// 	Quantity  int
// }

// type SQLiteStorage struct {
// 	dbProperty *sql.DB
// }

// // function for creating object (ekzemplyar)
// // output SQLiteStorage which implements the interface
// func NewStorage(db *sql.DB) *SQLiteStorage {
// 	return &SQLiteStorage{
// 		dbProperty: db,
// 	}
// }

// type StorageActions interface {
// 	GetPayments() ([]Payment, error)
// 	CreatePayment(payment Payment) (int, error)
// 	UpdatePayment(payment Payment) error
// 	DeletePayment(id int) error
// 	GetPaymentByID(id int) (Payment, error)
// 	CreateItem(item Item) error
// }

// // method get
// func (s *SQLiteStorage) GetPayments() ([]Payment, error) {
// 	//query to db
// 	fmt.Println("GetPayments called in storage")
// 	rows, err := s.dbProperty.Query("SELECT id, amount, description, created_at FROM payments")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close() //to distonnect the connection to the db

// 	//read each database entry and fills the object
// 	var payments []Payment
// 	for rows.Next() {
// 		var p Payment
// 		err := rows.Scan(&p.ID, &p.Amount, &p.Description, &p.CreatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		payments = append(payments, p)
// 	}

// 	return payments, nil
// }

// // method create
// func (s *SQLiteStorage) CreatePayment(payment Payment) (int, error) {
// 	fmt.Println("CreatePayment called in storage")

// 	query, values := buildInsertQuery("payments", payment)
// 	fmt.Println(query, ",", values)
// 	result, err := s.dbProperty.Exec(query, values)
// 	if err != nil {
// 		if err.Error() == "UNIQUE constraint failed: payments.id" {
// 			return 0, errors.New("payment with this id already exists")
// 		}
// 		return 0, err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return int(id), nil
// }

// // method update
// func (s *SQLiteStorage) UpdatePayment(payment Payment) error {
// 	fmt.Println("UpdatePayment called in storage")
// 	result, err := s.dbProperty.Exec(
// 		"UPDATE payments SET amount = ?, description = ?, created_at = ? WHERE id = ?",
// 		payment.Amount, payment.Description, payment.CreatedAt, payment.ID,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// checking the update result (for cases where id not found)
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return errors.New("payment not found")
// 	}
// 	return nil
// }

// // method delete
// func (s *SQLiteStorage) DeletePayment(id int) error {
// 	fmt.Println("DeletePayment called in storage")
// 	result, err := s.dbProperty.Exec("DELETE FROM payments WHERE id = ?", id)
// 	if err != nil {
// 		return err
// 	}

// 	// checking the update result (for cases where id not found)
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return errors.New("payment not found")
// 	}
// 	return nil
// }

// // method get by id
// func (s *SQLiteStorage) GetPaymentByID(id int) (Payment, error) {
// 	fmt.Println("GetPaymentsByID called in storage")

// 	//query to db
// 	var p Payment
// 	row := s.dbProperty.QueryRow("SELECT id, amount, description, created_at FROM payments where id = ?", id)

// 	err := row.Scan(&p.ID, &p.Amount, &p.Description, &p.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Payment{}, errors.New("payments not found")
// 		}
// 		return Payment{}, err
// 	}

// 	return p, nil
// }

// func (s *SQLiteStorage) CreateItem(item Item) error {
// 	fmt.Println("CreateItem called in storage")
// 	_, err := s.dbProperty.Exec(`
//         INSERT INTO items (payment_id, name, price, quantity)
//         VALUES (?, ?, ?, ?)`,
// 		item.PaymentID, item.Name, item.Price, item.Quantity)
// 	return err
// }

// func buildInsertQuery(table string, payment Payment) (string, string) {
// 	var columns []string
// 	var placeholders []string
// 	var values []interface{}

// 	// Проверяем каждое поле и добавляем его в запрос, если оно не nil
// 	if payment.ID != 0 {
// 		columns = append(columns, "id")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.ID)
// 	}
// 	if payment.Amount != 0 {
// 		columns = append(columns, "amount")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.Amount)
// 	}
// 	if payment.Description != "" {
// 		columns = append(columns, "description")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.Description)
// 	}
// 	if payment.CreatedAt.IsZero() {
// 		columns = append(columns, "created_at")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.CreatedAt)
// 	}
// 	if payment.Currency != "" {
// 		columns = append(columns, "currency")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.Currency)
// 	} else {
// 		columns = append(columns, "currency")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, "USD")
// 	}
// 	if payment.ShopID != 0 {
// 		columns = append(columns, "shop_id")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.ShopID)
// 	}
// 	if payment.Address != "" {
// 		columns = append(columns, "address")
// 		placeholders = append(placeholders, "?")
// 		values = append(values, *payment.Address)
// 	}

// 	query := fmt.Sprintf(
// 		"INSERT INTO %s (%s) VALUES (%s)",
// 		table,
// 		strings.Join(columns, ", "),
// 		strings.Join(placeholders, ", "),
// 	)

// 	var valuesStringArray []string
// 	for _, v := range values {
// 		valuesStringArray = append(valuesStringArray, fmt.Sprintf("%v", v))
// 	}
// 	valuesString := strings.Join(valuesStringArray, ", ")

// 	return query, valuesString
// }
