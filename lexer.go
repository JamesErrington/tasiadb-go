package main

import (
	"strings"
	"unicode/utf8"
)

type Location struct {
	line   uint
	column uint
}

type Symbol rune

const (
	SYMBOL_EOF          rune   = -1
	SYMBOL_SEMI_COLON   Symbol = ';'
	SYMBOL_COMMA        Symbol = ','
	SYMBOL_LEFT_PAREN   Symbol = '('
	SYMBOL_RIGHT_PAREN  Symbol = ')'
	SYMBOL_SINGLE_QUOTE Symbol = '\''
	SYMBOL_UNDERSCORE   Symbol = '_'
)

type Keyword string

const (
	KEYWORD_CREATE Keyword = "CREATE"
	KEYWORD_TABLE  Keyword = "TABLE"
	KEYWORD_INT    Keyword = "INT"
	KEYWORD_TEXT   Keyword = "TEXT"
)

var KEYWORDS = [...]Keyword{KEYWORD_CREATE, KEYWORD_TABLE, KEYWORD_INT, KEYWORD_TEXT}

var DATA_TYPES = [...]Keyword{KEYWORD_INT, KEYWORD_TEXT}

type TokenType uint

const (
	TOKEN_SYMBOL TokenType = iota
	TOKEN_KEYWORD
	TOKEN_IDENTIFIER
	TOKEN_INT
	TOKEN_TEXT
)

func is_whitespace(char rune) bool {
	return char == ' '
}

func is_alphabetical(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func is_digit(char rune) bool {
	return char >= '0' && char <= '9'
}

type Token struct {
	_type    TokenType
	value    string
	location Location
}

func (token *Token) is_keyword(keyword Keyword) bool {
	return token._type == TOKEN_KEYWORD && token.value == string(keyword)
}

func (token *Token) is_symbol(symbol Symbol) bool {
	return token._type == TOKEN_SYMBOL && token.value == string(symbol)
}

func (token *Token) is_data_type() bool {
	if token._type == TOKEN_KEYWORD {
		for _, data_type := range DATA_TYPES {
			if token.value == string(data_type) {
				return true
			}
		}
	}

	return false
}

type Cursor struct {
	index    int
	location Location
}

func (cursor *Cursor) increment() {
	cursor.index++
	cursor.location.column++
}

func Lex(source string) []Token {
	// @Optimization Could we use some sort of arena-style preallocation to avoid using append?
	var tokens []Token
	cursor := Cursor{}

	for {
		char := next_rune(source, cursor.index)
		if char == SYMBOL_EOF {
			break
		}

		if is_whitespace(char) {
			cursor.increment()
			continue
		}

		switch Symbol(char) {
		case SYMBOL_SEMI_COLON, SYMBOL_COMMA, SYMBOL_LEFT_PAREN, SYMBOL_RIGHT_PAREN:
			tokens = append(tokens, Token{_type: TOKEN_SYMBOL, value: string(char), location: cursor.location})
			cursor.increment()
		case SYMBOL_SINGLE_QUOTE:
			token, _ := lex_text(source, &cursor)
			tokens = append(tokens, token)
		default:
			switch {
			case is_digit(char):
				token, _ := lex_number(source, &cursor)
				tokens = append(tokens, token)
			case is_alphabetical(char):
				init_cursor := cursor
				token, ok := lex_keyword(source, &cursor)

				if !ok {
					cursor = init_cursor
					token = lex_identifier(source, &cursor)
				}

				tokens = append(tokens, token)
			}
		}
	}

	return tokens
}

func next_rune(source string, index int) rune {
	if index >= len(source) {
		return SYMBOL_EOF
	}

	char := rune(source[index])

	if char >= utf8.RuneSelf {
		panic("Source must be UTF-8 encoded!")
	}

	return char
}

func lex_text(source string, cursor *Cursor) (Token, error) {
	init_cursor := *cursor

	cursor.increment()

	for cursor.index < len(source) {
		char := next_rune(source, cursor.index)

		if Symbol(char) == SYMBOL_SINGLE_QUOTE {
			cursor.increment()
			return Token{TOKEN_TEXT, source[init_cursor.index+1 : cursor.index-1], init_cursor.location}, nil
		}

		cursor.increment()
	}

	panic("Error lexing text!")
}

func lex_number(source string, cursor *Cursor) (Token, error) {
	init_cursor := *cursor

	for cursor.index < len(source) {
		char := next_rune(source, cursor.index)

		if is_whitespace(char) {
			break
		}

		if is_digit(char) {
			cursor.increment()
			continue
		}

		panic("Error lexing number!")
	}

	return Token{TOKEN_INT, source[init_cursor.index:cursor.index], init_cursor.location}, nil
}

func lex_keyword(source string, cursor *Cursor) (Token, bool) {
	init_cursor := *cursor

	for cursor.index < len(source) {
		char := next_rune(source, cursor.index)

		if is_alphabetical(char) {
			cursor.increment()
			continue
		}

		break
	}

	value := source[init_cursor.index:cursor.index]
	for _, keyword := range KEYWORDS {
		if strings.EqualFold(value, string(keyword)) {
			return Token{TOKEN_KEYWORD, string(keyword), init_cursor.location}, true
		}
	}

	return Token{}, false
}

func lex_identifier(source string, cursor *Cursor) Token {
	init_cursor := *cursor

	for cursor.index < len(source) {
		char := next_rune(source, cursor.index)

		if is_alphabetical(char) || is_digit(char) {
			cursor.increment()
			continue
		}

		break
	}

	return Token{TOKEN_IDENTIFIER, source[init_cursor.index:cursor.index], init_cursor.location}
}
