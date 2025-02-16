package main

import (
	handler "avito-internship/internal/handler"
	"avito-internship/internal/repository"
	"avito-internship/internal/repository/postgres"
	"avito-internship/internal/server"
	"avito-internship/internal/services"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("произошла ошибка при инициализации конфигурационного файла: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("произошла ошибка при загрузке переменных окружения: %s", err.Error())
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("произошла ошибка при инициализации базы данных: %s", err)
	}

	repos := repository.NewRepository(db)
	services := services.NewServices(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("произошла ошибка при запуске http-сервера: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("произошла ошибка при выключении сервера: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("произошла ошибка при закрытии соединения с базой данных: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
