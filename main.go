package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	background_service "payments_service/background_services"
	"payments_service/config"
	"payments_service/external/currency_service"
	"payments_service/handlers"
	"payments_service/internal/ws"
	"payments_service/logger"
	"payments_service/middleware"
	"payments_service/routes"
	"payments_service/services"
	"payments_service/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func main() {
	//start parsing config
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v", err)
	}

	//logging setup
	logger.InitLogger(config.LogLevel)
	// logInstance := logger.NewLoggerStruct(config.LogLevel)

	// Контекст с отменой для graceful shutdown -
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // when func main is over - all go routines will finish

	// Обработчик сигналов для завершения
	sigChan := make(chan os.Signal, 1) //create chanel
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Подключение к БД
	db, err := sql.Open("sqlite3", "./payments.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB sql.Open")
	}
	defer db.Close()

	// Применение миграций
	if err := goose.Up(db, "./storage/migrations"); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	// DB check
	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("error db ping")
	}

	//Redis init
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatal().Err(err).Msg("Redis connection failed")
	}

	cacheType := config.CacheType
	log.Info().Msgf("cache type: %s", cacheType)
	log.Info().Msgf("Config: %v", config)

	//token secret
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	//auth middleware

	// Инициализация зависимостей
	paymentStorage := storage.NewPaymentStorage(db)
	bonusStorage := storage.NewBonusStorage(db)
	userStorage := storage.NewUserStorage(db)
	wsHub := ws.NewHub()
	go wsHub.Run()

	cachePayment := storage.NewPaymentCache(paymentStorage, cacheType, redisClient)
	cacheBonus := storage.NewBonusCache(bonusStorage, cacheType, redisClient)

	currencyService := currency_service.NewCurrencyAPI("https://api.exchangerate-api.com/v4")
	paymentService := services.NewPaymentService(cachePayment, currencyService, wsHub)
	bonusService := services.NewBonusService(cacheBonus)
	parse := services.NewParseService(paymentStorage)
	token := services.NewTokenStruct(jwtSecret)
	userService := services.NewUserService(userStorage, token)
	authorization := middleware.NewAuthMiddleware(token) //authorization middleware

	paymentHandler := handlers.NewPaymentHandler(paymentService, parse)
	bonusHandler := handlers.NewBonusHandler(bonusService)
	userHandler := handlers.NewUserHandler(userService)
	wsHandler := handlers.NewWSHandler(wsHub)

	bckgrnd_serv := background_service.NewBackgroundService(paymentStorage)
	paymentRoutes := routes.NewPaymentRoutes(paymentHandler, authorization)
	bonusRoutes := routes.NewBonusRoutes(bonusHandler)
	userRoutes := routes.NewUserRoutes(userHandler)
	wsRouter := routes.NewWSRoutes(wsHandler, authorization)

	// Запуск мониторинга директории в отдельной горутине
	dirName := "./files"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // is always executed at the end of the function
		paymentHandler.ProcessExistingFiles(ctx, dirName)
	}()

	//state machine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := bckgrnd_serv.CheckStatuses(ctx); err != nil {
			log.Fatal().Err(err).Msg("Background service error in function CheckStatuses")
		}
	}()

	// Настройка маршрутов
	mainHttpRouter := routes.MainRouter(paymentRoutes, bonusRoutes, userRoutes)
	mainWSRouter := routes.MainWSRouter(wsRouter)
	fmt.Printf("mainWSRouter: %+v\n", mainWSRouter)

	// Настройка HTTP-сервера
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mainHttpRouter,
	}

	// Запуск HTTP-сервера в горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("The server is running on port: 8080...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Error HTTP server")
		}
	}()

	//настройка WS сервера
	wsServer := &http.Server{
		Addr:    ":8081",
		Handler: mainWSRouter,
	}

	// Запуск WS-сервера в горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("The server is running on port: 8081...")
		if err := wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Error WebSocket server")
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Info().Msg("Получен сигнал завершения, останавливаем сервер...")

	// Остановка контекста и сервера
	cancel() //finish all go routines, that use context

	// Graceful shutdown для HTTP-сервера
	/* Создаётся новый контекст shutdownCtx с логикой:
	Когда истекает 5 секунд с момента создания контекста,
	статус контекста автоматически меняется на "cancelled".
	Это означает, что все горутины и операции,
	которые следят за этим контекстом, должны завершиться или отмениться.*/
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	//server.Shutdown() - for finish http server
	/* server.Shutdown(shutdownCtx) =
	Cервер будет пытаться завершить свою работу в течение 5 секунд.
	Если сервер не успевает завершить все запросы,
	он завершится принудительно по истечении этого времени */
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}
	if err := wsServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Info().Msg("Приложение завершено.")
}
