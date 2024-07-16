package rds

import "gm/global"

// 清理redis缓存
func DelRedisCacheByKey(key ...string) (err error) {
	r := global.Redis
	_, err = r.Del(key...)
	return
}
