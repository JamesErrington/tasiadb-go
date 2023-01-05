package lexer

import (
	"unicode/utf8"
)

const (
	SYMBOL_EOF          = -1
	SYMBOL_SEMI_COLON   = ';'
	SYMBOL_COMMA        = ','
	SYMBOL_LEFT_PAREN   = '('
	SYMBOL_RIGHT_PAREN  = ')'
	SYMBOL_SINGLE_QUOTE = '\''
	SYMBOL_UNDERSCORE   = '_'
	SYMBOL_ASTERISK     = '*'
	SYMBOL_DOT          = '.'
	SYMBOL_SPACE        = ' '
	SYMBOL_TAB          = '\t'
	SYMBOL_NEWLINE      = '\n'
)

type TokenType uint8

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_SEMI_COLON
	TOKEN_COMMA
	TOKEN_LEFT_PAREN
	TOKEN_RIGHT_PAREN
	TOKEN_ASTERISK

	TOKEN_KEYWORD_CREATE
	TOKEN_KEYWORD_TABLE
	TOKEN_KEYWORD_NUMBER
	TOKEN_KEYWORD_TEXT
	TOKEN_KEYWORD_BOOLEAN
	TOKEN_KEYWORD_INSERT
	TOKEN_KEYWORD_INTO
	TOKEN_KEYWORD_VALUES
	TOKEN_KEYWORD_TRUE
	TOKEN_KEYWORD_FALSE
	TOKEN_KEYWORD_SELECT
	TOKEN_KEYWORD_FROM

	TOKEN_IDENTIFIER
	TOKEN_LITERAL_NUMBER
	TOKEN_LITERAL_TEXT
)

func is_whitespace(char rune) bool {
	switch char {
	case SYMBOL_SPACE, SYMBOL_NEWLINE, SYMBOL_TAB:
		return true
	default:
		return false
	}
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

func MakeToken(_type TokenType, value string, offset int) Token {
	return Token{_type, value, offset}
}

func (token Token) IsTokenType(token_type TokenType) bool {
	return token._type == token_type
}

func (token Token) IsDataType() bool {
	return token._type == TOKEN_KEYWORD_BOOLEAN || token._type == TOKEN_KEYWORD_NUMBER || token._type == TOKEN_KEYWORD_TEXT
}

func (token Token) IsValueType() bool {
	return token._type == TOKEN_LITERAL_NUMBER || token._type == TOKEN_LITERAL_TEXT || token._type == TOKEN_KEYWORD_FALSE || token._type == TOKEN_KEYWORD_TRUE
}

func (token Token) Value() string {
	return token.value
}

func (token Token) Offset() int {
	return token.offset
}

type Lexer struct {
	source        string
	source_length int
	index         int
	start         int
}

func (lexer *Lexer) CurrentIndex() int {
	return lexer.index
}

func NewLexer(source string) *Lexer {
	return &Lexer{source, len(source), -1, 0}
}

func (lexer *Lexer) NextToken() (Token, bool) {
	for lexer.index < lexer.source_length {
		char := lexer.next_rune()

		if char < rune(SYMBOL_EOF) || char > utf8.MaxRune {
			return Token{TOKEN_ERROR, string(char), lexer.index}, false
		}

		if is_whitespace(char) {
			continue
		}

		switch char {
		case SYMBOL_EOF:
			return Token{TOKEN_EOF, "", lexer.index}, false
		case SYMBOL_SEMI_COLON:
			return Token{TOKEN_SEMI_COLON, "", lexer.index}, false
		case SYMBOL_COMMA:
			return Token{TOKEN_COMMA, "", lexer.index}, false
		case SYMBOL_LEFT_PAREN:
			return Token{TOKEN_LEFT_PAREN, "", lexer.index}, false
		case SYMBOL_RIGHT_PAREN:
			return Token{TOKEN_RIGHT_PAREN, "", lexer.index}, false
		case SYMBOL_ASTERISK:
			return Token{TOKEN_ASTERISK, "", lexer.index}, false
		case SYMBOL_SINGLE_QUOTE:
			token := lexer.lex_text()
			return token, false
		default:
			switch {
			case is_digit(char) || char == rune(SYMBOL_DOT):
				token := lexer.lex_number()
				return token, false
			case is_alphabetical(char):
				token := lexer.lex_keyword_or_identifier()
				return token, false
			default:
				return Token{TOKEN_ERROR, "Unidentified token", lexer.index}, false
			}
		}
	}

	return Token{}, true
}

func (lexer *Lexer) next_rune() rune {
	lexer.index += 1

	if lexer.index >= lexer.source_length {
		return rune(SYMBOL_EOF)
	}

	char := rune(lexer.source[lexer.index])

	return char
}

func (lexer *Lexer) lex_number() Token {
	lexer.start = lexer.index
	lexer.index -= 1
	seen_dot := false
	has_error := false

	for lexer.index < lexer.source_length {
		char := lexer.next_rune()

		if char == rune(SYMBOL_DOT) {
			if seen_dot {
				has_error = true
			}

			seen_dot = true
			continue
		}

		if !is_digit(char) {
			lexer.index -= 1
			break
		}
	}

	if seen_dot && lexer.index == lexer.start {
		has_error = true
	}

	if has_error {
		return Token{TOKEN_ERROR, "Invalid number literal", lexer.start}
	}

	return Token{TOKEN_LITERAL_NUMBER, lexer.source[lexer.start : lexer.index+1], lexer.start}
}

func (lexer *Lexer) lex_text() Token {
	lexer.start = lexer.index

	lexer.next_rune()

	for lexer.index < lexer.source_length {
		char := lexer.next_rune()

		if char == rune(SYMBOL_SINGLE_QUOTE) {
			return Token{TOKEN_LITERAL_TEXT, lexer.source[lexer.start+1 : lexer.index], lexer.start}
		}
	}

	lexer.index -= 1
	return Token{TOKEN_ERROR, "Non-terminated text literal", lexer.index}
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
		return lexer.check_keyword(1, 6, "OOLEAN", TOKEN_KEYWORD_BOOLEAN)
	case 'C':
		return lexer.check_keyword(1, 5, "REATE", TOKEN_KEYWORD_CREATE)
	case 'F':
		if lexer.index-lexer.start > 1 {
			char = rune(lexer.source[lexer.start+1])
			switch to_upper_rune(char) {
			case 'A':
				return lexer.check_keyword(2, 3, "LSE", TOKEN_KEYWORD_FALSE)
			case 'R':
				return lexer.check_keyword(2, 2, "OM", TOKEN_KEYWORD_FROM)
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
						return lexer.check_keyword(3, 3, "ERT", TOKEN_KEYWORD_INSERT)
					case 'T':
						return lexer.check_keyword(3, 1, "O", TOKEN_KEYWORD_INTO)
					}
				}
			}
		}
	case 'N':
		return lexer.check_keyword(1, 5, "UMBER", TOKEN_KEYWORD_NUMBER)
	case 'S':
		return lexer.check_keyword(1, 5, "ELECT", TOKEN_KEYWORD_SELECT)
	case 'T':
		if lexer.index-lexer.start > 1 {
			char = rune(lexer.source[lexer.start+1])
			switch to_upper_rune(char) {
			case 'A':
				return lexer.check_keyword(2, 3, "BLE", TOKEN_KEYWORD_TABLE)
			case 'E':
				return lexer.check_keyword(2, 2, "XT", TOKEN_KEYWORD_TEXT)
			case 'R':
				return lexer.check_keyword(2, 2, "UE", TOKEN_KEYWORD_TRUE)
			}
		}
	case 'V':
		return lexer.check_keyword(1, 5, "ALUES", TOKEN_KEYWORD_VALUES)
	}

	token := Token{TOKEN_IDENTIFIER, lexer.source[lexer.start : lexer.index+1], lexer.start}
	return token
}

func (lexer *Lexer) check_keyword(start int, length int, suffix string, token_type TokenType) Token {
	if lexer.index+1-lexer.start == start+length {
		for i, keyword_char := range suffix {
			char := to_upper_rune(rune(lexer.source[lexer.start+start+i]))
			if char != keyword_char {
				break
			}
		}
		return Token{token_type, "", lexer.start}
	}
	return Token{TOKEN_IDENTIFIER, lexer.source[lexer.start : lexer.index+1], lexer.start}
}
