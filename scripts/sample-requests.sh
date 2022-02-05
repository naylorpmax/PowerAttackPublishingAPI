# !/bin/sh

curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"whatever"}' | python -m json.tool

curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/monster/lookup \
	-d '{"Name":"whatever"}' | python -m json.tool

curl -X POST \
	-H "content-type: something-else" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"whatever"}' | python -m json.tool

curl -X POST \
	-H "content-type: something-else" \
	http://localhost:8080/monster/lookup \
	-d '{"Name":"whatever"}' | python -m json.tool
