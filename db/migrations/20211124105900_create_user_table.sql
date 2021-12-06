-- migrate:up
create table users
(
    id         bigserial                                          not null
        constraint users_pk
            primary key,
    first_name varchar(255),
    last_name varchar(255),
    email varchar(255)                                         not null,
    password varchar(255)                                         not null,
    image_url varchar(255),
    created_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP not null
);
--migrate down