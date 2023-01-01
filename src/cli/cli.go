package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JamesErrington/tasiadb/src/lexer"
)

const (
	META_CHAR    = "."
	EXIT_COMMAND = "exit"
)

func RunRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		print("tasiadb> ")

		scanner.Scan()
		input := scanner.Text()

		if strings.HasPrefix(input, META_CHAR) {
			do_meta_command(input[1:])
			continue
		}

		do_sql_command(input)
	}
}

func do_meta_command(command string) {
	switch command {
	case EXIT_COMMAND:
		os.Exit(0)
	default:
		fmt.Println("Error: unknown command or invalid arguments: ", command)
	}
}

func do_sql_command(command string) {
	lexer := lexer.NewLexer(command)
	tokens := lexer.Lex()
	fmt.Println(tokens)
}
