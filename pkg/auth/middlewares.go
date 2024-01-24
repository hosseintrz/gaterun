package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hosseintrz/gaterun/pkg/cache/redis"
	"github.com/hosseintrz/gaterun/pkg/router"
)

func APIKeyMiddleware(next router.Middleware) router.Middleware {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiKeyStr := ""

		apikeys, ok := r.Header["Gaterun-Api-Key"]
		if ok {
			apiKeyStr = apikeys[0]
		} else if r.URL.Query().Has("apikey") {
			apiKeyStr = r.URL.Query().Get("apikey")
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, "API_KEY not provided")
			return
		}

		redisClient := redis.GetClient(ctx)
		apiKeyJson, err := redisClient.Get(ctx, apiKeyStr).Result()
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, "API_KEY is not valid")
			return
		}

		var apiKey ApiKey
		err = json.Unmarshal([]byte(apiKeyJson), &apiKey)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, "API_KEY is not valid")
			return
		}

		rw.Header().Del(HEADERS_ANONYMOUS)
		rw.Header().Set(HEADERS_CONSUMER_ID, fmt.Sprintf("%d", apiKey.ConsumerId))
		rw.Header().Set(HEADERS_CONSUMER_CUSTOM_ID, string(apiKey.ConsumerCustomId))
		rw.Header().Set(HEADERS_CONSUMER_USERNAME, string(apiKey.ConsumerUsername))

		next(rw, r)
	}
}

// func APIKeyMiddleware(rw http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	apiKeyStr := ""

// 	apikeys, ok := r.Header["Gaterun-Api-Key"]
// 	if ok {
// 		apiKeyStr = apikeys[0]
// 	} else if r.URL.Query().Has("apikey") {
// 		apiKeyStr = r.URL.Query().Get("apikey")
// 	} else {
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(rw, "API_KEY not provided")
// 		return
// 	}

// 	redisClient := redis.GetClient(ctx)
// 	apiKeyJson, err := redisClient.Get(ctx, apiKeyStr).Result()
// 	if err != nil {
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(rw, "API_KEY is not valid")
// 		return
// 	}

// 	var apiKey ApiKey
// 	err = json.Unmarshal([]byte(apiKeyJson), &apiKey)
// 	if err != nil {
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(rw, "API_KEY is not valid")
// 		return
// 	}

// 	rw.Header().Del(HEADERS_ANONYMOUS)
// 	rw.Header().Set(HEADERS_CONSUMER_ID, fmt.Sprintf("%d", apiKey.ConsumerId))
// 	rw.Header().Set(HEADERS_CONSUMER_CUSTOM_ID, string(apiKey.ConsumerCustomId))
// 	rw.Header().Set(HEADERS_CONSUMER_USERNAME, string(apiKey.ConsumerUsername))

// }

func BasicMiddleware(next router.Middleware) router.Middleware {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("basic auth not implemented")
		next(rw, r)
	}
}

func JWTMiddleware(next router.Middleware) router.Middleware {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("jwt not implemented")
		next(rw, r)
	}
}
