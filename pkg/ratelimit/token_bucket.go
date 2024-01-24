package ratelimit

import (
	"net/http"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/router"
)

func TokenBucketMiddlewareWithCfg(cfg models.RateLimitConfig) router.MiddlewareFactory {
	return func(next router.Middleware) router.Middleware {
		return func(rw http.ResponseWriter, r *http.Request) {
			// ctx := r.Context()
			// var key string

			// userId := r.Header.Get(auth.HEADERS_CONSUMER_ID)
			// ip := r.RemoteAddr

			// if cfg.Target == models.TargetID {
			// 	key = userId
			// } else {
			// 	key = ip
			// }

			// redisCli := redis.GetClient(ctx)
			// value, err := redisCli.Get(ctx, key).Result()
			// if err != nil {
			// 	rw.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// count, err := strconv.Atoi(value)
			// if err != nil {
			// 	rw.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }

			// if count+1 > cfg.RequestsPerUnit {
			// 	rw.WriteHeader(http.StatusTooManyRequests)
			// 	fmt.Fprintf(rw, "too many requests")
			// 	return
			// } else {
			// 	err = redisCli.Incr(ctx, key).Err()
			// 	if err != nil {
			// 		rw.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}
			// }

			next(rw, r)
		}
	}
}
