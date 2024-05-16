tasks:
	curl localhost:5555/tasks | python -m json.tool

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

get_nodes:
	curl localhost:5555/nodes | python -m json.tool

delete_task:
	curl -v --request DELETE \
	'localhost:5555/tasks/bb1d59ef-9fc1-4e4b-a44d-db571eeed203'

echo:
	docker run -p 7777:7777 --name echo timboring/echo-server:latest

test_echo:
	curl -X POST http://localhost:7777/ -d '{"Msg": "hello world"}'

health:
	curl -v http://localhost:7777/health

healthfail:
	 curl -v http://localhost:7777/healthfail

task1:
	curl -v -X POST localhost:5555/tasks -d @task1.json

delete1:
	curl -X DELETE localhost:5555/tasks/bb1d59ef-9fc1-4e4b-a44d-db571eeed203

w1:
	go run main.go worker -p 5556
w2:
	go run main.go worker -p 5557
w3:
	go run main.go worker -p 5558

m:
	go run main.go manager -w 'localhost:5556,localhost:5557,localhost:5558'

run_task:
	go run main.go run --filename task1.json

stop_task:
	go run main.go stop bb1d59ef-9fc1-4e4b-a44d-db571eeed203

status:
	go run main.go status

nodes:
	go run main.go node