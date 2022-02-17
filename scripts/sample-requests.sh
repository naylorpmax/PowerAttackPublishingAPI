# !/bin/sh

# multiple spells returned
curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"Holy"}' | python -m json.tool

# single spell returned - exact match
curl -X POST \
	-H "content-type: application/json" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"Rings a Bell"}' | python -m json.tool

# bad request: incorrect content-type
curl -X POST \
	-H "content-type: something-else" \
	http://localhost:8080/spell/lookup \
	-d '{"Name":"whatever"}' | python -m json.tool
