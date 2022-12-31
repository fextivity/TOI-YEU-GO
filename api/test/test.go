package test

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
	e.POST("/api/test/add", wrp.AddTest)
	e.POST("/api/test/delete", wrp.DeleteTest)
	e.GET("/api/test/get", wrp.GetATest)
	e.GET("/api/test/get_all", wrp.AllTests)
}