package IRedis

type IRedis interface {
	GetRedisSync(key string, field string) string
	SetRedisSync(key string, field string, value interface{})
	IncRedisSyn(key string, field string, number int64) int64
	GetName() string
}

type redisDB struct {
	i map[string]IRedis
}

var db = redisDB{map[string]IRedis{}}

func SetRedisAgent(d IRedis) {
	db.i[d.GetName()] = d
}

func GetRedisSync(serviceName string, key string, field string) string {
	return db.i[serviceName].GetRedisSync(key, field)
}

func SetRedisSync(serviceName string, key string, field string, value interface{}) {
	db.i[serviceName].SetRedisSync(key, field, value)
}

func IncRedisSyn(serviceName string, key string, field string, number int64) int64 {
	return db.i[serviceName].IncRedisSyn(key, field, number)
}
