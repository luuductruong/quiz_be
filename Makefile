dto:
	docker compose -f docker-compose.yml up quiz_be_proto_builder
dbs:
	docker compose -f docker-compose.yml up -d quiz_be_postgres
proto:
	make -C services/core/proto
