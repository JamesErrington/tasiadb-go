package lexer

import (
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
	tokens := GenerateTokenSlice(";,( )*")
	expected := []Token{
		{TOKEN_SEMI_COLON, "", 0}, {TOKEN_COMMA, "", 1}, {TOKEN_LEFT_PAREN, "", 2},
		{TOKEN_RIGHT_PAREN, "", 4}, {TOKEN_ASTERISK, "", 5}, {TOKEN_EOF, "", 6},
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
		{TOKEN_LITERAL_NUMBER, "1", 0}, {TOKEN_LITERAL_NUMBER, "2.34", 2}, {TOKEN_LITERAL_NUMBER, "500", 7},
		{TOKEN_LITERAL_NUMBER, "06", 11}, {TOKEN_LITERAL_NUMBER, "07.80", 14}, {TOKEN_LITERAL_NUMBER, ".9", 20},
		{TOKEN_LITERAL_NUMBER, "1.", 23}, {TOKEN_EOF, "", 25},
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
		{TOKEN_LITERAL_TEXT, "a", 0}, {TOKEN_LITERAL_TEXT, "b12", 4},
		{TOKEN_LITERAL_TEXT, "cd3_4ef", 10}, {TOKEN_LITERAL_TEXT, ";,()*.", 20},
		{TOKEN_EOF, "", 28},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexTextErrors(t *testing.T) {
	tokens := GenerateTokenSlice("'abcd")
	expected := []Token{{TOKEN_ERROR, "Non-terminated text literal", 4}, {TOKEN_EOF, "", 5}}

	assert.Equal(t, expected, tokens)
}

func TestLexKeywordUpper(t *testing.T) {
	tokens := GenerateTokenSlice("BOOLEAN CREATE FALSE FROM INSERT INTO NUMBER SELECT TABLE TEXT TRUE VALUES")
	expected := []Token{
		{TOKEN_KEYWORD_BOOLEAN, "", 0}, {TOKEN_KEYWORD_CREATE, "", 8}, {TOKEN_KEYWORD_FALSE, "", 15},
		{TOKEN_KEYWORD_FROM, "", 21}, {TOKEN_KEYWORD_INSERT, "", 26}, {TOKEN_KEYWORD_INTO, "", 33},
		{TOKEN_KEYWORD_NUMBER, "", 38}, {TOKEN_KEYWORD_SELECT, "", 45}, {TOKEN_KEYWORD_TABLE, "", 52},
		{TOKEN_KEYWORD_TEXT, "", 58}, {TOKEN_KEYWORD_TRUE, "", 63}, {TOKEN_KEYWORD_VALUES, "", 68},
		{TOKEN_EOF, "", 74},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexKeywordLower(t *testing.T) {
	tokens := GenerateTokenSlice("boolean create false from insert into number select table text true values")
	expected := []Token{
		{TOKEN_KEYWORD_BOOLEAN, "", 0}, {TOKEN_KEYWORD_CREATE, "", 8}, {TOKEN_KEYWORD_FALSE, "", 15},
		{TOKEN_KEYWORD_FROM, "", 21}, {TOKEN_KEYWORD_INSERT, "", 26}, {TOKEN_KEYWORD_INTO, "", 33},
		{TOKEN_KEYWORD_NUMBER, "", 38}, {TOKEN_KEYWORD_SELECT, "", 45}, {TOKEN_KEYWORD_TABLE, "", 52},
		{TOKEN_KEYWORD_TEXT, "", 58}, {TOKEN_KEYWORD_TRUE, "", 63}, {TOKEN_KEYWORD_VALUES, "", 68},
		{TOKEN_EOF, "", 74},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexKeywordMixed(t *testing.T) {
	tokens := GenerateTokenSlice("boOLEan Create falsE fRom INSERt InTo nUmbeR SEleCt taBle TExT trUE vAlUeS")
	expected := []Token{
		{TOKEN_KEYWORD_BOOLEAN, "", 0}, {TOKEN_KEYWORD_CREATE, "", 8}, {TOKEN_KEYWORD_FALSE, "", 15},
		{TOKEN_KEYWORD_FROM, "", 21}, {TOKEN_KEYWORD_INSERT, "", 26}, {TOKEN_KEYWORD_INTO, "", 33},
		{TOKEN_KEYWORD_NUMBER, "", 38}, {TOKEN_KEYWORD_SELECT, "", 45}, {TOKEN_KEYWORD_TABLE, "", 52},
		{TOKEN_KEYWORD_TEXT, "", 58}, {TOKEN_KEYWORD_TRUE, "", 63}, {TOKEN_KEYWORD_VALUES, "", 68},
		{TOKEN_EOF, "", 74},
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
		{TOKEN_KEYWORD_CREATE, "", 0}, {TOKEN_KEYWORD_TABLE, "", 7}, {TOKEN_IDENTIFIER, "t", 13},
		{TOKEN_LEFT_PAREN, "", 15}, {TOKEN_IDENTIFIER, "c_1", 16}, {TOKEN_KEYWORD_NUMBER, "", 20},
		{TOKEN_COMMA, "", 26}, {TOKEN_IDENTIFIER, "c_2", 28}, {TOKEN_KEYWORD_TEXT, "", 32},
		{TOKEN_RIGHT_PAREN, "", 36}, {TOKEN_SEMI_COLON, "", 37}, {TOKEN_EOF, "", 38},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexInsertInto(t *testing.T) {
	tokens := GenerateTokenSlice("insert into t values (c_1 10.5, c_2 'Hello $ % !');")
	expected := []Token{
		{TOKEN_KEYWORD_INSERT, "", 0}, {TOKEN_KEYWORD_INTO, "", 7}, {TOKEN_IDENTIFIER, "t", 12},
		{TOKEN_KEYWORD_VALUES, "", 14}, {TOKEN_LEFT_PAREN, "", 21}, {TOKEN_IDENTIFIER, "c_1", 22},
		{TOKEN_LITERAL_NUMBER, "10.5", 26}, {TOKEN_COMMA, "", 30}, {TOKEN_IDENTIFIER, "c_2", 32},
		{TOKEN_LITERAL_TEXT, "Hello $ % !", 36}, {TOKEN_RIGHT_PAREN, "", 49}, {TOKEN_SEMI_COLON, "", 50},
		{TOKEN_EOF, "", 51},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexSelect(t *testing.T) {
	tokens := GenerateTokenSlice("SELECT c_1, c_2 FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD_SELECT, "", 0}, {TOKEN_IDENTIFIER, "c_1", 7}, {TOKEN_COMMA, "", 10},
		{TOKEN_IDENTIFIER, "c_2", 12}, {TOKEN_KEYWORD_FROM, "", 16}, {TOKEN_IDENTIFIER, "t", 21},
		{TOKEN_SEMI_COLON, "", 22}, {TOKEN_EOF, "", 23},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexSelectStar(t *testing.T) {
	tokens := GenerateTokenSlice("SELECT * FROM t;")
	expected := []Token{
		{TOKEN_KEYWORD_SELECT, "", 0}, {TOKEN_ASTERISK, "", 7}, {TOKEN_KEYWORD_FROM, "", 9},
		{TOKEN_IDENTIFIER, "t", 14}, {TOKEN_SEMI_COLON, "", 15}, {TOKEN_EOF, "", 16},
	}

	assert.Equal(t, expected, tokens)
}

func TestLexNonUtf8(t *testing.T) {
	tokens := GenerateTokenSlice("\xc5")
	expected := []Token{{TOKEN_ERROR, "Unidentified token", 0}, {TOKEN_EOF, "", 1}}

	assert.Equal(t, expected, tokens)
}
