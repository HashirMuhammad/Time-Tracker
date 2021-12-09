-- migrate:up
create table projects
(
    id          bigserial not null
        constraint projects_id_pk
            primary key,
    client_name varchar(255),
    started_by  varchar(255),
    title       varchar(255),
    description varchar(255)
);
--migrate down