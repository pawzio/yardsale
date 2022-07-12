catalog-run-api:
	${COMPOSE_CMD} run --rm --service-ports catalog-go sh -c "go run cmd/serverd/*.go"
catalog-test:
	${COMPOSE_CMD} run --rm catalog-go sh -c "go test -coverprofile=c.out -failfast -timeout 5m ./..."
catalog-go-vendor:
	${COMPOSE_CMD} run --rm catalog-go sh -c "go mod tidy && go mod vendor"
catalog-go-generate:
	${COMPOSE_CMD} run --rm catalog-go sh -c "go generate ./..."
catalog-setup: catalog-go-vendor catalog-go-generate

generate-codacy-coverage-report-go:
	cat ${SVC_NAME}/c.out > ${SVC_NAME}/filtered-coverage.out

teardown:
	${COMPOSE_CMD} down -v
	${COMPOSE_CMD} rm --force --stop -v

ifndef PROJECT_NAME:
PROJECT_NAME := yardsale
endif

ifndef DOCKER_BIN:
DOCKER_BIN := docker
endif

ifndef DOCKER_COMPOSE_BIN:
DOCKER_COMPOSE_BIN := docker-compose
endif

COMPOSE_CMD := PROJECT_NAME=${PROJECT_NAME} ${DOCKER_COMPOSE_BIN} -f build/docker-compose.yaml -p ${PROJECT_NAME}
