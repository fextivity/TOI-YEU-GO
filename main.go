package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/trxbach/TOI-YEU-GO/admin"
	"github.com/trxbach/TOI-YEU-GO/database"
	"github.com/trxbach/TOI-YEU-GO/helper"
	"github.com/trxbach/TOI-YEU-GO/test"
)

func main() {
	// Load dotenv file
	// godotenv.Load("default.env")

	// Get pointer to database
	db, err := database.New(nil)
	helper.FatalOnErr(err)
	defer db.Close()

	err = database.CreateTables(db)
	helper.FatalOnErr(err)

	e := echo.New()
	admin.New(e, db)
	// answer.New(e, db)
	// question.New(e, db)
	test.New(e, db)

	log.Fatal(e.Start(":1323"))
}