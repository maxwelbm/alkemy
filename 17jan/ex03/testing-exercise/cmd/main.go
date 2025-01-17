package main

import (
	"fmt"
	"testdoubles/internal/application"
)

func main() {
	// env
	// ...

	// app
	// - config
	fmt.Println("http://localhost:8080")
	app := application.NewApplicationDefault(":8080", true)
	// - tear down
	// defer app.TearDown()
	// - set up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
