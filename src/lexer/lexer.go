package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Symbol rune

const (
	SYMBOL_EOF          Symbol = -1
	SYMBOL_SEMI_COLON   Symbol = ';'
	SYMBOL_COMMA        Symbol = ','
	SYMBOL_LEFT_PAREN   Symbol = '('
	SYMBOL_RIGHT_PAREN  Symbol = ')'
	SYMBOL_SINGLE_QUOTE Symbol = '\''
	SYMBOL_UNDERSCORE   Symbol = '_'
	SYMBOL_ASTERISK     Symbol = '*'
	SYMBOL_DOT          Symbol = '.'
)

type Keyword string

const (
	KEYWORD_CREATE  Keyword = "CREATE"
	KEYWORD_TABLE   Keyword = "TABLE"
	KEYWORD_NUMBER  Keyword = "NUMBER"
	KEYWORD_TEXT    Keyword = "TEXT"
	KEYWORD_BOOLEAN Keyword = "BOOLEAN"
	KEYWORD_INSERT  Keyword = "INSERT"
	KEYWORD_INTO    Keyword = "INTO"
	KEYWORD_VALUES  Keyword = "VALUES"
	KEYWORD_TRUE    Keyword = "TRUE"
	KEYWORD_FALSE   Keyword = "FALSE"
	KEYWORD_SELECT  Keyword = "SELECT"
	KEYWORD_FROM    Keyword = "FROM"
)

type TokenType uint8

const (
	TOKEN_SYMBOL TokenType = iota
	TOKEN_KEYWORD
	TOKEN_IDENTIFIER
	TOKEN_NUMBER_LITERAL
	TOKEN_TEXT_LITERAL
	TOKEN_EOF
)

func is_whitespace(char rune) bool {
	switch char {
	case ' ', '\n', '\t':
		return true
	default:
		return false
	}
}

func is_eof(char rune) bool {
	return char == rune(SYMBOL_EOF)
}

func is_alphabetical(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func is_digit(char rune) bool {
	return char >= '0' && char <= '9'
}

func is_alphanumeric(char rune) bool {
	return is_alphabetical(char) || is_digit(char) || char == '_'
}

func to_upper_rune(char rune) rune {
	if char >= 'a' && char <= 'z' {
		return char - 32
	}

	return char
}

type Token struct {
	_type  TokenType
	value  string
	offset int
}

func (token Token) IsType(token_type TokenType) bool {
	return token._type == token_type
}

type Lexer struct {
	source        string
	source_length int
	index         int
	start         int
}

type LexerError struct {
	message string
	index   int
}

func (lexer *Lexer) throw_error(message string) {
	panic(LexerError{message, lexer.index})
}

func (lexer *Lexer) handle_error() {
	if err := recover(); err != nil {
		if le, ok := err.(LexerError); ok {
			lexeme := string(lexer.source[le.index])
			fmt.Printf("Lexer error near \"%v\": %v\n", lexeme, le.message)
			fmt.Printf("%s\n", lexer.source)
			if le.index < 15 {
				fmt.Printf("%*v\n", le.index+15, "^--- error here")
			} else {
				fmt.Printf("%*v\n", le.index+1, "error here ---^")
			}
		} else {
			panic(err)
		}
	}
}

func NewLexer(source string) *Lexer {
	return &Lexer{source, len(source), -1, 0}
}

func (lexer *Lexer) Lex() []Token {
	defer lexer.handle_error()

	var tokens []Token

	for lexer.index = -1; ; {
		char := lexer.next_rune()

		if is_eof(char) {
			tokens = append(tokens, Token{TOKEN_EOF, "", lexer.index})
			return tokens
		}

		if is_whitespace(char) {
			continue
		}

		switch Symbol(char) {
		case
			SYMBOL_SEMI_COLON,
			SYMBOL_COMMA,
			SYMBOL_LEFT_PAREN,
			SYMBOL_RIGHT_PAREN,
			SYMBOL_ASTERISK:
			tokens = append(tokens, Token{TOKEN_SYMBOL, string(char), lexer.index})
		case SYMBOL_SINGLE_QUOTE:
			token := lexer.lex_text()
			tokens = append(tokens, token)
		default:
			switch {
			case is_digit(char) || char == rune(SYMBOL_DOT):
				token := lexer.lex_number()
				tokens = append(tokens, token)
			case is_alphabetical(char):
				token := lexer.lex_keyword_or_identifier()
				tokens = append(tokens, token)
			default:
				lexer.throw_error("Unexpected token")
			}
		}
	}
}

func (lexer *Lexer) next_rune() rune {
	lexer.index += 1

	if lexer.index >= lexer.source_length {
		return rune(SYMBOL_EOF)
	}

	char := rune(lexer.source[lexer.index])

	if char < 0 || char >= utf8.RuneSelf {
		lexer.throw_error("Encountered non UTF-8 character")
	}

	return char
}

func (lexer *Lexer) lex_number() Token {
	lexer.start = lexer.index
	lexer.index -= 1
	seen_dot := false

	for lexer.index < lexer.source_length {
		char := lexer.next_rune()

		if is_digit(char) {
			continue
		}

		if char == rune(SYMBOL_DOT) {
			if seen_dot {
				lexer.throw_error("Invalid numeric literal")
			}

			seen_dot = true
			continue
		}

		lexer.index -= 1
		break
	}

	if seen_dot && lexer.index == lexer.start {
		lexer.throw_error("Invalid numeric literal")
	}

	return Token{TOKEN_NUMBER_LITERAL, lexer.source[lexer.start : lexer.index+1], lexer.start}
}

func (lexer *Lexer) lex_text() Token {
	lexer.start = lexer.index

	lexer.next_rune()

	for lexer.index < lexer.source_length {
		char := lexer.next_rune()

		if char == rune(SYMBOL_SINGLE_QUOTE) {
			return Token{TOKEN_TEXT_LITERAL, lexer.source[lexer.start+1 : lexer.index], lexer.start}
		}
	}

	lexer.index -= 1
	lexer.throw_error("Non-terminated text literal")
	panic("")
}

func (lexer *Lexer) lex_keyword_or_identifier() Token {
	lexer.start = lexer.index

	for {
		char := lexer.next_rune()

		if !is_alphanumeric(char) {
			lexer.index -= 1
			break
		}
	}

	char := rune(lexer.source[lexer.start])
	switch to_upper_rune(char) {
	case 'B':
		return lexer.check_keyword(1, 6, "OOLEAN")
	case 'C':
		return lexer.check_keyword(1, 5, "REATE")
	case 'F':
		if lexer.index-lexer.start > 1 {
			char = rune(lexer.source[lexer.start+1])
			switch to_upper_rune(char) {
			case 'A':
				return lexer.check_keyword(2, 3, "LSE")
			case 'R':
				return lexer.check_keyword(2, 2, "OM")
			}
		}
	case 'I':
		if lexer.index-lexer.start > 1 {
			char = rune(lexer.source[lexer.start+1])
			switch to_upper_rune(char) {
			case 'N':
				if lexer.index-lexer.start > 2 {
					char = rune(lexer.source[lexer.start+2])
					switch to_upper_rune(char) {
					case 'S':
						return lexer.check_keyword(3, 3, "ERT")
					case 'T':
						return lexer.check_keyword(3, 1, "O")
					}
				}
			}
		}
	case 'N':
		return lexer.check_keyword(1, 5, "UMBER")
	case 'S':
		return lexer.check_keyword(1, 5, "ELECT")
	case 'T':
		if lexer.index-lexer.start > 1 {
			char = rune(lexer.source[lexer.start+1])
			switch to_upper_rune(char) {
			case 'A':
				return lexer.check_keyword(2, 3, "BLE")
			case 'E':
				return lexer.check_keyword(2, 2, "XT")
			case 'R':
				return lexer.check_keyword(2, 2, "UE")
			}
		}
	case 'V':
		return lexer.check_keyword(1, 5, "ALUES")
	}

	token := Token{TOKEN_IDENTIFIER, lexer.source[lexer.start : lexer.index+1], lexer.start}
	return token
}

func (lexer *Lexer) check_keyword(start int, length int, suffix string) Token {
	if lexer.index+1-lexer.start == start+length {
		for i, keyword_char := range suffix {
			char := to_upper_rune(rune(lexer.source[lexer.start+start+i]))
			if char != keyword_char {
				break
			}
		}
		return Token{TOKEN_KEYWORD, strings.ToUpper(lexer.source[lexer.start : lexer.index+1]), lexer.start}
	}
	return Token{TOKEN_IDENTIFIER, lexer.source[lexer.start : lexer.index+1], lexer.start}
}
