local key = KEYS[1]
--用户输入的code
local expectedCode = ARGV[1]
local cntKey = key..":cnt"
local code = redis.call("get",key)
local cnt = tonumber(redis.call("get",cntKey))

if cnt<=0 then
    --一直输入错误
    return -1
elseif expectedCode==code then
    --输入正确 设置这个key不再用
    redis.call("set",cntKey,-1)
    return 0
else
    --可验证次数 -1
    redis.call("decr",cntKey)
    return -2
end