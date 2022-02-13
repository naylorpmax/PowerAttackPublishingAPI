APP_NAME:=homebrew-api
APP_VER:=0.0.1
IMAGE_NAME?=homebrew-py
IMAGE_VERSION?=0.0.1

.PHONY: build-image
build-image:	## build docker image
	docker build -t ${IMAGE_NAME}:${IMAGE_VERSION} .

.PHONY: db
db:	## create database
	docker rm -f /homebrew-db || true
	docker-compose up -d homebrew-db

.PHONY: db-init
db-init: db
	sleep 5 && docker exec -it \
		homebrew-db \
		psql "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}" \
		-f /app/sql/db-init.sql

.PHONY: db-conn
db-conn: ## connect to database
	docker exec -it \
		homebrew-db \
		psql "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}"

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
		--network=local \
		-p 8080:8080 \
		-e PATREON_CLIENT_ID="${PATREON_CLIENT_ID}" \
		-e PATREON_CLIENT_SECRET="${PATREON_CLIENT_SECRET}" \
		${APP_NAME}:${APP_VER}	 		

# 		--service-ports \
# 		--name homebrew-api \
