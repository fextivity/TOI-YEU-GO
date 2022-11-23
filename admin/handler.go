package admin

import (
	"database/sql"
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
	rows, err := db.Query(`SELECT t.id, t.name, t.start, t.end, q.id, q.content, c.id, c.content, c.is_answer
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
			idt sql.NullInt32
			t_name sql.NullString
			t_start sql.NullInt64
			t_end sql.NullInt64
			idq sql.NullInt32
			q_content sql.NullString
			idc sql.NullInt32
			c_content sql.NullString
			c_isanswer sql.NullBool
		)
		err := rows.Scan(&idt, &t_name, &t_start, &t_end, &idq, &q_content, &idc, &c_content, &c_isanswer)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !idt.Valid {
			continue;
		}
		if len(arrTest) == 0 || arrTest[len(arrTest) - 1].Id != int(idt.Int32) {
			arrTest = append(arrTest, test.Test{Id: int(idt.Int32)})
			if t_name.Valid {
				arrTest[len(arrTest) - 1].Name = t_name.String
			}
			if t_start.Valid {
				arrTest[len(arrTest) - 1].StartTime = t_start.Int64
			}
			if t_end.Valid {
				arrTest[len(arrTest) - 1].EndTime = t_end.Int64
			}
		}

		if !idq.Valid {
			continue;
		}
		arrQuestion := &arrTest[len(arrTest) - 1].Questions
		if len(*arrQuestion) == 0 || (*arrQuestion)[len(*arrQuestion) - 1].Id != int(idq.Int32) {
			*arrQuestion = append(*arrQuestion, question.Question{Id: int(idq.Int32), Content: q_content.String, Idt: int(idt.Int32)})
			
		}

		if !idc.Valid {
			continue;
		}
		arrChoice := &((*arrQuestion)[len(*arrQuestion) - 1].Choices)
		if len(*arrChoice) == 0 || (*arrChoice)[len(*arrChoice) - 1].Id != int(idc.Int32) {
			*arrChoice = append(*arrChoice, choice.Choice{Id: int(idc.Int32), Content: c_content.String, IsAnswer: c_isanswer.Bool, Idq: int(idq.Int32)})
		}
	}
	return c.JSON(http.StatusOK, arrTest)
}