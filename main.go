package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Sale struct {
	Product int
	Volume  int
	Date    string
}

type Client struct {
	ID       int
	FIO      string
	Login    string
	Birthday string // строка в формате YYYYMMDD
	Email    string
}

func (s Sale) String() string {
	return fmt.Sprintf("Product: %d Volume: %d Date:%s", s.Product, s.Volume, s.Date)
}

func (c Client) String() string {
	return fmt.Sprintf("ID: %d FIO: %s Login: %s Birthday: %s Email: %s",
		c.ID, c.FIO, c.Login, c.Birthday, c.Email)
}

func selectSales(db *sql.DB, id int) ([]Sale, error) {
	var sales []Sale

	rows, err := db.Query("select product, volume, date from sales where client = :client_id", sql.Named("client_id", id))

	if err != nil {
		return nil, fmt.Errorf("request execution error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var sale Sale
		rows.Scan(&sale.Product, &sale.Volume, &sale.Date)

		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("errors in reading results: %v", err)
	}

	return sales, nil
}

func insertClient(db *sql.DB, client Client) (int, error) {

	tx, err := db.Begin()

	if err != nil {
		return 0, fmt.Errorf("error opening transaction: %w", err)
	}

	defer tx.Rollback()

	res, err := tx.Exec("insert into clients(fio, login, birthday, email) values(:fio, :login, :birthday, :email)",
		sql.Named("fio", client.FIO),
		sql.Named("login", client.Login),
		sql.Named("birthday", client.Birthday),
		sql.Named("email", client.Email))

	if err != nil {
		return 0, fmt.Errorf("error adding data: %w", err)
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("error receiving lastId: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return 0, fmt.Errorf("commit command error: %w", err)
	}
	return int(lastID), nil

}

func updateClientLogin(db *sql.DB, newLogin string, id int) error {

	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("error opening transaction: %w", err)
	}

	defer tx.Rollback()

	res, err := tx.Exec("update clients set login = :login where id = :id",
		sql.Named("login", newLogin),
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("data update error: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting result: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit command error: %w", err)
	}
	return nil
}

func deleteClient(db *sql.DB, id int) error {
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("error opening transaction: %w", err)
	}

	defer tx.Rollback()

	res, err := tx.Exec("delete from clients where id = :id",
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("data deletion error: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting result: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit command error: %w", err)
	}
	return nil

}

func selectClient(db *sql.DB, id int) (Client, error) {
	client := Client{}

	row := db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&client.ID, &client.FIO, &client.Login, &client.Birthday, &client.Email)

	return client, err
}

func main() {
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

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)

	Newclient := Client{
		FIO:      "John Doe",
		Login:    "JDFPerson",
		Birthday: "19700101",
		Email:    "ThefirstpersonJD@gmail.com",
	}

	id, err := insertClient(db, Newclient)

	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client)

	newLogin := "AgentSmith@gmail.com"
	err = updateClientLogin(db, newLogin, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err = selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client)

	err = deleteClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = selectClient(db, id)
	if err != nil {
		fmt.Println(err)
		return
	}
}
