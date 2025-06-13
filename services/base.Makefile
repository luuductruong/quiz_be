export DATABASE_HOST := quiz_be_postgres
export DATABASE_PORT := 5432
export DATABASE_USER := quiz_be
export DATABASE_PASSWORD := quiz_be

dbs:
	if [ "$${DATABASE_NAME}" != "" ]; then \
	    docker run --rm --network=quiz_be_default \
	        -e PGPASSWORD=$${DATABASE_PASSWORD} \
	        postgres:15 \
	        psql -h $${DATABASE_HOST} -U $${DATABASE_USER} -d postgres \
	        -tc "SELECT 1 FROM pg_database WHERE datname = '$${DATABASE_NAME}'" \
	    | grep -q 1 \
	    || docker run --rm --network=quiz_be_default \
	        -e PGPASSWORD=$${DATABASE_PASSWORD} \
	        postgres:15 \
	        psql -h $${DATABASE_HOST} -U $${DATABASE_USER} -d postgres \
	        -c "CREATE DATABASE $${DATABASE_NAME};"; \
	    if [ "$${ENV}" != "" ]; then \
	        cp $$(pwd)/external/db/seed/$${ENV}/*.sql $$(pwd)/external/db/migration; \
	    else \
	        cp $$(pwd)/external/db/seed/dev/*.sql $$(pwd)/external/db/migration; \
	    fi; \
	    docker run --rm --network=quiz_be_default --name flyway \
	        -v $$(pwd)/external/db/migration:/flyway/sql \
	        -v $$(pwd)/external/db:/flyway/conf \
	        flyway/flyway \
	        -url=jdbc:postgresql://$${DATABASE_HOST}:$${DATABASE_PORT}/$${DATABASE_NAME} \
	        -user=$${DATABASE_USER} -password=$${DATABASE_PASSWORD} \
	        migrate info; \
	fi;
