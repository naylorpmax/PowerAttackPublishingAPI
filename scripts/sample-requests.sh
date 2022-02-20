# !/bin/sh

# multiple spells returned by name
curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"Holy"}' | python3 -m json.tool

# multiple spells returned by name and level
curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"Holy","Level":"1"}' | python3 -m json.tool

# single spell returned - exact match
curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"Rings a Bell"}' | python3 -m json.tool

# bad request: incorrect content-type
curl -X POST \
	-H "content-type: something-else" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"whatever"}' | python3 -m json.tool

curl -X GET \
	-H "content-type: application/json" \
	-H "accept: application/json" \
	http://localhost:8080/login

