package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/database"
	"github.com/trxbach/TOI-YEU-GO/helper"
)

func (wrp *Wrapper) ResetDatabase(c echo.Context) error {
	err := database.ResetTables(wrp.db)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Reset database\n")
}

func (wrp *Wrapper) DeleteAllQuestions(c echo.Context) error {
	db := wrp.db
	_, err := db.Query("DELETE FROM questions")
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Deleted questions\n")
}

func (wrp *Wrapper) DeleteAllAnswers(c echo.Context) error {
	db := wrp.db
	_, err := db.Query("DELETE FROM answers")
	helper.FatalOnErr(err)
	return c.String(http.StatusOK, "Deleted answers\n")
}

func (wrp *Wrapper) AllQuestions(c echo.Context) error {
	type answer struct {
		Id int `json:"id"`
		Content string `json:"content"`
	}
	type question struct {
		Id int `json:"id"`
		Content string `json:"content"`
		Choices []answer `json:"choices"`
	}
	db := wrp.db
	rows, err := db.Query("SELECT q.id, a.id, q.content, a.content FROM questions AS q LEFT JOIN answers AS a ON q.id = a.idq ORDER BY q.id")
	helper.FatalOnErr(err)
	defer rows.Close()

	var questions []question
	setQuestion := make(map[int]struct{})
	var void struct{}
	for rows.Next() {
		var idq, ida int
		var cq, ca string
		err := rows.Scan(&idq, &ida, &cq, &ca)
		helper.FatalOnErr(err)

		if _, exist := setQuestion[idq]; !exist {
			setQuestion[idq] = void
			questions = append(questions, question{Id: idq, Content: cq, Choices: []answer{}})
		}
		choices := &questions[len(questions) - 1].Choices
		*choices = append(*choices, answer{Id: ida, Content: ca})
	}
	return c.JSON(http.StatusOK, questions)
}