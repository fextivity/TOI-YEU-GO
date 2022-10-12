package question

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/helper"
)

func (wrp *Wrapper) InsertQuestionSql(statement string, answer string) error {
	db := wrp.db
	_, err := db.Exec("INSERT INTO questions (statement, answer) VALUES (?, ?)", statement, answer)
	return err
}

func (wrp *Wrapper) AddQuestion(c echo.Context) error {
	statement := c.FormValue("statement")
	answer := c.FormValue("answer")
	err := wrp.InsertQuestionSql(statement, answer)
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Added statement = \""+statement+"\", answer = \""+answer+"\"\n")
}

func (wrp *Wrapper) AllQuestions(c echo.Context) error {
	db := wrp.db
	rows, err := db.Query("SELECT statement, answer FROM questions")
	helper.FatalOnErr(err)
	defer rows.Close()

	var string_display strings.Builder
	for rows.Next() {
		var statement, answer string
		err := rows.Scan(&statement, &answer)
		helper.FatalOnErr(err)
		string_display.WriteString("statement = \"" + statement + "\", answer = \"" + answer + "\"\n")
	}
	return c.String(http.StatusOK, string_display.String())
}