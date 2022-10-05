package main

import (
	"database/sql"
	_ "fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"net/http"

	"github.com/labstack/echo/v4"
)

type DbWrapper struct {
	db *sql.DB
}

func FatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitTables(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS questions (id INTEGER PRIMARY KEY AUTO_INCREMENT, statement TEXT, answer TEXT)")
	FatalOnErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY AUTO_INCREMENT, account TEXT, password TEXT)")
	FatalOnErr(err)
	return nil
}

func (dbw *DbWrapper) AddQuestionSql(statement string, answer string) error {
	db := dbw.db
	_, err := db.Exec("INSERT INTO questions (statement, answer) VALUES (?, ?)", statement, answer)
	return err
}

func (dbw *DbWrapper) AddQuestion(c echo.Context) error {
	statement := c.FormValue("statement")
	answer := c.FormValue("answer")
	err := dbw.AddQuestionSql(statement, answer)
	FatalOnErr(err)
	return c.String(http.StatusOK, "Added statement = \""+statement+"\", answer = \""+answer+"\"\n")
}

func (dbw *DbWrapper) AllQuestions(c echo.Context) error {
	db := dbw.db
	rows, err := db.Query("SELECT statement, answer FROM questions")
	FatalOnErr(err)
	defer rows.Close()

	var string_display strings.Builder
	for rows.Next() {
		var statement, answer string
		err := rows.Scan(&statement, &answer)
		FatalOnErr(err)
		string_display.WriteString("statement = \"" + statement + "\", answer = \"" + answer + "\"\n")
	}
	return c.String(http.StatusOK, string_display.String())
}

func (dbw *DbWrapper) DeleteAllQuestions(c echo.Context) error {
	db := dbw.db
	_, err := db.Query("DELETE FROM questions")
	FatalOnErr(err)
	return c.String(http.StatusOK, "Deleted")
}

func main() {
	var dbw DbWrapper
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/db")
	FatalOnErr(err)
	defer db.Close()
	dbw.db = db

	err = InitTables(db)
	FatalOnErr(err)

	e := echo.New()
	e.POST("/api/question/add", dbw.AddQuestion)
	e.GET("/api/question/all", dbw.AllQuestions)
	e.GET("/api/question/deleteall", dbw.DeleteAllQuestions)

	log.Fatal(e.Start(":1323"))
}
