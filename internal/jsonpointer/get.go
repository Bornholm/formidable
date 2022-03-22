package jsonpointer

import (
	"strconv"

	"github.com/pkg/errors"
)

func get(doc interface{}, tokens []string) (interface{}, error) {
	if len(tokens) == 0 {
		return doc, nil
	}

	currentToken := tokens[0]

	if doc == nil {
		return nil, errors.Wrapf(ErrNotFound, "pointer '%s' not found on document", tokensToString(tokens))
	}

	switch typedDoc := doc.(type) {
	case map[string]interface{}:
		value, exists := typedDoc[currentToken]
		if !exists {
			return nil, errors.Wrapf(ErrNotFound, "pointer '%s' not found on document", tokensToString(tokens))
		}

		value, err := get(value, tokens[1:])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return value, nil

	case []interface{}:
		if currentToken == NonExistentMemberToken {
			return nil, errors.WithStack(ErrOutOfBounds)
		}

		index, err := strconv.ParseUint(currentToken, 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if len(typedDoc) <= int(index) {
			return nil, errors.WithStack(ErrOutOfBounds)
		}

		value := typedDoc[index]

		value, err = get(value, tokens[1:])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return value, nil

	default:
		return nil, errors.Wrapf(ErrUnexpectedType, "unexpected type '%T'", typedDoc)
	}
}
