package parser

import lex "github.com/JamesErrington/tasiadb/src/lexer"

type NodeType int8

const (
	NODE_NUMBER_VALUE NodeType = iota
	NODE_TEXT_VALUE
	NODE_BOOLEAN_VALUE
	NODE_CREATE_TABLE_STATEMENT
	NODE_INSERT_STATEMENT
	NODE_SELECT_STATEMENT
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
	start       int
	table_name  lex.Token
	column_defs []ColumnDefinition
}

func (s *CreateTableStatement) Pos() int {
	return s.start
}

type InsertStatement struct {
	_type      NodeType
	start      int
	table_name lex.Token
	values     []ColumnValue
}

func (s *InsertStatement) Pos() int {
	return s.start
}

type SelectStatement struct {
	_type      NodeType
	start      int
	columns    []lex.Token
	table_name lex.Token
}

func (s *SelectStatement) Pos() int {
	return s.start
}

type ColumnDefinition struct {
	colum_name  lex.Token
	column_type lex.Token
}

type ColumnValue struct {
	column_name  lex.Token
	column_value lex.Token
}
