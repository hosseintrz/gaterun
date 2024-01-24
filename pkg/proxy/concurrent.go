package proxy

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/sirupsen/logrus"
)

var (
	ErrConcurrentMiddleware = errors.New("concurrent middleware added while concurrent_calls set to 1")
	ErrNilResponse          = errors.New("response was nil")
)

func ConcurrentMiddleware(cfg *models.BackendConfig) Middleware {
	return func(next ...Proxy) Proxy {
		if cfg.ConcurrentCalls == 1 {
			panic(ErrConcurrentMiddleware)
		}

		timeout := time.Duration(cfg.Timeout)
		cnt := cfg.ConcurrentCalls

		return func(ctx context.Context, req *Request) (response *Response, err error) {
			defer next[0](ctx, req)

			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			results := make(chan *Response, cnt)
			failures := make(chan error, cnt)

			wg := sync.WaitGroup{}

			for i := int32(0); i < cfg.ConcurrentCalls; i++ {
				go func() {
					wg.Add(1)
					defer wg.Done()
					processConcurrentCall(ctxWithTimeout, next[0], req, results, failures)
				}()
			}

			done := make(chan bool)

			go func() {
				wg.Wait()
				done <- true
			}()

			for {
				select {
				case response = <-results:
					return
				case err = <-failures:
					logrus.Error(err)
				case <-done:
					break
				}
			}
		}
	}
}

func processConcurrentCall(ctx context.Context, proxy Proxy, req *Request, response chan *Response, failures chan error) {
	localCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := proxy(localCtx, req)
	if err != nil {
		failures <- err
		return
	}

	if res == nil {
		err = ErrNilResponse
		failures <- err
		return
	}

	select {
	case <-ctx.Done():
		failures <- ctx.Err()
	case response <- res:
	}
}
