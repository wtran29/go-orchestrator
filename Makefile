tasks:
	curl localhost:5555/tasks | python -m json.tool

stats:
	curl localhost:5555/stats | python -m json.tool

docker_ps:
	docker ps --format "table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Names}}"