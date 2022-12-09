package main

import (
	"errors"
)

type Ast struct {
	statements []*Statement
}

type StatementType uint

const (
	CREATE_TABLE_STATEMENT StatementType = iota
)

type Statement struct {
	_type                  StatementType
	create_table_statement *CreateTableStatement
}

type CreateTableStatement struct {
	name    Token
	columns []ColumnDefinition
}

type ColumnDefinition struct {
	name      Token
	data_type Token
}

func Parse(tokens []Token) Ast {
	var statements []*Statement

	for index := 0; index < len(tokens); {
		statement, new_index, err := parse_create(tokens[index:], index)
		if err != nil {
			panic(err)
		}

		statements = append(statements, statement)
		index = new_index
	}

	return Ast{statements}
}

func next_token(tokens []Token, index int) (*Token, int) {
	if index >= len(tokens) {
		return nil, index
	}

	token := tokens[index]
	index++

	return &token, index
}

// CREATE TABLE table_name ([column_name column_type]);
func parse_create(tokens []Token, index int) (*Statement, int, error) {
	token, index := next_token(tokens, index)
	if token == nil || !token.is_keyword(KEYWORD_CREATE) {
		return nil, index, errors.New("")
	}

	token, index = next_token(tokens, index)
	if token == nil || !token.is_keyword(KEYWORD_TABLE) {
		return nil, index, errors.New("")
	}

	token, index = next_token(tokens, index)
	if token == nil || token._type != TOKEN_IDENTIFIER {
		return nil, index, errors.New("")
	}
	table_name_token := token

	token, index = next_token(tokens, index)
	if token == nil || !token.is_symbol(SYMBOL_LEFT_PAREN) {
		return nil, index, errors.New("")
	}

	columns, index, err := parse_column_definitions(tokens, index)
	if err != nil {
		return nil, index, errors.New("")
	}

	token, index = next_token(tokens, index)
	if token == nil || !token.is_symbol(SYMBOL_RIGHT_PAREN) {
		return nil, index, errors.New("")
	}

	token, index = next_token(tokens, index)
	if token == nil || !token.is_symbol(SYMBOL_SEMI_COLON) {
		return nil, index, errors.New("")
	}

	return &Statement{
		_type:                  CREATE_TABLE_STATEMENT,
		create_table_statement: &CreateTableStatement{*table_name_token, columns},
	}, index, nil
}

func parse_column_definitions(tokens []Token, index int) ([]ColumnDefinition, int, error) {
	var columns []ColumnDefinition

	var token *Token
	for {
		token, index = next_token(tokens, index)
		if token == nil || token._type != TOKEN_IDENTIFIER {
			return nil, index, errors.New("")
		}
		column_name_token := token

		token, index = next_token(tokens, index)
		if token == nil || !token.is_data_type() {
			return nil, index, errors.New("")
		}
		column_data_type := token

		columns = append(columns, ColumnDefinition{name: *column_name_token, data_type: *column_data_type})

		token, index = next_token(tokens, index)
		if token.is_symbol(SYMBOL_RIGHT_PAREN) {
			index--
			break
		}

		if !token.is_symbol(SYMBOL_COMMA) {
			return nil, index, errors.New("unable to parse column definitions - expected ')' or ','")
		}
	}

	return columns, index, nil
}
