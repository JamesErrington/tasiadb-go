package main

import (
	"strings"
	"testing"
)

func (token *Token) test_token_equals(other Token) bool {
	return token._type == other._type && token.value == other.value && token.location.column == other.location.column
}

func TestLexSymbol(t *testing.T) {
	symbols := []string{string(SYMBOL_SEMI_COLON), string(SYMBOL_COMMA), string(SYMBOL_LEFT_PAREN), string(SYMBOL_RIGHT_PAREN)}
	source := strings.Join(symbols, " ")
	tokens := Lex(source)

	if len(tokens) != len(symbols) {
		t.Errorf("Expected %d tokens, got %d: %v", len(symbols), len(tokens), tokens)
	}

	for i, symbol := range symbols {
		token := tokens[i]
		expected := Token{TOKEN_SYMBOL, symbol, Location{0, uint(2 * i)}}
		if !token.test_token_equals(expected) {
			t.Errorf("Expected %v, got %v", expected, token)
		}
	}
}

func TestLexText(t *testing.T) {
	source := "'abc ABC 123 ()\";_*'"
	tokens := Lex(source)

	if len(tokens) != 1 {
		t.Errorf("Expected 1 token, got %d: %v", len(tokens), tokens)
	}

	expected := Token{TOKEN_TEXT, "abc ABC 123 ()\";_*", Location{0, 0}}
	if !tokens[0].test_token_equals(expected) {
		t.Errorf("Expected %v, got %v", expected, tokens[0])
	}
}

func TestLexNumber(t *testing.T) {
	source := "42"
	tokens := Lex(source)

	if len(tokens) != 1 {
		t.Errorf("Expected 1 token, got %d: %v", len(tokens), tokens)
	}

	expected := Token{TOKEN_INT, "42", Location{0, 0}}
	if !tokens[0].test_token_equals(expected) {
		t.Errorf("Expected %v, got %v", expected, tokens[0])
	}
}
