package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("tasiadb> ")

		scanner.Scan()
		input := scanner.Text()

		if strings.HasPrefix(input, ".") {
			do_meta_command(input)
			continue
		}

		do_sql_command(input)
	}
}

func do_meta_command(command string) {
	switch command {
	case ".exit":
		os.Exit(0)
	default:
		fmt.Println("Error: unknown command or invalid arguments: ", command)
	}
}

func do_sql_command(command string) {
	tokens := Lex(command)

	for _, token := range tokens {
		fmt.Println(token)
	}

	ast := Parse(tokens)

	for _, statement := range ast.statements {
		fmt.Println(*statement.create_table_statement)
	}
}
