
create table users (
    id serial primary key not null,
    username varchar(100) not null,
    email varchar(100) not null,
    passwordHash text not null,
    created timestamp not null,
    modified timestamp not null
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO datacatdbuser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO datacatdbuser;