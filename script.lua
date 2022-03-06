-- 这个命令等价于 set key1 argv1 EX argv2
-- 比如下面这个栗子,设置age是18过期时间是60
-- set age 18 EX 60
redis.call('SET',KEYS[1],ARGV[1])
redis.call('EXPIRE',KEYS[1],ARGV[3])

redis.call('SET',KEYS[2],ARGV[2])
redis.call('EXPIRE',KEYS[2],ARGV[3])
return 1