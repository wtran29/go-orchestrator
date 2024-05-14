tasks:
	curl localhost:5556/tasks | python -m json.tool

stats:
	curl localhost:5555/stats | python -m json.tool

docker_ps:
	docker ps --format "table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Names}}"

post_task:
	curl -v --request POST \
	--header 'Content-Type: application/json' \
	--data @task.json \
	localhost:5556/tasks | python -m json.tool

get_wtask:
	curl -v localhost:5555/tasks | python -m json.tool

get_mtask:
	curl -v localhost:5556/tasks | python -m json.tool

delete_task:
	curl -v --request DELETE \
	'localhost:5556/tasks/21b23589-5d2d-4731-b5c9-a97e9832d021'

echo:
	docker run -p 7777:7777 --name echo timboring/echo-server:latest

test_echo:
	curl -X POST http://localhost:7777/ -d '{"Msg": "hello world"}'

health:
	curl -v http://localhost:7777/health

healthfail:
	 curl -v http://localhost:7777/healthfail

task2:
	curl -v -X POST localhost:5556/tasks -d @task2.json
