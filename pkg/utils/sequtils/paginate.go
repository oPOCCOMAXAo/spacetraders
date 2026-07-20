package sequtils

import "context"

// Paginate wraps page/limit pagination in an iter.Seq2-style generator.
//
// It walks pages starting from startPage (default 1) with the given limit
// (default defaultLimit) until a page returns fewer than limit items or
// an empty page, yielding each element in turn. A fetch error is yielded
// in-stream and stops iteration.
func Paginate[T any](
	ctx context.Context,
	startPage, limit, defaultLimit int,
	fetch func(page, limit int) ([]*T, error),
) func(yield func(*T, error) bool) {
	page := startPage
	if page <= 0 {
		page = 1
	}

	pageLimit := limit
	if pageLimit <= 0 {
		pageLimit = defaultLimit
	}

	return func(yield func(*T, error) bool) {
		for ctx.Err() == nil {
			items, err := fetch(page, pageLimit)
			if err != nil {
				yield(nil, err)

				return
			}

			if len(items) == 0 {
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			if len(items) < pageLimit {
				return
			}

			page++
		}
	}
}

// ToPtrSlice returns a slice of pointers into the given items.
func ToPtrSlice[T any](items []T) []*T {
	res := make([]*T, len(items))
	for i := range items {
		res[i] = &items[i]
	}

	return res
}
