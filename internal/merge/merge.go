package merge

import (
	"reflect"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

var (
	ErrNonPointerDst        = errors.New("dst is not a pointer")
	ErrUnsupportedMerge     = errors.New("unsupported merge")
	ErrUnexpectedFailedCast = errors.New("unexpected failed cast")
)

func Merge(dst interface{}, sources ...interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.WithStack(ErrNonPointerDst)
	}

	dstPointedKind := reflect.Indirect(reflect.ValueOf(dst)).Elem().Kind()

	for _, src := range sources {
		srcKind := reflect.ValueOf(src).Kind()

		switch dstPointedKind {
		case reflect.Map:
			if srcKind != dstPointedKind {
				return errors.WithStack(unsupportedMergeError(dstPointedKind, srcKind))
			}

			if err := mergeMaps(dst, src); err != nil {
				return errors.WithStack(err)
			}

		case reflect.Slice:
			if srcKind != dstPointedKind {
				return errors.WithStack(unsupportedMergeError(dstPointedKind, srcKind))
			}

			if err := mergeSlices(dst, src); err != nil {
				return errors.WithStack(err)
			}

		default:
			return errors.WithStack(unsupportedMergeError(dstPointedKind, srcKind))
		}
	}

	return nil
}

func unsupportedMergeError(dstKind reflect.Kind, defaultsKind reflect.Kind) error {
	return errors.Wrapf(ErrUnsupportedMerge, "could not merge '%s' with defaults '%s'", dstKind, defaultsKind)
}

func mergeMaps(dst interface{}, defaults interface{}) error {
	dstMap, ok := reflect.Indirect(reflect.ValueOf(dst)).Elem().Interface().(map[string]interface{})
	if !ok {
		return errors.WithStack(ErrUnexpectedFailedCast)
	}

	defaultsMap, ok := defaults.(map[string]interface{})
	if !ok {
		return errors.WithStack(ErrUnexpectedFailedCast)
	}

	if err := mergo.Merge(&dstMap, defaultsMap, mergo.WithOverride); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func mergeSlices(dst interface{}, defaults interface{}) error {
	dstSlice, ok := reflect.Indirect(reflect.ValueOf(dst)).Elem().Interface().([]interface{})
	if !ok {
		return errors.WithStack(ErrUnexpectedFailedCast)
	}

	defaultsSlice, ok := defaults.([]interface{})
	if !ok {
		return errors.WithStack(ErrUnexpectedFailedCast)
	}

	if err := mergo.Merge(&dstSlice, defaultsSlice, mergo.WithOverride); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
