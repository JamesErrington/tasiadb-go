package parser

import (
	"fmt"

	lex "github.com/JamesErrington/tasiadb/src/lexer"
)

type Parser struct {
	source        string
	source_length int
	lexer         *lex.Lexer
	current       lex.Token
	previous      lex.Token
}

func NewParser(source string) *Parser {
	lexer := lex.NewLexer(source)
	return &Parser{source, len(source), lexer, lex.Token{}, lex.Token{}}
}

func (parser *Parser) handle_error() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func (parser *Parser) Parse() []Statement {
	defer parser.handle_error()

	var statements []Statement

	for !parser.current.IsType(lex.TOKEN_EOF) {
		parser.advance()
		statement := parser.parse_statement()
		statements = append(statements, statement)
		parser.consume_token(lex.TOKEN_SEMI_COLON, "Expected semi colon at end of statement")
	}

	return statements
}

func (parser *Parser) advance() {
	parser.previous = parser.current

	for {
		token, finished := parser.lexer.NextToken()
		parser.current = token

		if finished || !token.IsType(lex.TOKEN_ERROR) {
			break
		}

		panic(token)
	}
}

func (parser *Parser) match_token(token_type lex.TokenType) bool {
	if parser.current.IsType(token_type) {
		parser.advance()
		return true
	}

	return false
}

func (parser *Parser) consume_token(token_type lex.TokenType, message string) {
	if parser.current.IsType(token_type) {
		parser.advance()
		return
	}

	panic(message)
}

func (parser *Parser) parse_statement() Statement {
	var statement Statement
	if parser.match_token(lex.TOKEN_KEYWORD_CREATE) {
		statement = parser.parse_create_statement()
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

	definitions := parser.parse_column_definitions()

	return CreateTableStatement{NODE_CREATE_TABLE_STATEMENT, 0, table_name_token, definitions}
}

func (parser *Parser) parse_column_definitions() []ColumnDefinition {
	var definitions []ColumnDefinition

	parser.consume_token(lex.TOKEN_LEFT_PAREN, "Expected '('")
	for {
		if !parser.match_token(lex.TOKEN_IDENTIFIER) {
			panic("Expected identifier")
		}
		column_name_token := parser.previous

		parser.advance()
		column_type_token := parser.previous
		if !(column_type_token.IsType(lex.TOKEN_KEYWORD_BOOLEAN) || column_type_token.IsType(lex.TOKEN_KEYWORD_NUMBER) || column_type_token.IsType(lex.TOKEN_KEYWORD_TEXT)) {
			panic("Expected Type")
		}

		definition := ColumnDefinition{column_name_token, column_type_token}
		definitions = append(definitions, definition)

		if parser.match_token(lex.TOKEN_COMMA) {
			continue
		}

		if parser.match_token(lex.TOKEN_RIGHT_PAREN) {
			break
		}

		panic("Expected ',' or ')")
	}

	return definitions
}
