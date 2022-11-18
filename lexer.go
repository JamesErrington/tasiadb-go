package main

import (
	"errors"
	"strings"
)

type Location struct {
	line   uint32
	column uint32
}

const (
	KEYWORD_SELECT = "SELECT"
	KEYWORD_FROM   = "FROM"
	KEYWORD_INSERT = "INSERT"
	KEYWORD_INTO   = "INTO"
	KEYWORD_VALUES = "VALUES"
	KEYWORD_INT    = "INT"
	KEYWORD_TEXT   = "TEXT"
)

var KEYWORDS = [...]string{KEYWORD_SELECT, KEYWORD_FROM, KEYWORD_INSERT, KEYWORD_INTO, KEYWORD_VALUES, KEYWORD_INT, KEYWORD_TEXT}

const (
	SYMBOL_SEMI         byte = ';'
	SYMBOL_ASTERISK     byte = '*'
	SYMBOL_COMMA        byte = ','
	SYMBOL_LEFT_PAREN   byte = '('
	SYMBOL_RIGHT_PAREN  byte = ')'
	SYMBOL_SINGLE_QUOTE byte = '\''
	SYMBOL_SPACE        byte = ' '
	SYMBOL_UNDERSCORE   byte = '_'
)

type TokenType uint

const (
	TOKEN_KEYWORD TokenType = iota
	TOKEN_SYMBOL
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_INT
)

type Token struct {
	_type    TokenType
	value    string
	location Location
}

type Cursor struct {
	pointer  uint32
	location Location
}

func (cursor *Cursor) increment() {
	cursor.pointer++
	cursor.location.column++
}

func (cursor *Cursor) new_line() {
	cursor.pointer++
	cursor.location.line++
	cursor.location.column = 0
}

func lex(source string) ([]*Token, error) {
	tokens := []*Token{}
	cursor := Cursor{}

	for cursor.pointer < uint32(len(source)) {
		char := source[cursor.pointer]

		switch char {
		case SYMBOL_SPACE:
			cursor.increment()
			continue
		case SYMBOL_SEMI, SYMBOL_ASTERISK, SYMBOL_COMMA, SYMBOL_LEFT_PAREN, SYMBOL_RIGHT_PAREN:
			tokens = append(tokens, &Token{_type: TOKEN_SYMBOL, value: string(char), location: cursor.location})
			cursor.increment()
		case SYMBOL_SINGLE_QUOTE:
			token, err := lex_string(source, &cursor)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		default:
			if is_digit(char) {
				token, err := lex_number(source, &cursor)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, token)
			} else if is_alpha(char) {
				start_cursor := cursor
				token, err := lex_keyword(source, &cursor)
				if err != nil {
					cursor = start_cursor
					token, err = lex_identifier(source, &cursor)
					if err != nil {
						return nil, err
					}
				}
				tokens = append(tokens, token)
			}
		}
	}

	return tokens, nil
}

func lex_string(source string, cursor *Cursor) (*Token, error) {
	start_location := cursor.location

	cursor.increment()

	var builder strings.Builder

	for cursor.pointer < uint32(len(source)) {
		char := source[cursor.pointer]

		if char == SYMBOL_SINGLE_QUOTE {
			cursor.increment()
			return &Token{_type: TOKEN_STRING, value: builder.String(), location: start_location}, nil
		}

		builder.WriteByte(char)
		cursor.increment()
	}

	return nil, errors.New("error parsing string: " + builder.String())
}

func lex_number(source string, cursor *Cursor) (*Token, error) {
	start_location := cursor.location

	var builder strings.Builder

	for cursor.pointer < uint32(len(source)) {
		char := source[cursor.pointer]

		if char == SYMBOL_SPACE {
			break
		}

		if is_digit(char) {
			builder.WriteByte(char)
			cursor.increment()
			continue
		}

		break
	}

	return &Token{_type: TOKEN_INT, value: builder.String(), location: start_location}, nil
}

func lex_keyword(source string, cursor *Cursor) (*Token, error) {
	start_location := cursor.location

	var builder strings.Builder

	for cursor.pointer < uint32(len(source)) {
		char := source[cursor.pointer]

		if char == SYMBOL_SPACE {
			break
		}

		if is_alpha(char) {
			builder.WriteByte(char)
			cursor.increment()
			continue
		}

		return nil, errors.New("error parsing keyword: " + source[start_location.column:cursor.pointer+1])
	}

	value := builder.String()
	for _, keyword := range KEYWORDS {
		if strings.EqualFold(value, keyword) {
			return &Token{_type: TOKEN_KEYWORD, value: keyword, location: start_location}, nil
		}
	}

	return nil, errors.New("unknown keyword: " + value)
}

func lex_identifier(source string, cursor *Cursor) (*Token, error) {
	start_location := cursor.location

	if !is_alpha(source[start_location.column]) {
		return nil, errors.New("error parsing identifier: must start with an alphabetical character")
	}

	var builder strings.Builder

	for cursor.pointer < uint32(len(source)) {
		char := source[cursor.pointer]

		if is_alpha(char) || is_digit(char) || char == SYMBOL_UNDERSCORE {
			builder.WriteByte(char)
			cursor.increment()
			continue
		}

		break
	}

	return &Token{_type: TOKEN_IDENTIFIER, value: builder.String(), location: start_location}, nil
}

func is_digit(char byte) bool {
	return char >= '0' && char <= '9'
}

func is_alpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}
