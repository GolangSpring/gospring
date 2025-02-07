COMPOSE_FILES := $(wildcard ./dev/docker-compose.*.yaml)
NETWORK_NAME := dev

# Load environment variables from dev.env before running up or down
env:
	@echo "Loading environment variables from dev.env..."
	@bash -c 'source ./dev/dev.env && export $(grep -v "^#" ./dev/dev.env | xargs)'


up: env network
	@echo "Loading environment variables from dev.env..."
		for file in $(COMPOSE_FILES); do \
			echo "Starting $$file..."; \
			docker-compose -f $$file up -d; \
		done

down:
	@echo "Stopping all services..."
	for file in $(COMPOSE_FILES); do \
		echo "Stopping $$file..."; \
		docker-compose -f $$file down -v; \
	done
	@docker network rm $(NETWORK_NAME) || true

network:
	@docker network create $(NETWORK_NAME) || true