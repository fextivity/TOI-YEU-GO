package answer

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/helper"
)

func (wrp *Wrapper) insertAnswerSql(content string, idq int) error {
	db := wrp.db
	_, err := db.Exec("INSERT INTO answers (content, idq) VALUES (?, ?)", content, idq)
	return err
}

func (wrp *Wrapper) AddAnswer(c echo.Context) error {
	content := c.FormValue("content")
	idq, err := strconv.Atoi(c.FormValue("idq"))
	helper.FatalOnErr(err)
	err = wrp.insertAnswerSql(content, idq)
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Added answer, content = \"" + content + "\", idq = " + strconv.Itoa(idq) + "\n")
}

