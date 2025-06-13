ARG SERVICE_NAME=quiz
FROM flyway/flyway:10.11.0

ARG SERVICE_NAME
ARG BRANCH
COPY services/${SERVICE_NAME}/external/db/flyway.conf /flyway/conf/flyway.conf
COPY services/${SERVICE_NAME}/external/db/migration/*.sql /flyway/sql/
COPY services/${SERVICE_NAME}/external/db/seed/${BRANCH}/*.sql /flyway/sql/
