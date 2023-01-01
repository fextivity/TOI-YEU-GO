package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/trxbach/TOI-YEU-GO/api/question"
	"github.com/trxbach/TOI-YEU-GO/choice"
	"github.com/trxbach/TOI-YEU-GO/database"
)

type Test struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartTime int64  `json:"start"`
	EndTime   int64  `json:"end"`
	// sometimes we intentionally leave this out, for example when we only need the test's info
	Questions []question.Question `json:"questions,omitempty"`
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
	client := http.Client{}
	revalidateURL := fmt.Sprintf("%s/api/internal/revalidate_exam?secret=%s&examId=%d", os.Getenv("FRONTEND_ADDR"), os.Getenv("FRONTEND_API_TOKEN"), T.Id)
	revalidateReq, err := http.NewRequest("GET", revalidateURL, nil)
	if err != nil {
		return err
	}
	revalidateRes, err := client.Do(revalidateReq)
	if err != nil {
		return err
	}
	if revalidateRes.StatusCode == 200 {
		fmt.Printf("Added test and revalidated the frontend successfully.\n")
	} else {
		fmt.Printf("An error occurred while revalidating the frontend. Status code: %d. Please check the frontend's log.\n", revalidateRes.StatusCode)
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
	id, err := strconv.Atoi(sid)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = DeleteTestSql(wrp.db, id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "deleted")
}

func (wrp *Wrapper) GetATest(c echo.Context) error {
	db := wrp.db
	id := c.QueryParam("id")
	omit_questions, err := strconv.ParseBool(c.QueryParam("omit_questions"))
	if err != nil {
		omit_questions = false
	}
	var (
		idq        sql.NullInt32
		q_content  sql.NullString
		idc        sql.NullInt32
		c_content  sql.NullString
		c_isanswer sql.NullBool
	)
	var test Test

	if !omit_questions {
		rows, err := db.Query(`SELECT t.id, t.name, t.start, t.end, q.id, q.content, c.id, c.content, c.is_answer
						FROM tests AS t 
						LEFT JOIN questions AS q ON t.id = q.idt
						LEFT JOIN choices AS c ON q.id = c.idq
						WHERE t.id = ?
						ORDER BY t.id, q.id, c.id`, id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for rows.Next() {
			err := rows.Scan(&test.Id, &test.Name, &test.StartTime, &test.EndTime, &idq, &q_content, &idc, &c_content, &c_isanswer)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			if !idq.Valid {
				continue
			}
			arrQuestion := &test.Questions
			if len(*arrQuestion) == 0 || (*arrQuestion)[len(*arrQuestion)-1].Id != int(idq.Int32) {
				*arrQuestion = append(*arrQuestion, question.Question{Id: int(idq.Int32), Content: q_content.String, Idt: test.Id})
			}
			if !idc.Valid {
				continue
			}
			arrChoice := &((*arrQuestion)[len(*arrQuestion)-1].Choices)
			if len(*arrChoice) == 0 || (*arrChoice)[len(*arrChoice)-1].Id != int(idc.Int32) {
				*arrChoice = append(*arrChoice, choice.Choice{Id: int(idc.Int32), Content: c_content.String, IsAnswer: c_isanswer.Bool, Idq: int(idq.Int32)})
			}
		}
		return c.JSON(http.StatusOK, test)
	}
	row := db.QueryRow(`SELECT t.id, t.name, t.start, t.end
							FROM tests AS t 
							WHERE t.id = ?`, id)
	switch err := row.Scan(&test.Id, &test.Name, &test.StartTime, &test.EndTime); err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, "Test not found.")
	case nil:
		break
	default:
		return c.JSON(http.StatusInternalServerError, "Internal server error.")
	}
	return c.JSON(http.StatusOK, test)
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

	arrTest := []Test{}
	for rows.Next() {
		var (
			idt        sql.NullInt32
			t_name     sql.NullString
			t_start    sql.NullInt64
			t_end      sql.NullInt64
			idq        sql.NullInt32
			q_content  sql.NullString
			idc        sql.NullInt32
			c_content  sql.NullString
			c_isanswer sql.NullBool
		)
		err := rows.Scan(&idt, &t_name, &t_start, &t_end, &idq, &q_content, &idc, &c_content, &c_isanswer)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !idt.Valid {
			continue
		}
		if len(arrTest) == 0 || arrTest[len(arrTest)-1].Id != int(idt.Int32) {
			arrTest = append(arrTest, Test{Id: int(idt.Int32)})
			if t_name.Valid {
				arrTest[len(arrTest)-1].Name = t_name.String
			}
			if t_start.Valid {
				arrTest[len(arrTest)-1].StartTime = t_start.Int64
			}
			if t_end.Valid {
				arrTest[len(arrTest)-1].EndTime = t_end.Int64
			}
		}

		if !idq.Valid {
			continue
		}

		arrQuestion := &arrTest[len(arrTest)-1].Questions
		if len(*arrQuestion) == 0 || (*arrQuestion)[len(*arrQuestion)-1].Id != int(idq.Int32) {
			*arrQuestion = append(*arrQuestion, question.Question{Id: int(idq.Int32), Content: q_content.String, Idt: int(idt.Int32)})
		}

		if !idc.Valid {
			continue
		}

		arrChoice := &((*arrQuestion)[len(*arrQuestion)-1].Choices)
		if len(*arrChoice) == 0 || (*arrChoice)[len(*arrChoice)-1].Id != int(idc.Int32) {
			*arrChoice = append(*arrChoice, choice.Choice{Id: int(idc.Int32), Content: c_content.String, IsAnswer: c_isanswer.Bool, Idq: int(idq.Int32)})
		}
	}
	return c.JSON(http.StatusOK, arrTest)
}
