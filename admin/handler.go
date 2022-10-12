package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/helper"
)

func (wrp *Wrapper) DeleteAllQuestions(c echo.Context) error {
	db := wrp.db
	_, err := db.Query("DELETE FROM questions")
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Deleted")
}