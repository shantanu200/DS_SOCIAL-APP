USER_BINARY=userApp
TWEET_BINARY=tweetApp
NOTIFICATION_BINARY=notificationApp
USER_RELATION_BINARY=userRelationApp
TIMELINE_BINARY = timeLineApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_user build_tweet build_notification build_user_relation build_timeline
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"  
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

build_user:
	@echo "Building user binary..."
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"


build_tweet:
	@echo "Building tweet binary..."
	cd ../tweet-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TWEET_BINARY} ./cmd/api
	@echo "Done!"


build_notification:
	@echo "Building notification binary..."
	cd ../notification-service && env GOOS=linux CGO_ENABLED=0 go build -o ${NOTIFICATION_BINARY} ./cmd/socket
	@echo "Done!"

build_user_relation:
	@echo "Building user-relation binary..."
	cd ../user-relation-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_RELATION_BINARY} ./cmd/api
	@echo "Done!"

build_timeline:
	@echo "Building timeline binary..."
	cd ../timeline-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TIMELINE_BINARY} ./cmd/api
	@echo "Done!"


