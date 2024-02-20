package auth

import (
	"concert_pre-poster/internal/cache/redisCache"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var loginFormTmpl = `
  <html>
	<body>
	<form action="/get_cookie" method="post">
	  Login: <input type="text" name="login">
	  Password: <input type="password" name="password">
	  <input type="submit" value="Login">
	</form>
	</body>
  </html>
  `

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

var sessions = map[string]string{}

func CookieAuth(w http.ResponseWriter, r *http.Request) {

	sessionID, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		w.Write([]byte(loginFormTmpl))
		return
	} else if err != nil {
		PanicOnErr(err)
	}

	session := redisCache.NewRedisCache(redisClient)

	username, err := session.GetValue(context.Background(), sessionID.Value)

	if err != nil {
		fmt.Fprint(w, "Session not found")
	} else {
		fmt.Fprint(w, "Welcome, "+username)
	}
}

func GetCookie(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	inputLogin := r.Form["login"][0]
	expiration := time.Now().Add(365 * 24 * time.Hour)


	sessionID := RandStringRunes(32)
	sessions[sessionID] = inputLogin

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
