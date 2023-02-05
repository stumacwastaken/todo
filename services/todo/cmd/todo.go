package main

import (
	"fmt"

	"github.com/stumacwastaken/todo/cmd/commands"
)

func main() {
	err := commands.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
