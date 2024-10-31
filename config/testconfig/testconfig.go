package testconfig

import "example.com/example/config"

var testConf *config.Config = &config.Config{
	Env: "test",
	DB: config.DatabseConfig{
		DatabaseUri: "file::memory:?cache=shared",
	},
	Host: "example.com",
	Port: 8888,
}

func ReloadTestConfig() *config.Config {
	return config.ReloadTestConfig(testConf)
}
