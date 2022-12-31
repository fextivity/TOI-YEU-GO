package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/database"
)

func (wrp *Wrapper) ResetDatabase(c echo.Context) error {
	err := database.ResetTables(wrp.db)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Reset database\n")
}