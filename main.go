package main

import (
	"log"

	a "github.com/duartqx/gochatws/api/fiber"
	r "github.com/duartqx/gochatws/infrastructure/repositories/sqlite"
)

func main() {

	db, err := r.GetDbConnection("db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := a.Setup(db, []byte("secret"))

	log.Fatalln(app.Listen(":8000"))
}
