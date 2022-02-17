APP_NAME:=homebrew-api
APP_VER:=0.0.1
IMAGE_NAME?=homebrew-py
IMAGE_VERSION?=0.0.1

DB_URL:="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}"

.PHONY: build-image
build-image:	## build docker image
	docker build -t ${IMAGE_NAME}:${IMAGE_VERSION} -f py.dockerfile .

.PHONY: db
db:	## create database
	docker rm -f /postgres || true
	docker-compose up -d postgres

.PHONY: db-init
db-init: db
	sleep 5 && docker exec -it \
		postgres \
		psql ${DB_URL} \
		-f /app/sql/db-init.sql

.PHONY: db-conn
db-conn: ## connect to database
	docker exec -it \
		postgres \
		psql ${DB_URL}

.PHONY: db-write
db-write: ## write data to database
	docker-compose run \
		-e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
		-e DATA_PATH="${DATA_PATH}" \
		-e TABLE_NAME="${TABLE_NAME}" \
		homebrew-write

.PHONY: build
build: 
	docker build \
		-t ${APP_NAME}:${APP_VER} . \
		-f go.dockerfile \
		--no-cache

.PHONY: api
api: build
	docker rm -f "/homebrew-api" || true
	docker run -d --name ${APP_NAME} \
		--entrypoint "/app/homebrew-api" \
		--network=powerattackpublishingapi_local \
		-p 8080:8080 \
		-e PATREON_CLIENT_ID="${PATREON_CLIENT_ID}" \
		-e PATREON_CLIENT_SECRET="${PATREON_CLIENT_SECRET}" \
		-e DB_URL=${DB_URL} \
		${APP_NAME}:${APP_VER}
	docker logs ${APP_NAME} -f
