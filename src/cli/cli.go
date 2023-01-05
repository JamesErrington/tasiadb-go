package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JamesErrington/tasiadb/src/parser"
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
	parser := parser.NewParser(command)
	statements := parser.Parse()

	for _, statement := range statements {
		fmt.Println(statement.Content)
	}
}

// func (lexer *Lexer) handle_error() {
// 	if err := recover(); err != nil {
// 		if le, ok := err.(LexerError); ok {
// 			lexeme := string(lexer.source[le.index])
// 			fmt.Printf("Lexer error near \"%v\": %v\n", lexeme, le.message)
// 			fmt.Printf("%s\n", lexer.source)
// 			if le.index < 15 {
// 				fmt.Printf("%*v\n", le.index+15, "^--- error here")
// 			} else {
// 				fmt.Printf("%*v\n", le.index+1, "error here ---^")
// 			}
// 		} else {
// 			panic(err)
// 		}
// 	}
// }
