package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Service struct {
	redis *redis.Client
}

var (
	redisAddr = "localhost:6379"
	httpPort  = ":8080"
)

func main() {

	redis := newRedisClient()
	service := Service{
		redis: redis,
	}

	// start http server
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.GET("/", service.getRequest)
	router.Run(httpPort)

}

func (service *Service) getRequest(ctx *gin.Context) {

	queryX, ok := ctx.GetQuery("x")
	if !ok {
		fmt.Println("X is empty")
		ctx.String(http.StatusBadRequest, "X is empty")
		return
	}
	x, err := strconv.Atoi(queryX)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusBadRequest, "X is not int")
		return
	}
	if x == 0 {
		fmt.Println("X is empty")
		ctx.String(http.StatusBadRequest, "X is empty")
		return
	}

	queryY, ok := ctx.GetQuery("y")
	if !ok {
		fmt.Println("Y is empty")
		ctx.String(http.StatusBadRequest, "Y is empty")
		return
	}
	y, err := strconv.Atoi(queryY)
	if err != nil {
		fmt.Println(err)
		ctx.String(http.StatusBadRequest, "Y is not int")
		return
	}

	if y <= x {
		fmt.Println("X must be less than Y")
		ctx.String(http.StatusBadRequest, "X must be less than Y")
		return
	}

	fmt.Println(x, y)

	response := service.fibonachi(x, y)

	ctx.String(http.StatusOK, fmt.Sprintln(response))
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ping, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(ping) < 0 {
		fmt.Println(err)
		return nil
	}
	return client
}

func (service *Service) fibonachi(x, y int) []int {

	redisSort := &redis.Sort{
		Order: "asc",
	}

	redisSlice := service.redis.Sort("fibo", redisSort)
	slice := redisSlice.Val()

	result := make([]int, 0)

	if len(slice) <= 0 {
		service.redis.LPush("fibo", 0)
		service.redis.LPush("fibo", 1)
		slice = []string{"0", "1"}
	}

	for _, value := range slice {
		v, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, v)
	}

	if len(result) < y {
		for i := len(result); i < y; i++ {
			elems := result[i-2] + result[i-1]
			result = append(result, elems)
			service.redis.LPush("fibo", elems)
		}
	}

	return result[x-1 : y]
}
