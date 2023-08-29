--你的验证码在 redis 上的 key
local key = KEYS[1]
--验证次数
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
--验证码有效时间
local ttl = tonumber(redis.call("ttl",key))

if ttl == -1 then
    --误操作 key 存在 但是并没有过期时间
    return -2
elseif ttl==-2 or ttl<540 then
    -- =-2代表key不存在  小于540代表已经过了一分多钟 小于十分钟 所以可以再次发送
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    -- 预期实现
    return 0
else
    -- 已经发送了验证码，并且时间小于一分钟
    return -1
end