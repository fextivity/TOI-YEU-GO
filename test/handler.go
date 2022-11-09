package test

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/database"
	"github.com/trxbach/TOI-YEU-GO/question"
)

type Test struct{
	Id int `json:"id"`
	Content string `json:"content"`
	Questions []question.Question `json:"questions"`
}

func InsertTestSql(db *database.DB, T *Test) error {
	res, err := db.Exec("INSERT INTO tests (content) VALUES (?)", T.Content)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	T.Id = int(id)
	for i := range T.Questions {
		Q := &T.Questions[i]
		Q.Idt = T.Id
		err := question.InsertQuestionSql(db, Q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wrp *Wrapper) AddTest(c echo.Context) error {
	var T Test
	if err := c.Bind(&T); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := InsertTestSql(wrp.db, &T); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, T)
}