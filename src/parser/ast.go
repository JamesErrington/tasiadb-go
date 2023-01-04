package parser

import "github.com/JamesErrington/tasiadb/src/lexer"

type NodeType int8

const (
	NODE_NUMBER_VALUE NodeType = iota
	NODE_TEXT_VALUE
	NODE_BOOLEAN_VALUE
	NODE_CREATE_TABLE_STATEMENT
)

type Node interface {
	Pos() int
}

type Statement struct {
	Content Node
}

func (s *Statement) Pos() int {
	return s.Content.Pos()
}

type CreateTableStatement struct {
	_type       NodeType
	create      int
	table_name  lexer.Token
	column_defs []ColumnDefinition
}

func (s *CreateTableStatement) Pos() int {
	return s.create
}

type InsertStatement struct {
	_type      NodeType
	table_name lexer.Token
	values     []ColumnValue
}

type ColumnDefinition struct {
	colum_name  lexer.Token
	column_type lexer.Token
}

type ColumnValue struct {
	column_name  lexer.Token
	column_value lexer.Token
}
