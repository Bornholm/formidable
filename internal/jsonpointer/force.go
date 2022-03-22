package jsonpointer

import (
	"strconv"

	"github.com/pkg/errors"
)

func force(doc interface{}, tokens []string, value interface{}) (interface{}, error) {
	if len(tokens) == 0 {
		return value, nil
	}

	currentToken := tokens[0]

	switch typedDoc := doc.(type) {

	case map[string]interface{}:
		nestedDoc, exists := typedDoc[currentToken]
		if !exists {
			if len(tokens) == 1 {
				typedDoc[currentToken] = value

				return typedDoc, nil
			}

			nextToken := tokens[1]
			if isArrayIndexToken(nextToken) {
				nestedDoc = make([]interface{}, 0)
			} else {
				nestedDoc = make(map[string]interface{})
			}
		}

		nestedDoc, err := force(nestedDoc, tokens[1:], value)
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
				for i := len(typedDoc); i <= int(index); i++ {
					typedDoc = append(typedDoc, nil)
				}
			}

			nestedDoc = typedDoc[index]
		}

		nestedDoc, err = force(nestedDoc, tokens[1:], value)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		typedDoc[index] = nestedDoc

		return typedDoc, nil

	default:
		overrideDoc := map[string]interface{}{}
		overrideDoc[currentToken] = value

		var nestedDoc interface{}

		if len(tokens) > 1 && isArrayIndexToken(tokens[1]) {
			nestedDoc = make([]interface{}, 0)
		} else {
			nestedDoc = make(map[string]interface{})
		}

		nestedDoc, err := force(nestedDoc, tokens[1:], value)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		overrideDoc[currentToken] = nestedDoc

		return overrideDoc, nil
	}
}

func isArrayIndexToken(token string) bool {
	if token == NonExistentMemberToken {
		return true
	}

	if _, err := strconv.ParseUint(token, 10, 64); err != nil {
		return false
	}

	return true
}
