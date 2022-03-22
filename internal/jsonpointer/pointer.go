package jsonpointer

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	TokenSeparator         = "/"
	NonExistentMemberToken = "-"
)

type Pointer struct {
	tokens []string
}

func (p *Pointer) Get(doc interface{}) (interface{}, error) {
	value, err := get(doc, p.tokens)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return value, nil
}

func (p *Pointer) Set(doc interface{}, value interface{}) (interface{}, error) {
	doc, err := set(doc, p.tokens, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}

func (p *Pointer) Force(doc interface{}, value interface{}) (interface{}, error) {
	doc, err := force(doc, p.tokens, value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}

func (p *Pointer) Delete(doc interface{}) (interface{}, error) {
	doc, err := del(doc, p.tokens)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return doc, nil
}

func New(raw string) *Pointer {
	tokens := decodeTokens(raw)

	return &Pointer{tokens}
}

func tokensToString(tokens []string) string {
	escapedTokens := make([]string, 0)

	for _, t := range tokens {
		escapedTokens = append(escapedTokens, escapeToken(t))
	}

	return TokenSeparator + strings.Join(escapedTokens, TokenSeparator)
}

func escapeToken(token string) string {
	token = strings.ReplaceAll(token, "/", "~1")
	token = strings.ReplaceAll(token, "~", "~0")

	return token
}

func unescapeToken(token string) string {
	token = strings.ReplaceAll(token, "~1", "/")
	token = strings.ReplaceAll(token, "~0", "~")

	return token
}

func decodeTokens(raw string) []string {
	tokens := strings.Split(raw, TokenSeparator)

	for i, t := range tokens {
		tokens[i] = unescapeToken(t)
	}

	return tokens[1:]
}
