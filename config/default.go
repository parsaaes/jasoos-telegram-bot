package config

// Default return default configuration
// nolint: gomnd
func Default() Config {
	return Config{
		Token: "AwesomeBotToken",
		Words: []string{
			"Hello",
			"Home",
			"Raha",
			"Room",
		},
	}
}
