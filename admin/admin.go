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
	e.GET("/api/admin/reset_database", wrp.ResetDatabase)
	e.GET("/api/admin/delete_questions", wrp.DeleteAllQuestions)
	e.GET("/api/admin/delete_answers", wrp.DeleteAllAnswers)
	e.GET("/api/admin/all", wrp.AllQuestions)
}