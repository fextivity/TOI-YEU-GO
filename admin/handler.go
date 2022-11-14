package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/choice"
	"github.com/trxbach/TOI-YEU-GO/database"
	"github.com/trxbach/TOI-YEU-GO/question"
	"github.com/trxbach/TOI-YEU-GO/test"
)

func (wrp *Wrapper) ResetDatabase(c echo.Context) error {
	err := database.ResetTables(wrp.db)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Reset database\n")
}

func (wrp *Wrapper) AllTests(c echo.Context) error {
	db := wrp.db
	rows, err := db.Query(`SELECT t.id, t.name, t.start, t.end, q.id, q.content, c.id, c.content
						  FROM tests AS t
						  LEFT JOIN questions AS q ON t.id = q.idt
						  LEFT JOIN choices AS c ON q.id = c.idq
						  ORDER BY t.id, q.id, c.id`)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	arrTest := []test.Test{}
	for rows.Next() {
		var (
			idt int
			t_name string
			t_start int64
			t_end int64
			idq int
			q_content string
			idc int
			c_content string
		)
		err := rows.Scan(&idt, &t_name, &t_start, &t_end, &idq, &q_content, &idc, &c_content)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if len(arrTest) == 0 || arrTest[len(arrTest) - 1].Id != idt {
			arrTest = append(arrTest, test.Test{Id: idt, Name: t_name, StartTime: t_start, EndTime: t_end})
		}
		arrQuestion := &arrTest[len(arrTest) - 1].Questions
		if len(*arrQuestion) == 0 || (*arrQuestion)[len(*arrQuestion) - 1].Id != idq {
			*arrQuestion = append(*arrQuestion, question.Question{Id: idq, Content: q_content, Idt: idt})
		}
		arrChoice := &((*arrQuestion)[len(*arrQuestion) - 1].Choices)
		if len(*arrChoice) == 0 || (*arrChoice)[len(*arrChoice) - 1].Id != idq {
			*arrChoice = append(*arrChoice, choice.Choice{Id: idc, Content: c_content, Idq: idq})
		}
	}
	return c.JSON(http.StatusOK, arrTest)
}