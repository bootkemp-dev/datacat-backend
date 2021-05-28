
create table users (
    id serial primary key not null,
    username varchar(50) unique not null,
    email varchar(100) unique not null,
    passwordHash text not null,
    created timestamp not null,
    modified timestamp not null
);

create table jobs (
    id serial primary key not null,
    jobName varchar(50) not null,
    jobUrl text not null,
    frequency integer not null,
    userid integer not null references users(id)
);

create table jobLog(
    id serial primary key not null,
    jobID integer not null references jobs(id),
    down boolean,
    downFor timestamp,
    timeChecked timestamp not null
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO datacatdbuser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO datacatdbuser;