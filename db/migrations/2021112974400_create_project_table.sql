-- migrate:up
create table tasks
(
    id         bigserial                                          not null
        constraint tasks_id_pk
            primary key,
    user_id bigint
        constraint tasks_users_fk
            references users(id),
    description varchar(255),
    started_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    ended_at timestamp with time zone default CURRENT_TIMESTAMP not null
);
--migrate down