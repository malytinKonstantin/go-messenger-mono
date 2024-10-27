package database

import (
	"github.com/malytinKonstantin/go-messenger-mono/friendship-service/internal/models"
	"github.com/mindstand/gogm/v2"
	"github.com/spf13/viper"
)

var Gogm *gogm.Gogm

func InitGogm() error {
	config := gogm.Config{
		IndexStrategy: gogm.IGNORE_INDEX,
		Username:      viper.GetString("NEO4J_USER"),
		Password:      viper.GetString("NEO4J_PASSWORD"),
		Host:          viper.GetString("DATABASE_HOST"),
		Port:          viper.GetString("DATABASE_PORT"),
		Protocol:      "bolt",
		PoolSize:      50,
		IsCluster:     false,
		LogLevel:      "DEBUG",
	}

	var err error
	Gogm, err = gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &models.User{}, &models.FriendRequest{})
	if err != nil {
		return err
	}

	return nil
}
