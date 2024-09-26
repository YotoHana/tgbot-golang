package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Title struct {
	Title string
}

func ConnectToDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "internal/database/todo.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to SQLite database!")
	return db, nil
}

func CreateTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS todo (
		chatID BIGINT,
		title TEXT,
		description TEXT
	);`
	_, err := db.Exec(query)
	fmt.Println(err)
}

func InsertData (db *sql.DB, chatID int64, title string) error {
	query := "INSERT INTO todo (chatID, title) VALUES ($1, $2);"
	_, err := db.Exec(query, chatID, title)
	if err != nil {
		return err
	}
	return nil
}

func QueryData (db *sql.DB, chatID int64) ([]Title, error) {
	query := "SELECT title FROM todo WHERE chatID = $1;"
	rows, err := db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var titles []Title
	for rows.Next() {
		var title Title
		err = rows.Scan(&title.Title)
		if err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	return titles, nil
}

func DeleteData (db *sql.DB, chatID int64, title string) error {
	query := "DELETE FROM todo WHERE title = $1 AND chatID = $2"
	_, err := db.Exec(query, title, chatID)
	if err != nil {
		return err
	}
	return nil
}