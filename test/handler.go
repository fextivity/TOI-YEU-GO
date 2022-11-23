package test

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/database"
	"github.com/trxbach/TOI-YEU-GO/question"
)

type Test struct{
	Id int `json:"id"`
	Name string `json:"name"`
	StartTime int64 `json:"start"`
	EndTime int64 `json:"end"`
	Questions []question.Question `json:"questions"`
}

func InsertTestSql(db *database.DB, T *Test) error {
	res, err := db.Exec(`INSERT INTO tests (name, start, end) VALUES (?, ?, ?)`, T.Name, T.StartTime, T.EndTime)
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

func DeleteTestSql(db *database.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tests WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := db.Query(`SELECT id FROM questions WHERE questions.idt = ?`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var idq int
		err := rows.Scan(&idq)
		if err != nil {
			return err
		}

		err = question.DeleteQuestionSql(db, idq)
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

func (wrp *Wrapper) DeleteTest(c echo.Context) error {
	var id int
	sid := c.FormValue("id")
	id, err := strconv.Atoi(sid);
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	
	err = DeleteTestSql(wrp.db, id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "deleted")
}