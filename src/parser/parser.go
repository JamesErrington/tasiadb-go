package parser

import (
	"fmt"

	lex "github.com/JamesErrington/tasiadb/src/lexer"
)

type Parser struct {
	source        string
	source_length int
	lexer         *lex.Lexer
	start         int
	current       lex.Token
	previous      lex.Token
}

func NewParser(source string) *Parser {
	lexer := lex.NewLexer(source)
	return &Parser{source, len(source), lexer, 0, lex.Token{}, lex.Token{}}
}

func (parser *Parser) handle_error() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func (parser *Parser) Parse() []Statement {
	defer parser.handle_error()

	var statements []Statement

	parser.start = 0
	parser.advance()
	for !parser.current.IsTokenType(lex.TOKEN_EOF) {

		statement := parser.parse_statement()
		statements = append(statements, statement)

		parser.consume_token(lex.TOKEN_SEMI_COLON, "Expected semi colon at end of statement")
		parser.start = parser.lexer.StartIndex()
	}

	return statements
}

func (parser *Parser) advance() {
	parser.previous = parser.current

	for {
		token, finished := parser.lexer.NextToken()
		parser.current = token

		if finished || !token.IsTokenType(lex.TOKEN_ERROR) {
			break
		}

		panic(token)
	}
}

func (parser *Parser) match_token(token_type lex.TokenType) bool {
	if parser.current.IsTokenType(token_type) {
		parser.advance()
		return true
	}

	return false
}

func (parser *Parser) consume_token(token_type lex.TokenType, message string) {
	if parser.current.IsTokenType(token_type) {
		parser.advance()
		return
	}

	panic(message)
}

func (parser *Parser) parse_statement() Statement {
	var statement Statement

	if parser.match_token(lex.TOKEN_KEYWORD_CREATE) {
		return parser.parse_create_statement()
	}

	if parser.match_token(lex.TOKEN_KEYWORD_INSERT) {
		return parser.parse_insert_statement()
	}

	if parser.match_token(lex.TOKEN_KEYWORD_SELECT) {
		return parser.parse_select_statement()
	}

	return statement
}

func (parser *Parser) parse_create_statement() Statement {
	if parser.match_token(lex.TOKEN_KEYWORD_TABLE) {
		content := parser.parse_create_table_statement()
		return Statement{&content}
	}

	panic("Unhandled CREATE statement")
}

func (parser *Parser) parse_create_table_statement() CreateTableStatement {
	parser.consume_token(lex.TOKEN_IDENTIFIER, "Expected identifier")
	table_name_token := parser.previous

	column_names, column_types := parser.parse_column_definitions()

	return CreateTableStatement{NODE_CREATE_TABLE_STATEMENT, parser.start, table_name_token, column_names, column_types}
}

func (parser *Parser) parse_column_definitions() ([]lex.Token, []lex.Token) {
	var column_names []lex.Token
	var column_types []lex.Token

	parser.consume_token(lex.TOKEN_LEFT_PAREN, "Expected '('")
	for {
		if !parser.match_token(lex.TOKEN_IDENTIFIER) {
			panic("Expected identifier")
		}
		column_names = append(column_names, parser.previous)

		parser.advance()
		column_type_token := parser.previous
		if !column_type_token.IsDataType() {
			panic("Expected Type")
		}

		column_types = append(column_types, parser.previous)

		if parser.match_token(lex.TOKEN_COMMA) {
			continue
		}

		if parser.match_token(lex.TOKEN_RIGHT_PAREN) {
			break
		}

		panic("Expected ',' or ')")
	}

	return column_names, column_types
}

func (parser *Parser) parse_insert_statement() Statement {
	parser.consume_token(lex.TOKEN_KEYWORD_INTO, "Expected INTO")

	parser.consume_token(lex.TOKEN_IDENTIFIER, "Expected identifier")
	table_name_token := parser.previous

	var column_names []lex.Token

	if parser.match_token(lex.TOKEN_LEFT_PAREN) {
		for {
			parser.consume_token(lex.TOKEN_IDENTIFIER, "Expected identifier")
			column_names = append(column_names, parser.previous)

			if parser.match_token(lex.TOKEN_COMMA) {
				continue
			}

			if parser.match_token(lex.TOKEN_RIGHT_PAREN) {
				break
			}

			panic("Expected ',' or ')")
		}
	}

	parser.consume_token(lex.TOKEN_KEYWORD_VALUES, "Expected VALUES")

	var column_values []lex.Token

	parser.consume_token(lex.TOKEN_LEFT_PAREN, "Expected '('")
	for {

		parser.advance()
		column_value_token := parser.previous
		if !column_value_token.IsValueType() {
			panic("Expected Type")
		}

		column_values = append(column_values, column_value_token)

		if parser.match_token(lex.TOKEN_COMMA) {
			continue
		}

		if parser.match_token(lex.TOKEN_RIGHT_PAREN) {
			break
		}

		panic("Expected ',' or ')")
	}

	content := InsertStatement{NODE_INSERT_STATEMENT, parser.start, table_name_token, column_names, column_values}
	return Statement{&content}
}

func (parser *Parser) parse_select_statement() Statement {
	var columns []lex.Token
	seen_from := false

	if parser.match_token(lex.TOKEN_ASTERISK) {
		columns = append(columns, parser.previous)
	} else {
		for {
			if !parser.match_token(lex.TOKEN_IDENTIFIER) {
				panic("Expected identifier")
			}

			columns = append(columns, parser.previous)

			if parser.match_token(lex.TOKEN_COMMA) {
				continue
			}

			if parser.match_token(lex.TOKEN_KEYWORD_FROM) {
				seen_from = true
				break
			}

			panic("Expected ',' or ')")
		}
	}

	if !seen_from {
		parser.consume_token(lex.TOKEN_KEYWORD_FROM, "Expected FROM")
	}

	if !parser.match_token(lex.TOKEN_IDENTIFIER) {
		panic("Expected identifier")
	}
	table_name_token := parser.previous

	content := SelectStatement{NODE_SELECT_STATEMENT, parser.start, columns, table_name_token}
	return Statement{&content}
}
