package config

import (
	"fmt"
	"strings"

	"github.com/luno/luno-go/decimal"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.SetConfigFile(viper.GetString("CONFIG_PATH"))
	err = viper.MergeInConfig()
	if err != nil {
		panic(err)
	}
}

func GetEnvironmentValue() (env string) {
	env = viper.GetString("ENVIRONMENT")
	return
}

func GetLunoConfig() (lc lunoConfig) {

	if err := viper.UnmarshalKey("luno", &lc); err != nil {
		panic(err)
	}

	return
}

func GetLunoCredentials() (id, secret string) {
	cfg := GetLunoConfig()
	id = cfg.Id
	secret = cfg.Secret
	return
}

func (assets LunoAssetMap) GetThreshold(asset string) (threshold decimal.Decimal, err error) {
	if a, ok := assets[strings.ToLower(asset)]; !ok {
		err = fmt.Errorf("can't convert %s: asset not found in config file", asset)
	} else {
		threshold, err = decimal.NewFromString(a.Threshold)
	}
	return
}
