package concurrent

import (
	"sync"
)

type Result[I any, T any] struct {
	Req I
	Res T
	Err error
}

type Fetcher[I, T any] func(I) (T, error)

func Fetch[I, T any](f Fetcher[I, T], ids []I) []Result[I, T] {
	results := make([]Result[I, T], 0, len(ids))
	ch := make(chan Result[I, T])
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		r := Result[I, T]{Req: id}
		go func() {
			res, err := f(id)
			if err != nil {
				r.Err = err
			} else {
				r.Res = res
			}
			ch <- r
		}()
	}
	go func() {
		for r := range ch {
			results = append(results, r)
			wg.Done()
		}
	}()
	wg.Wait()
	return results
}
