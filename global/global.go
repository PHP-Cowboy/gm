package global

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gm/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	User         *gorm.DB
	Game         *gorm.DB
	Pay          *gorm.DB
	Log          *gorm.DB
	Redis        *RedisCli
	ServerConfig = &config.ServerConfig{}
	Logger       = make(map[string]*zap.SugaredLogger, 0)
	GoCache      *cache.Cache

	RedisPool *redis.Pool
	GameDB    *sqlx.DB
	NDB       *sqlx.DB
	PayDB     *sqlx.DB
	LogDB     *sqlx.DB
)
