package auth

import (
	"concert_pre-poster/internal/cache/redisCache"
	"concert_pre-poster/pkg/util"
	"context"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"

	"github.com/spf13/viper"
)

var redisClient *redis.Client

func init() {

	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// Получаем значения из конфигурации
	host := viper.GetString("redis.host")
	password := viper.GetString("redis.password")
	count := viper.GetInt("redis.countOfDataBase")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,     // Адрес вашего Redis сервера
		Password: password, // Пароль Redis сервера, если есть
		DB:       count,    // Номер базы данных Redis
	})

	ctx := context.Background()
	// Проверка соединения с Redis
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
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
			_, err = session.GetValue(context.Background(), sessionID.Value)

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

	sessionID := util.RandStringRunes(32)

	session := redisCache.NewRedisCache(redisClient)

	err := session.SetValue(context.Background(), sessionID, inputLogin, 0)

	if err != nil {
		log.Fatal(err)
	}

	cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)

}
