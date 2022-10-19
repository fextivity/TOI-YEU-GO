package question

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/helper"
)

func (wrp *Wrapper) insertQuestionSql(content string) (int, error) {
	db := wrp.db
	res, err := db.Exec("INSERT INTO questions (content) VALUES (?)", content)
	idq, err2 := res.LastInsertId()
	helper.FatalOnErr(err2)
	return int(idq), err
}

func (wrp *Wrapper) AddQuestion(c echo.Context) error {
	content := c.FormValue("content")
	idq, err := wrp.insertQuestionSql(content)
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Added question, content = \"" + content + "\", idq = " + strconv.Itoa(idq) + "\n")
}