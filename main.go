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
	tokens, err := lex(command)
	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
