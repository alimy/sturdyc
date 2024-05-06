package sturdyc

import (
	"context"
	"errors"
)

func wrap[T, V any](fetchFn FetchFn[V]) FetchFn[T] {
	return func(ctx context.Context) (T, error) {
		res, err := fetchFn(ctx)
		if err != nil {
			var zero T
			return zero, err
		}
		if val, ok := any(res).(T); ok {
			return val, nil
		}
		var zero T
		return zero, errors.New("invalid response type")
	}
}

func unwrap[T, V any](val T, err error) (V, error) {
	v, ok := any(val).(V)
	if !ok {
		return v, errors.New("invalid response type")
	}
	return v, err
}

func wrapBatch[T, V any](fetchFn BatchFetchFn[V]) BatchFetchFn[T] {
	return func(ctx context.Context, ids []string) (map[string]T, error) {
		resV, err := fetchFn(ctx, ids)
		if err != nil {
			return map[string]T{}, err
		}

		resT := make(map[string]T, len(resV))
		for id, v := range resV {
			if val, ok := any(v).(T); ok {
				resT[id] = val
			}
		}

		return resT, nil
	}
}

func unwrapBatch[T, V any](values map[string]T, err error) (map[string]V, error) {
	vals := make(map[string]V, len(values))
	for id, v := range values {
		if val, ok := any(v).(V); ok {
			vals[id] = val
		}
	}
	return vals, err
}
