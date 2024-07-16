package initialize

import (
	"github.com/patrickmn/go-cache"
	"gm/global"
	"time"
)

func init() {
	global.GoCache = cache.New(5*time.Second, 30*time.Second) // 默认过期时间5s；清理间隔30s，即每30s钟会自动清理过期的键值对
}
