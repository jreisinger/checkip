package main

import (
	"os"

	"github.com/jreisinger/checkip/api"
	"github.com/jreisinger/checkip/cmd"
)

func main() {
	if len(os.Args[1:]) == 0 {
		api.Serve(":8000", "/")
	} else {
		cmd.Exec()
	}
}
