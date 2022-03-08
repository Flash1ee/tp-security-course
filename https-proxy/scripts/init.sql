create table if not exists requests(
    id bigserial primary key,
    method text,
    path text,
    get_params jsonb,
    headers jsonb,
    cookies jsonb,
    post_params jsonb
);
create table if not exists responses(
    id bigserial primary key,
    request_id bigint references requests(id),
    code int,
    message text,
    headers jsonb,
    body text
);