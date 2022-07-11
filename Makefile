catalog-api-run:
	${COMPOSE} run --rm --service-ports catalog-go sh -c "go run main.go"

catalog-test:
	${COMPOSE} run --rm catalog-go sh -c "go test -coverprofile=c.out -failfast -timeout 5m ./..."

ifndef PROJECT_NAME:
PROJECT_NAME := yardsale
endif

ifndef DOCKER_BIN:
DOCKER_BIN := docker
endif

ifndef DOCKER_COMPOSE_BIN:
DOCKER_COMPOSE_BIN := docker-compose
endif

COMPOSE := PROJECT_NAME=${PROJECT_NAME} ${DOCKER_COMPOSE_BIN} -f build/docker-compose.base.yaml -f build/docker-compose.local.yaml -p ${PROJECT_NAME}
