package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	todo "to-do-list"
	"to-do-list/pkg/handler"
	"to-do-list/pkg/repository"
	"to-do-list/pkg/service"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// аннотация для swagger
// @title           Todo Service API
// @version         1.0
// @description     API Server for TodoList App

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// устанавливаем формат вывода для логгера
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// инициализация кофигурационного файла с помощью viper
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initalizing configs: %s", err.Error())
	}

	// применяется godotenv для чтения пароля из файла .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// инициализируем базу данных с помощью метода из модуля repository
	// используем конфигурациооный файл с помощью viper
	// используем переменную окруженя для получения пароля с помощью godotenv
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatal("failed to initialize db: %s", err.Error())
	}

	// Создаем экземпляры основных объектов и объявляем зависимости в нужном порядке
	// repos зависит от базы данных
	// services зависит от repos
	// handlers зависит от services
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// инициализируется экземпляр сервиса
	srv := new(todo.Server)

	// запускаем сервис, используем порт из config/config/yml
	// с помощью пакета viper
	// для обработки спользуем модуль handlers, метод InitRoutes.
	// Данный метод возвращает указатель на gin.Engine, который реализует
	// интерфейс handler из пакет http, поэтому его можно передать,как аргумент сервера
	//
	// запускаем сервис в горутине, чтобы можно было потом плавно завершить его работу
	// плавное завершение работы приложения должно гарантировать,
	// что мы перестали принимать запросы, но при этом закончили
	// все текущие обработки запросов и операции в БД
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatal("error occured while running http server: %s", err.Error())
		}
	}()

	// чтобы функция main не прекращала свою работу сразу, добавим блокировку с помощью канала
	quit := make(chan os.Signal, 1)

	// запись в канал происходит когда процесс приложения получит сигнал от системы
	// до этого момента приложение будет работать
	// SIGTERM - сигнал, применяемый в POSIX-системах для запроса завершения процесса
	// SIGINT - сигнал, применяемый в POSIX-системах для остановки процесса пользователем с терминала
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// читаем из канала, блокирующего выполнение главной горутины
	<-quit

	// при завершении работы остановим сервер
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	// закрываем соединение с БД
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

// инициализируем конфигурационные файлы с помощью viper
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
