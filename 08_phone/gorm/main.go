package gorm

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/maximepeschard/gophercises/08_phone/phone"
)

const (
	schema = "exercise_08"
	table  = "phone_numbers"
)

type phoneNumber struct {
	gorm.Model
	Number string
}

func (phoneNumber) TableName() string {
	return fmt.Sprintf("%s.%s", schema, table)
}

func setup(dbName string) (*gorm.DB, error) {
	connStr := fmt.Sprintf("dbname=%s sslmode=disable", dbName)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema)).Error; err != nil {
		return nil, err
	}

	return db, nil
}

func teardown(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error
}

func Main() {
	db, err := setup("postgres")
	check(err)
	defer db.Close()
	defer teardown(db)

	check(db.AutoMigrate(&phoneNumber{}).Error)

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
	for _, number := range rawNumbers {
		db.Create(&phoneNumber{Number: number})
	}

	var phoneNumbers []*phoneNumber
	check(db.Find(&phoneNumbers).Error)
	for _, pn := range phoneNumbers {
		normalized, err := phone.Normalize(pn.Number)
		check(err)

		var existing phoneNumber
		result := db.First(&existing, "number = ?", normalized)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			check(err)
		}

		if result.RecordNotFound() && pn.Number != normalized {
			pn.Number = normalized
			db.Save(&pn)
		} else if existing.ID != pn.ID {
			db.Delete(&pn)
		}
	}

	check(db.Find(&phoneNumbers).Error)
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
