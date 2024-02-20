package auth

import (
	"concert_pre-poster/internal/cache/redisCache"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"time"
)


var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес вашего Redis сервера
		Password: "",               // Пароль Redis сервера, если есть
		DB:       0,                // Номер базы данных Redis
	})

	ctx := context.Background()
	// Проверка соединения с Redis
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/login" && r.URL.Path != "/get_cookie" {
			sessionID, err := r.Cookie("session_id")

			// Если куки не существует или ошибка, перенаправляем пользователя на страницу аутентификации
			if err == http.ErrNoCookie {
				http.ServeFile(w, r, "./templates/login.html") 
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			session := redisCache.NewRedisCache(redisClient)
			username, err := session.GetValue(context.Background(), sessionID.Value)

			fmt.Println(username)
			// Если сессия не найдена в Redis, перенаправляем пользователя на страницу аутентификации
			if err != nil {
				http.ServeFile(w, r, "./templates/login.html") 
				return
			}
		}

		// Если пользователь авторизован, продолжаем обработку запроса
		next.ServeHTTP(w, r)
	})
}

func GetCookie(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	inputLogin := r.Form["login"][0]
	expiration := time.Now().Add(365 * 24 * time.Hour)

	sessionID := RandStringRunes(32)

	session := redisCache.NewRedisCache(redisClient)

	err := session.SetValue(context.Background(), sessionID, inputLogin, 0)

	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)

}

// PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
