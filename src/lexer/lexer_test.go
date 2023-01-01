package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexSymbol(t *testing.T) {
	lexer := NewLexer(",()*")
	expected := []Token{
		{TOKEN_SYMBOL, ",", 0}, {TOKEN_SYMBOL, "(", 1},
		{TOKEN_SYMBOL, ")", 2}, {TOKEN_SYMBOL, "*", 3},
		{TOKEN_EOF, "", 4},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexSymbolErrors(t *testing.T) {
	for _, test := range []string{"!", "/", "$"} {
		lexer := NewLexer(test)
		expected := []Token(nil)

		tokens := lexer.Lex()
		assert.Equal(t, expected, tokens)
	}
}

func TestLexNumber(t *testing.T) {
	lexer := NewLexer("1 2.34 500 06 07.80 .9 1.")

	expected := []Token{
		{TOKEN_NUMBER_LITERAL, "1", 0}, {TOKEN_NUMBER_LITERAL, "2.34", 2}, {TOKEN_NUMBER_LITERAL, "500", 7},
		{TOKEN_NUMBER_LITERAL, "06", 11}, {TOKEN_NUMBER_LITERAL, "07.80", 14}, {TOKEN_NUMBER_LITERAL, ".9", 20},
		{TOKEN_NUMBER_LITERAL, "1.", 23}, {TOKEN_EOF, "", 25},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexNumberErrors(t *testing.T) {
	for _, test := range []string{".", "1.2.3"} {
		lexer := NewLexer(test)
		expected := []Token(nil)

		tokens := lexer.Lex()
		assert.Equal(t, expected, tokens)
	}
}

func TestLexText(t *testing.T) {
	lexer := NewLexer("'a' 'b12' 'cd3_4ef' ';,()*.'")
	expected := []Token{
		{TOKEN_TEXT_LITERAL, "a", 0}, {TOKEN_TEXT_LITERAL, "b12", 4},
		{TOKEN_TEXT_LITERAL, "cd3_4ef", 10}, {TOKEN_TEXT_LITERAL, ";,()*.", 20},
		{TOKEN_EOF, "", 28},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexTextErrors(t *testing.T) {
	lexer := NewLexer("'abcd")
	expected := []Token(nil)

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexKeyword(t *testing.T) {
	for _, keyword := range []string{"CREATE", "create", "TABLE", "taBlE", "NUMBER", "TEXT", "BOOLEAN", "INSERT", "insert", "INTO",
		"VALUES", "ValueS", "TRUE", "FALSE", "faLse", "SELECT", "select", "FROM"} {
		lexer := NewLexer(keyword)
		expected := []Token{{TOKEN_KEYWORD, strings.ToUpper(keyword), 0}, {TOKEN_EOF, "", len(keyword)}}

		tokens := lexer.Lex()
		assert.Equal(t, expected, tokens)
	}
}

func TestLexKeywords(t *testing.T) {
	lexer := NewLexer("create CREATE inSeRt InSERT")
	expected := []Token{
		{TOKEN_KEYWORD, "CREATE", 0}, {TOKEN_KEYWORD, "CREATE", 7},
		{TOKEN_KEYWORD, "INSERT", 14}, {TOKEN_KEYWORD, "INSERT", 21},
		{TOKEN_EOF, "", 27},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexIdentifiers(t *testing.T) {
	lexer := NewLexer("table_1 column_2_b TABLE_3 false4 tabl tabler")
	expected := []Token{
		{TOKEN_IDENTIFIER, "table_1", 0}, {TOKEN_IDENTIFIER, "column_2_b", 8},
		{TOKEN_IDENTIFIER, "TABLE_3", 19}, {TOKEN_IDENTIFIER, "false4", 27},
		{TOKEN_IDENTIFIER, "tabl", 34}, {TOKEN_IDENTIFIER, "tabler", 39},
		{TOKEN_EOF, "", 45},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexCreateTable(t *testing.T) {
	lexer := NewLexer("CREATE TABLE t (c_1 NUMBER, c_2 TEXT);")
	expected := []Token{
		{TOKEN_KEYWORD, "CREATE", 0}, {TOKEN_KEYWORD, "TABLE", 7}, {TOKEN_IDENTIFIER, "t", 13},
		{TOKEN_SYMBOL, "(", 15}, {TOKEN_IDENTIFIER, "c_1", 16}, {TOKEN_KEYWORD, "NUMBER", 20},
		{TOKEN_SYMBOL, ",", 26}, {TOKEN_IDENTIFIER, "c_2", 28}, {TOKEN_KEYWORD, "TEXT", 32},
		{TOKEN_SYMBOL, ")", 36}, {TOKEN_SYMBOL, ";", 37}, {TOKEN_EOF, "", 38},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexInsertInto(t *testing.T) {
	lexer := NewLexer("insert into t values (c_1 10.5, c_2 'Hello $ % !');")
	expected := []Token{
		{TOKEN_KEYWORD, "INSERT", 0}, {TOKEN_KEYWORD, "INTO", 7}, {TOKEN_IDENTIFIER, "t", 12},
		{TOKEN_KEYWORD, "VALUES", 14}, {TOKEN_SYMBOL, "(", 21}, {TOKEN_IDENTIFIER, "c_1", 22},
		{TOKEN_NUMBER_LITERAL, "10.5", 26}, {TOKEN_SYMBOL, ",", 30}, {TOKEN_IDENTIFIER, "c_2", 32},
		{TOKEN_TEXT_LITERAL, "Hello $ % !", 36}, {TOKEN_SYMBOL, ")", 49}, {TOKEN_SYMBOL, ";", 50},
		{TOKEN_EOF, "", 51},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexSelect(t *testing.T) {
	lexer := NewLexer("SELECT c_1, c_2 FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD, "SELECT", 0}, {TOKEN_IDENTIFIER, "c_1", 7}, {TOKEN_SYMBOL, ",", 10},
		{TOKEN_IDENTIFIER, "c_2", 12}, {TOKEN_KEYWORD, "FROM", 16}, {TOKEN_IDENTIFIER, "t", 21},
		{TOKEN_SYMBOL, ";", 22}, {TOKEN_EOF, "", 23},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexSelectStar(t *testing.T) {
	lexer := NewLexer("SELECT * FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD, "SELECT", 0}, {TOKEN_SYMBOL, "*", 7}, {TOKEN_KEYWORD, "FROM", 9},
		{TOKEN_IDENTIFIER, "t", 14}, {TOKEN_SYMBOL, ";", 15}, {TOKEN_EOF, "", 16},
	}

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}

func TestLexNonUtf8(t *testing.T) {
	lexer := NewLexer("\xc5")
	expected := []Token(nil)

	tokens := lexer.Lex()
	assert.Equal(t, expected, tokens)
}
