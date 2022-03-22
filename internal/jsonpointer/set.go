package jsonpointer

import (
	"strconv"

	"github.com/pkg/errors"
)

func set(doc interface{}, tokens []string, value interface{}) (interface{}, error) {
	if len(tokens) == 0 {
		return value, nil
	}

	currentToken := tokens[0]

	switch typedDoc := doc.(type) {
	case map[string]interface{}:
		nestedDoc, exists := typedDoc[currentToken]
		if !exists {
			return nil, errors.Wrapf(ErrNotFound, "pointer '%s' not found on document", tokensToString(tokens))
		}

		nestedDoc, err := set(nestedDoc, tokens[1:], value)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		typedDoc[currentToken] = nestedDoc

		return typedDoc, nil

	case []interface{}:
		var (
			index     uint64
			nestedDoc interface{}
			err       error
		)

		if currentToken == NonExistentMemberToken {
			typedDoc = append(typedDoc, value)
			index = uint64(len(typedDoc) - 1)
		} else {
			index, err = strconv.ParseUint(currentToken, 10, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if len(typedDoc) <= int(index) {
				return nil, errors.WithStack(ErrOutOfBounds)
			}

			nestedDoc = typedDoc[index]
		}

		nestedDoc, err = set(nestedDoc, tokens[1:], value)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		typedDoc[index] = nestedDoc

		return typedDoc, nil

	default:
		return nil, errors.Wrapf(ErrUnexpectedType, "unexpected type '%T'", typedDoc)
	}
}
