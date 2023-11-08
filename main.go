package main

import (
	"log"
	"os"
	"os/signal"

	a "github.com/duartqx/gochatws/api/fiber"
	r "github.com/duartqx/gochatws/infrastructure/repositories/sqlite"
)

func main() {

	db, err := r.GetDbConnection("db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := a.
		GetNewAppBuilder().
		SetDb(db).
		SetPort(":8000").
		SetViewsPath("./presentation/views").
		SetStaticPath("./presentation/static").
		SetViewsBase("base").
		SetSecret("secret").
		Build()

	go func() {
		if err := app.Listen(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	app.Shutdown()
	log.Println("Shutting down")
	os.Exit(0)
}
