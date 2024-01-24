package ratelimit

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/auth"
	"github.com/hosseintrz/gaterun/pkg/cache/redis"
	"github.com/hosseintrz/gaterun/pkg/router"
)

func SlidingWindowLogMDFactory(cfg models.RateLimitConfig) router.MiddlewareFactory {
	sortedSetName := fmt.Sprintf("%s%d%d", cfg.Domain, cfg.Algorithm, (rand.Intn(1000) + 1000))
	return func(next router.Middleware) router.Middleware {
		return func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var userKey string

			userId := r.Header.Get(auth.HEADERS_CONSUMER_ID)
			ip := strings.Split(r.RemoteAddr, ":")[0]

			if cfg.Target == models.TargetID {
				userKey = userId
			} else {
				userKey = ip
			}

			redisCli := redis.GetClient(ctx)

			key := fmt.Sprintf("%s.%s", sortedSetName, userKey)
			now := time.Now()

			pipe := redisCli.TxPipeline()

			// Fetch all logs
			intervalStart := fmt.Sprintf("%d", now.Add(-cfg.Interval).UnixNano())
			pipe.ZRemRangeByScore(ctx, key, "-inf", intervalStart).Result()

			// Count the number of logs in the current interval
			pipe.ZCount(ctx, key, intervalStart, "+inf").Result()

			// Add the new log
			val := now.UnixNano()
			pipe.ZAdd(ctx, key, &goredis.Z{Score: float64(val), Member: val})

			// Execute the pipeline
			cmder, err := pipe.Exec(ctx)
			if err != nil {
				return
			}

			count := cmder[1].(*goredis.IntCmd).Val()
			if count >= int64(cfg.Threshold) {
				rw.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next(rw, r)
		}
	}
}
