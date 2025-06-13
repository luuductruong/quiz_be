CMDS = flyway dbs

SERVICES := $(shell find ./services -mindepth 1 -maxdepth 1 -type d | xargs -n 1 basename | grep -vE 'core')

dto:
	docker compose -f docker-compose.yml up quiz_be_proto_builder
deps:
	docker compose -f docker-compose.yml up -d quiz_be_postgres
proto:
	make -C services/core/proto

$(SERVICES):
	$(MAKE) -C services/$@ service SERVICE_NAME=$@

$(CMDS):
	make -C services $@