package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/trxbach/TOI-YEU-GO/admin"
	"github.com/trxbach/TOI-YEU-GO/answer"
	"github.com/trxbach/TOI-YEU-GO/db"
	"github.com/trxbach/TOI-YEU-GO/general"
	"github.com/trxbach/TOI-YEU-GO/helper"
	"github.com/trxbach/TOI-YEU-GO/question"
)

func main() {
	// Load dotenv file
	err := godotenv.Load("default.env")
	helper.FatalOnErr(err)

	// Get pointer to database
	db, err := db.New(nil)
	helper.FatalOnErr(err)
	defer db.Close()

	err = general.InitTables(db)
	helper.FatalOnErr(err)

	e := echo.New()
	// e.POST("/api/auth/login")
	admin.New(e, db)
	question.New(e, db)
	answer.New(e, db)

	log.Fatal(e.Start(":1323"))
}

// go run . -env="default.env"