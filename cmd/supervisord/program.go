package supervisord

import (
	"os"
	"strings"
)

func Program(args ...string) string {
	program := os.Args[0]
	if strings.HasPrefix(program, "/var/") {
		program = "go run ."
	}

	if len(args) > 0 {
		program += " " + strings.Join(args, " ")
	}

	return program
}
