package answer

import (
	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/db"
)

type Wrapper struct {
	db *db.DB
}

func New(e *echo.Echo, db *db.DB) {
	wrp := &Wrapper{
		db: db,
	}
	e.POST("/api/answer/add", wrp.AddAnswer)
}