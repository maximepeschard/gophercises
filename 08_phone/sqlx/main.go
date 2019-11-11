package sqlx

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maximepeschard/gophercises/08_phone/phone"
)

const (
	schema = "exercise_08"
	table  = "phone_numbers"
)

type phoneNumber struct {
	ID     int    `db:"id"`
	Number string `db:"number"`
}

func setup(dbName string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("dbname=%s sslmode=disable", dbName)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s.%s (id SERIAL, number TEXT NOT NULL)", schema, table))
	return db, err
}

func teardown(db *sqlx.DB) error {
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	return err
}

func addPhoneNumbers(db *sqlx.DB, numbers []string) error {
	tx, err := db.Begin()
	if err != nil {
		return nil
	}

	for _, number := range numbers {
		_, err := tx.Exec(fmt.Sprintf("INSERT INTO %s.%s (number) VALUES ($1);", schema, table), number)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func updatePhoneNumber(db *sqlx.DB, id int, number string) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE %s.%s SET number = $1 WHERE id = $2", schema, table), number, id)
	return err
}

func deletePhoneNumber(db *sqlx.DB, id int) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s.%s WHERE id = $1", schema, table), id)
	return err
}

func searchPhoneNumber(db *sqlx.DB, number string) (*phoneNumber, error) {
	var pn phoneNumber
	err := db.Get(&pn, fmt.Sprintf("SELECT id, number FROM %s.%s WHERE number = $1", schema, table), number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &pn, err
}

func getPhoneNumbers(db *sqlx.DB) ([]*phoneNumber, error) {
	var pns []*phoneNumber
	err := db.Select(&pns, fmt.Sprintf("SELECT id, number FROM %s.%s", schema, table))
	return pns, err
}

func Main() {
	db, err := setup("postgres")
	check(err)
	defer db.Close()
	defer teardown(db)

	rawNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	err = addPhoneNumbers(db, rawNumbers)
	check(err)

	phoneNumbers, err := getPhoneNumbers(db)
	check(err)
	for _, pn := range phoneNumbers {
		normalized, err := phone.Normalize(pn.Number)
		check(err)

		existing, err := searchPhoneNumber(db, normalized)
		check(err)

		if existing != nil && existing.ID != pn.ID {
			deletePhoneNumber(db, pn.ID)
		} else if pn.Number != normalized {
			updatePhoneNumber(db, pn.ID, normalized)
		}
	}

	phoneNumbers, err = getPhoneNumbers(db)
	check(err)
	fmt.Println("normalized numbers:")
	for _, pn := range phoneNumbers {
		fmt.Println(pn.Number)
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
