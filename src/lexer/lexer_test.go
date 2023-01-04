package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GenerateTokenSlice(source string) []Token {
	lexer := NewLexer(source)
	var tokens []Token

	for lexer.index < lexer.source_length {
		token, finished := lexer.NextToken()

		if finished {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func TestLexSymbol(t *testing.T) {
	tokens := GenerateTokenSlice(",()*")
	expected := []Token{
		{TOKEN_SYMBOL, ",", 0}, {TOKEN_SYMBOL, "(", 1},
		{TOKEN_SYMBOL, ")", 2}, {TOKEN_SYMBOL, "*", 3},
		{TOKEN_EOF, "", 4},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexSymbolErrors(t *testing.T) {
	tokens := GenerateTokenSlice("! / $")
	expected := []Token{
		{TOKEN_ERROR, "Unidentified token", 0}, {TOKEN_ERROR, "Unidentified token", 2}, {TOKEN_ERROR, "Unidentified token", 4},
		{TOKEN_EOF, "", 5},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexNumber(t *testing.T) {
	tokens := GenerateTokenSlice("1 2.34 500 06 07.80 .9 1.")
	expected := []Token{
		{TOKEN_NUMBER_LITERAL, "1", 0}, {TOKEN_NUMBER_LITERAL, "2.34", 2}, {TOKEN_NUMBER_LITERAL, "500", 7},
		{TOKEN_NUMBER_LITERAL, "06", 11}, {TOKEN_NUMBER_LITERAL, "07.80", 14}, {TOKEN_NUMBER_LITERAL, ".9", 20},
		{TOKEN_NUMBER_LITERAL, "1.", 23}, {TOKEN_EOF, "", 25},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexNumberErrors(t *testing.T) {
	tokens := GenerateTokenSlice(". 3..4")
	expected := []Token{
		{TOKEN_ERROR, "Invalid number literal", 0}, {TOKEN_ERROR, "Invalid number literal", 2}, {TOKEN_EOF, "", 6},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexText(t *testing.T) {
	tokens := GenerateTokenSlice("'a' 'b12' 'cd3_4ef' ';,()*.'")
	expected := []Token{
		{TOKEN_TEXT_LITERAL, "a", 0}, {TOKEN_TEXT_LITERAL, "b12", 4},
		{TOKEN_TEXT_LITERAL, "cd3_4ef", 10}, {TOKEN_TEXT_LITERAL, ";,()*.", 20},
		{TOKEN_EOF, "", 28},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexTextErrors(t *testing.T) {
	tokens := GenerateTokenSlice("'abcd")
	expected := []Token{{TOKEN_ERROR, "Non-terminated text literal", 4}, {TOKEN_EOF, "", 5}}

	assert.Equal(t, expected, tokens)
}

func TestLexKeyword(t *testing.T) {
	for _, keyword := range []string{"CREATE", "create", "TABLE", "taBlE", "NUMBER", "TEXT", "BOOLEAN", "INSERT", "insert", "INTO",
		"VALUES", "ValueS", "TRUE", "FALSE", "faLse", "SELECT", "select", "FROM"} {
		tokens := GenerateTokenSlice(keyword)
		expected := []Token{{TOKEN_KEYWORD, strings.ToUpper(keyword), 0}, {TOKEN_EOF, "", len(keyword)}}

		assert.Equal(t, expected, tokens)
	}
}

func TestLexKeywords(t *testing.T) {
	tokens := GenerateTokenSlice("create CREATE inSeRt InSERT")
	expected := []Token{
		{TOKEN_KEYWORD, "CREATE", 0}, {TOKEN_KEYWORD, "CREATE", 7},
		{TOKEN_KEYWORD, "INSERT", 14}, {TOKEN_KEYWORD, "INSERT", 21},
		{TOKEN_EOF, "", 27},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexIdentifiers(t *testing.T) {
	tokens := GenerateTokenSlice("table_1 column_2_b TABLE_3 false4 tabl tabler")
	expected := []Token{
		{TOKEN_IDENTIFIER, "table_1", 0}, {TOKEN_IDENTIFIER, "column_2_b", 8},
		{TOKEN_IDENTIFIER, "TABLE_3", 19}, {TOKEN_IDENTIFIER, "false4", 27},
		{TOKEN_IDENTIFIER, "tabl", 34}, {TOKEN_IDENTIFIER, "tabler", 39},
		{TOKEN_EOF, "", 45},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexCreateTable(t *testing.T) {
	tokens := GenerateTokenSlice("CREATE TABLE t (c_1 NUMBER, c_2 TEXT);")
	expected := []Token{
		{TOKEN_KEYWORD, "CREATE", 0}, {TOKEN_KEYWORD, "TABLE", 7}, {TOKEN_IDENTIFIER, "t", 13},
		{TOKEN_SYMBOL, "(", 15}, {TOKEN_IDENTIFIER, "c_1", 16}, {TOKEN_KEYWORD, "NUMBER", 20},
		{TOKEN_SYMBOL, ",", 26}, {TOKEN_IDENTIFIER, "c_2", 28}, {TOKEN_KEYWORD, "TEXT", 32},
		{TOKEN_SYMBOL, ")", 36}, {TOKEN_SYMBOL, ";", 37}, {TOKEN_EOF, "", 38},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexInsertInto(t *testing.T) {
	tokens := GenerateTokenSlice("insert into t values (c_1 10.5, c_2 'Hello $ % !');")
	expected := []Token{
		{TOKEN_KEYWORD, "INSERT", 0}, {TOKEN_KEYWORD, "INTO", 7}, {TOKEN_IDENTIFIER, "t", 12},
		{TOKEN_KEYWORD, "VALUES", 14}, {TOKEN_SYMBOL, "(", 21}, {TOKEN_IDENTIFIER, "c_1", 22},
		{TOKEN_NUMBER_LITERAL, "10.5", 26}, {TOKEN_SYMBOL, ",", 30}, {TOKEN_IDENTIFIER, "c_2", 32},
		{TOKEN_TEXT_LITERAL, "Hello $ % !", 36}, {TOKEN_SYMBOL, ")", 49}, {TOKEN_SYMBOL, ";", 50},
		{TOKEN_EOF, "", 51},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexSelect(t *testing.T) {
	tokens := GenerateTokenSlice("SELECT c_1, c_2 FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD, "SELECT", 0}, {TOKEN_IDENTIFIER, "c_1", 7}, {TOKEN_SYMBOL, ",", 10},
		{TOKEN_IDENTIFIER, "c_2", 12}, {TOKEN_KEYWORD, "FROM", 16}, {TOKEN_IDENTIFIER, "t", 21},
		{TOKEN_SYMBOL, ";", 22}, {TOKEN_EOF, "", 23},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexSelectStar(t *testing.T) {
	tokens := GenerateTokenSlice("SELECT * FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD, "SELECT", 0}, {TOKEN_SYMBOL, "*", 7}, {TOKEN_KEYWORD, "FROM", 9},
		{TOKEN_IDENTIFIER, "t", 14}, {TOKEN_SYMBOL, ";", 15}, {TOKEN_EOF, "", 16},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexNonUtf8(t *testing.T) {
	tokens := GenerateTokenSlice("\xc5")
	expected := []Token{{TOKEN_ERROR, "Unidentified token", 0}, {TOKEN_EOF, "", 1}}

	assert.Equal(t, expected, tokens)
}
