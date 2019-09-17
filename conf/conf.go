package conf

// RedisConfig ...
type RedisConfig struct {
	Host     string
	Password string
	Port     int
	DB       int
}

// DbConfig ...
type DbConfig struct {
	Driver  string
	Source  string
	LogMode bool
}

// LogConfig ...
type LogConfig struct {
	Level      int
	EnableFile bool
	LogPath    string
}

// ServerConfig ...
type ServerConfig struct {
	Port                   int
	TraceEnabled           bool
	CorsEnabled            bool
	AllowOrigins           []string
	JwtISS                 string
	JwtScrect              string
	HTTPServerEnabled      bool
	GRPCServerEnabled      bool
	WebsocketServerEnabled bool
}

// BaseConfig ...
type BaseConfig struct {
	// RedisConf ...
	RedisConf *RedisConfig
	// DBConf ...
	DBConf *DbConfig
	// LogConf ...
	LogConf *LogConfig
	// ServerConf ...
	ServerConf *ServerConfig
	// GlobalConf ...
	GlobalConf *GlobalConfig
}

// GlobalConfig ...
type GlobalConfig struct {
}
