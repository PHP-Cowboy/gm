package config

type ServerConfig struct {
	Mode       string      `mapstructure:"mode" json:"mode"`
	Port       int         `mapstructure:"port" json:"port"`
	ConfigPath string      `mapstructure:"configPath" json:"configPath"`
	MysqlInfo  MysqlConfig `mapstructure:"mysql" json:"mysql"`
	UserInfo   MysqlConfig `mapstructure:"user" json:"user"`
	GameInfo   MysqlConfig `mapstructure:"game" json:"game"`
	PayInfo    MysqlConfig `mapstructure:"pay" json:"pay"`
	LogInfo    MysqlConfig `mapstructure:"log" json:"log"`
	JwtInfo    JWTConfig   `mapstructure:"jwt" json:"jwt"`
	RedisInfo  RedisConfig `mapstructure:"redis" json:"redis"`
	Mct        Mct         `mapstructure:"mct" json:"mct"`
	InPay      InPay       `mapstructure:"inPay" json:"inPay"`
	Srv        Srv         `mapstructure:"srv" json:"srv"`
}

type Srv struct {
	Hall  string `mapstructure:"hall" json:"hall"`
	Token string `mapstructure:"token" json:"token"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	Db       int    `mapstructure:"db"`
	Expire   int    `mapstructure:"expire" json:"expire"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type Mct struct {
	Addr    string `mapstructure:"addr"`
	Auth    string `mapstructure:"auth"`
	TypeIdx int    `mapstructure:"typeIdx"`
	SvrId   int    `mapstructure:"sid"`
	LogCfg  string `mapstructure:"logCfg"`
}

type InPay struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
