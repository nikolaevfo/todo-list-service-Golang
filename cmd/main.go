package main

import (
	_ "github.com/lib/pq"
	"os"
	todo "to-do-list"
	"to-do-list/pkg/handler"
	"to-do-list/pkg/repository"
	"to-do-list/pkg/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// информация для swagger
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
	// интерфейс handler из пакет http, поэтому его можно передать,
	// как аргумент сервера
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatal("error occured while running http server: %s", err.Error())
	}
}

// инициализируем конфигурационные файлы с помощью viper
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
