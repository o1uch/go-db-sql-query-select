package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectSales(t *testing.T) {

	db, err := sql.Open("sqlite", "demo.db")

	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("failed to connect to the database: ", err)
		return
	}

	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()

	client := 208
	sales, err := selectSales(db, client)

	require.NoError(t, err)
	require.NotEmpty(t, sales)

	for _, sale := range sales {
		assert.NotEmpty(t, sale.Product)
		assert.NotEmpty(t, sale.Volume)
		assert.NotEmpty(t, sale.Date)
	}
}

func TestInsertUpdateDelete(t *testing.T) {

	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	newClient := Client{
		FIO:      "TEST",
		Login:    "TEST",
		Birthday: "TEST",
		Email:    "TEST",
	}

	// insert
	id, err := insertClient(db, newClient)

	require.NoError(t, err)
	require.NotEmpty(t, id)
	newClient.ID = int(id)

	got, err := selectClient(db, id)
	require.NoError(t, err)
	require.Equal(t, newClient, got)

	// update
	newLogin := "TEST_NEW"
	err = updateClientLogin(db, newLogin, id)
	require.NoError(t, err)

	got, err = selectClient(db, id)
	require.NoError(t, err)
	require.Equal(t, newLogin, got.Login)

	// delete
	err = deleteClient(db, id)
	require.NoError(t, err)

	got, err = selectClient(db, id)
	require.Equal(t, sql.ErrNoRows, err)
	require.Empty(t, got)
}
