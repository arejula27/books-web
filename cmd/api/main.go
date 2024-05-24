// Packahe main is the entry point of the application
package main

import (
	"books/internal/server"
	"fmt"
)

func main() {

	server := server.New()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
