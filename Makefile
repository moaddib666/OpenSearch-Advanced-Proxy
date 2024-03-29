# Define the path to the docker-compose file
COMPOSE_FILE := infrastrcuture/docker/docker-compose.yaml

# Default action when just 'make' is run
.PHONY: default
default: up

# Start the services in the background
.PHONY: up
up:
	@docker-compose -f $(COMPOSE_FILE) up -d
	@echo "Access the Dashboard at http://localhost:5601"

# Stop the services
.PHONY: down
down:
	docker-compose -f $(COMPOSE_FILE) down

# Stop and remove containers, networks, images, and volumes
.PHONY: clean
clean:
	docker-compose -f $(COMPOSE_FILE) down --rmi all --volumes

# View logs of the services
.PHONY: logs
logs:
	docker-compose -f $(COMPOSE_FILE) logs

# Build or rebuild services
.PHONY: build
build:
	docker-compose -f $(COMPOSE_FILE) build

# Configure the services
.PHONY: configure
configure:
	@if [ ! -d "./.local/config" ]; then \
		cp -prf examples/config ./.local/config; \
	fi
	@if [ ! -d "./.local/shard_config" ]; then \
		cp -prf examples/shard_config ./.local/shard_config; \
	fi
	@if [ ! -f "./.local/test.log" ]; then \
		cp -prf examples/test.log ./.local/test.log; \
	fi


