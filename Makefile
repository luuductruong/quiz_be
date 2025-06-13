ifneq (,$(wildcard .env))
	include .env
	export
endif

CMDS = dbs

dto:
	docker compose -f docker-compose.yml up quiz_be_proto_builder
deps:
	docker compose -f docker-compose.yml up -d quiz_be_postgres
proto:
	make -C services/core/proto

$(CMDS):
	make -C services $@