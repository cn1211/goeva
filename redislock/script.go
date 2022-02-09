
package redislock

// 释放key
var releaseScript = NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

// 增加过期时间
var addTTLScript = NewScript(1, `
	local ttl = redis.call("TTL", KEYS[1])
	if ttl > 0 then
		return redis.call("SETEX", KEYS[1], ARGV[1]+ttl , ARGV[2])
	else
		return -2
	end
`)
