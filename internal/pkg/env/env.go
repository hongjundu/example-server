package env

var Env struct {
	LogLevel int
	Port     int
	LogPath  string

	RedisMode       string
	RedisHost       string
	RedisMasterName string
	RedisDb         int
	RedisPassword   string

	MySqlHost     string
	MySqlDb       string
	MySqlUser     string
	MySqlPassword string
}
