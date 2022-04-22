package main

import (
	"fmt"
	"os"

	yandexAPI "github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/config"
	fileLogger "github.com/de1phin/music-transfer/internal/log/file_logger"
	"github.com/de1phin/music-transfer/internal/service/yandex"
	"github.com/de1phin/music-transfer/internal/storage/postgres"
)

func main() {

	config := config.NewConfig("./config", "config", "yaml")

	pDataSourceName := config.GetPosgresDataSourceName()

	psql, err := postgres.NewPostgresDatabase(pDataSourceName)
	if err != nil {
		panic(err)
	}
	table := postgres.NewTable[int64, yandexAPI.Credentials](psql, "Yandex", "id")

	logger, err := fileLogger.NewFileLogger("./log/t.log")
	if err != nil {
		logger.Log(err)
		os.Exit(1)
	}

	api := yandexAPI.NewYandexAPI(logger, config.GetYandexMagicToken())
	service := yandex.NewYandexService(api, table, logger)
	api.BindOnGetCredentials(service.OnGetCredentials)

	fmt.Println(api.GetAuthURL(123))
	fmt.Scanln()
	c, _ := table.Get(123)
	fmt.Println(api.GetMe(c))

	for {

	}
}
