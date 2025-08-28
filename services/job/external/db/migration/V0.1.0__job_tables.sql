create table if not exists job
(
    id              text                                             not null
    primary key,
    name            text                                             not null,
    data            jsonb,
    status          text                     default 'PENDING'::text not null,
    retry_count     integer                  default 0               not null,
    last_run_status text,
    last_error      text,
    last_run_at     timestamp with time zone,
    trace_id        text,
    exchange        text,
    routing_key     text,
    created_at      timestamp with time zone default now()           not null,
    updated_at      timestamp with time zone default now()           not null
    );
