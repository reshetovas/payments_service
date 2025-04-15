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
	"payments_service/logger"
	"payments_service/routes"
	"payments_service/services"
	"payments_service/storage"

	"github.com/gorilla/mux"
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
	defer cancel() // when func main is over - all go riutines will finish

	// Обработчик сигналов для завершения
	sigChan := make(chan os.Signal, 1) //create channel
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
	fmt.Println(cacheType)

	// Инициализация зависимостей
	paymentStorage := storage.NewPaymentStorage(db)
	bonusStorage := storage.NewBonusStorage(db)
	cachePayment := storage.NewPaymentCache(paymentStorage, cacheType, redisClient)
	//to do bonuses
	currencyService := currency_service.NewCurrencyAPI("https://api.exchangerate-api.com/v4")
	paymentService := services.NewPaymentService(cachePayment, currencyService)
	bonusService := services.NewBonusService(bonusStorage)
	parse := services.NewParseService(paymentStorage)
	paymentHandler := handlers.NewPaymentHandler(paymentService, parse)
	bonusHandler := handlers.NewBonusHandler(bonusService)
	bckgrnd_serv := background_service.NewBackgroundService(paymentStorage)
	paymentRoutes := routes.NewPaymentRoutes(paymentHandler)
	bonusRoutes := routes.NewBonusRoutes(bonusHandler)

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
	mainRouter := mux.NewRouter()

	paymentRouter := paymentRoutes.PaymentRouter()
	bonusRouter := bonusRoutes.BonusRouter()

	mainRouter.PathPrefix("/payment").Handler(paymentRouter)
	mainRouter.PathPrefix("/payments").Handler(paymentRouter)
	mainRouter.PathPrefix("/bonus").Handler(bonusRouter)
	mainRouter.PathPrefix("/bonuses").Handler(bonusRouter)

	// Настройка HTTP-сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: mainRouter,
	}

	// Запуск HTTP-сервера в горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("The server is running on port: 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Error http server")
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
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Info().Msg("Приложение завершено.")
}
