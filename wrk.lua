-- insert into id_generator_tab(business, current_id, step) value ('test_1', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_2', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_3', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_4', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_5', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_6', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_7', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_8', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_9', 0, 1000);
-- insert into id_generator_tab(business, current_id, step) value ('test_10', 0, 1000);

function request()
	method = "GET"

	id = math.random(1,10)

	path = "/alloc?business=test_" .. tostring(id)
	return wrk.format(method,path,nil,body)
end

function response(status, headers, body)
end