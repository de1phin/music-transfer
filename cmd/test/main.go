package main

import (
	"fmt"
	"os"

	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/log"
	fileLogger "github.com/de1phin/music-transfer/internal/log/file_logger"
)

var (
	logger log.Logger
	api    *yandex.YandexAPI
)

func OnGetCredentials(userID int64, credentials yandex.Credentials) {
	user, err := api.GetMe(credentials)
	if err != nil {
		logger.Log(err)
	}

	logger.Log("User =", user)
}

func main() {
	var err error
	logger, err = fileLogger.NewFileLogger("./log/t.log")
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	api = yandex.NewYandexAPI(logger)
	api.BindOnGetCredentials(OnGetCredentials)

	fmt.Println(api.GetAuthURL(123))

	for {

	}
}
