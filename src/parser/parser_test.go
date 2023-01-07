package parser

import (
	"testing"

	lex "github.com/JamesErrington/tasiadb/src/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseCreateTableSingleColumn(t *testing.T) {
	parser := NewParser("CREATE TABLE t (c_1 NUMBER);")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*CreateTableStatement)
	assert.Equal(t, &CreateTableStatement{
		NODE_CREATE_TABLE_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 13),
		[]lex.Token{lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 16)},
		[]lex.Token{lex.MakeToken(lex.TOKEN_KEYWORD_NUMBER, "", 20)},
	}, content)
}

func TestParseCreateTableMultiColumn(t *testing.T) {
	parser := NewParser("CREATE TABLE t (c_1 NUMBER, c_2 TEXT,c_3 BOOLEAN);")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*CreateTableStatement)
	assert.Equal(t, &CreateTableStatement{
		NODE_CREATE_TABLE_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 13),
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 16),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_2", 28),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_3", 37),
		},
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_KEYWORD_NUMBER, "", 20),
			lex.MakeToken(lex.TOKEN_KEYWORD_TEXT, "", 32),
			lex.MakeToken(lex.TOKEN_KEYWORD_BOOLEAN, "", 41),
		},
	}, content)
}

func TestParseInsertSingleColumn(t *testing.T) {
	parser := NewParser("INSERT INTO t (c_1) VALUES (10.5);")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*InsertStatement)
	assert.Equal(t, &InsertStatement{
		NODE_INSERT_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 12),
		[]lex.Token{lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 15)},
		[]lex.Token{lex.MakeToken(lex.TOKEN_LITERAL_NUMBER, "10.5", 28)},
	}, content)
}

func TestParseInsertMultiColumn(t *testing.T) {
	parser := NewParser("INSERT INTO t (c_1, c_2,c_3) VALUES (10.5, 'Hello',FALSE);")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*InsertStatement)
	assert.Equal(t, &InsertStatement{
		NODE_INSERT_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 12),
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 15),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_2", 20),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_3", 24),
		},
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_LITERAL_NUMBER, "10.5", 37),
			lex.MakeToken(lex.TOKEN_LITERAL_TEXT, "Hello", 43),
			lex.MakeToken(lex.TOKEN_KEYWORD_FALSE, "", 51),
		},
	}, content)
}

func TestParseInsertImplicitColumn(t *testing.T) {
	parser := NewParser("INSERT INTO t VALUES (10.5);")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*InsertStatement)
	assert.Equal(t, &InsertStatement{
		NODE_INSERT_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 12),
		nil,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_LITERAL_NUMBER, "10.5", 22),
		},
	}, content)
}

func TestParseSelectSingleColumn(t *testing.T) {
	parser := NewParser("SELECT c_1 FROM t;")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*SelectStatement)
	assert.Equal(t, &SelectStatement{
		NODE_SELECT_STATEMENT,
		0,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 7),
		},
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 16),
	}, content)
}

func TestParseSelectMultiColumn(t *testing.T) {
	parser := NewParser("SELECT c_1, c_2,c_3 FROM t;")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*SelectStatement)
	assert.Equal(t, &SelectStatement{
		NODE_SELECT_STATEMENT,
		0,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 7),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_2", 12),
			lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_3", 16),
		},
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 25),
	}, content)
}

func TestParseSelectStar(t *testing.T) {
	parser := NewParser("SELECT * FROM t;")
	result := parser.Parse()

	assert.Len(t, result, 1)
	content := result[0].Content.(*SelectStatement)
	assert.Equal(t, &SelectStatement{
		NODE_SELECT_STATEMENT,
		0,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_ASTERISK, "", 7),
		},
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 14),
	}, content)
}

func TestMultipleStatements(t *testing.T) {
	parser := NewParser("CREATE TABLE t (c_1 TEXT); INSERT INTO t VALUES ('James'); SELECT * FROM t;")
	result := parser.Parse()

	assert.Len(t, result, 3)
	create_stmt := result[0].Content.(*CreateTableStatement)
	assert.Equal(t, &CreateTableStatement{
		NODE_CREATE_TABLE_STATEMENT,
		0,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 13),
		[]lex.Token{lex.MakeToken(lex.TOKEN_IDENTIFIER, "c_1", 16)},
		[]lex.Token{lex.MakeToken(lex.TOKEN_KEYWORD_TEXT, "", 20)},
	}, create_stmt)

	insert_stmt := result[1].Content.(*InsertStatement)
	assert.Equal(t, &InsertStatement{
		NODE_INSERT_STATEMENT,
		27,
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 39),
		nil,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_LITERAL_TEXT, "James", 49),
		},
	}, insert_stmt)

	select_stmt := result[2].Content.(*SelectStatement)
	assert.Equal(t, &SelectStatement{
		NODE_SELECT_STATEMENT,
		59,
		[]lex.Token{
			lex.MakeToken(lex.TOKEN_ASTERISK, "", 66),
		},
		lex.MakeToken(lex.TOKEN_IDENTIFIER, "t", 73),
	}, select_stmt)
}
