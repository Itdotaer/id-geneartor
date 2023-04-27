function request()
	method = "GET"

	id = math.random(1,10)

	path = "/alloc?business=test_" .. tostring(id)
	return wrk.format(method,path,nil,body)
end

function response(status, headers, body)
--    print(body)
end