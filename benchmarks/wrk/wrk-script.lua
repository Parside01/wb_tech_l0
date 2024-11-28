local uuids = {}
local file = io.open("../uuids.txt", "r")
for line in file:lines() do table.insert(uuids, line)
end
file:close()
request = function()
    local id = uuids[math.random(#uids)]
    local path = "/api/v1/order/" .. id
    return wrk.format("GET", path)
end