package main

import (
	"fmt"
	"github.com/dshemin/otus_highload_architect/internal/application"
	"os"
)

func main() {
	if err := application.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
