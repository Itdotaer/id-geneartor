-- insert into id_generator_tab(business, current_id, step) value ('test_11', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_12', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_13', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_14', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_15', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_16', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_17', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_18', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_19', 0, 10000);
-- insert into id_generator_tab(business, current_id, step) value ('test_20', 0, 10000);

function request()
	method = "GET"

	id = math.random(11,20)

	path = "/alloc?business=test_" .. tostring(id)
	return wrk.format(method,path,nil,body)
end

function response(status, headers, body)
end