package admin

import (
	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/database"
)

type Wrapper struct {
	db *database.DB
}

func New(e *echo.Echo, db *database.DB) {
	wrp := &Wrapper{
		db: db,
	}
	e.POST("/api/admin/reset_database", wrp.ResetDatabase)
}
