package Util

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GetClient attempts to create and return a valid postgres client related to the supplied DSN and Schema (directory name)
func GetClient(DSN string, Schema string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   Schema + ".",
			SingularTable: false,
		}})
	return
}

// MustConnect will attempt to connect to the database and panic if it fails
func MustConnect() *gorm.DB {
	DSN, err := GetDSN()
	if err != nil {
		panic(err)
	}
	Schema, err := GetSchema()
	if err != nil {
		panic(err)
	}
	client, err := GetClient(DSN, Schema)
	if err != nil {
		panic(err)
	}
	return client
}
