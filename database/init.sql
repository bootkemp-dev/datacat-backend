
create table users (
    id serial primary key unique not null,
    username varchar(50) unique not null,
    email varchar(100) unique not null,
    passwordHash text not null,
    created timestamp not null,
    modified timestamp not null
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO datacatdbuser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO datacatdbuser;