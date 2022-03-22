package jsonpointer

import (
	"strconv"

	"github.com/pkg/errors"
)

func del(doc interface{}, tokens []string) (interface{}, error) {
	currentToken := tokens[0]

	switch typedDoc := doc.(type) {

	case map[string]interface{}:
		nestedDoc, exists := typedDoc[currentToken]
		if !exists {
			return doc, nil
		}

		if len(tokens) == 1 {
			delete(typedDoc, currentToken)

			return typedDoc, nil
		}

		nestedDoc, err := del(nestedDoc, tokens[1:])
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
			index = uint64(len(typedDoc) - 1)
		} else {
			index, err = strconv.ParseUint(currentToken, 10, 64)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			if len(typedDoc) <= int(index) {
				return typedDoc, nil
			}
		}

		if len(tokens) == 1 {
			typedDoc = append(typedDoc[:index], typedDoc[index+1:]...)

			return typedDoc, nil
		}

		nestedDoc, err = del(nestedDoc, tokens[1:])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		typedDoc[index] = nestedDoc

		return typedDoc, nil

	default:
		return typedDoc, nil
	}
}
