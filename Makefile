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

get_task:
	curl -v localhost:5556/tasks | python -m json.tool

delete_task:
	curl -v --request DELETE \
	'localhost:5556/tasks/21b23589-5d2d-4731-b5c9-a97e9832d021'