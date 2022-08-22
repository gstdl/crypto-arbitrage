package config

type configs map[string]interface{}

type lunoConfig struct {
	Id     string       `mapstructure:"key_id"`
	Secret string       `mapstructure:"secret"`
	Assets LunoAssetMap `mapstructure:"assets"`
	Pairs  []string
}

type LunoAssetMap map[string]LunoAsset

type LunoAsset struct {
	Ticker    string `mapstructure:"ticker"`
	AddressID string `mapstructure:"spot_wallet_address"`
	Threshold string `mapstructure:"threshold"`
}
