--用户id
local userId    = tostring(KEYS[1])
--订单集合
local orderSet=tostring(KEYS[2])
-- 商品库存key
local goodsTotal=tostring(ARGV[1])
--订单队列
local orderList=tostring(ARGV[2])

-- 是否已经抢购到了,如果是返回
local hasBuy = tonumber(redis.call("sIsMember", orderSet, userId))
if hasBuy ~= 0 then
    return 0
end

-- 库存的数量
local total=tonumber(redis.call("GET", goodsTotal))
--return total
-- 是否已经没有库存了,如果是返回
if total <= 0 then
    return 0
end

-- 可以下单
local flag

-- 增加至订单队列
flag = redis.call("LPUSH", orderList, userId)

-- 增加至用户集合
flag = redis.call("SADD", orderSet, userId)

-- 库存数减1
flag = redis.call("DECR", goodsTotal)
-- 返回当时缓存的数量
return total

--[[


--  多行注释
]]