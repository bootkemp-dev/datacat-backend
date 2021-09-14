
create table users (
    id serial primary key not null,
    username varchar(50) unique not null,
    email varchar(100) unique not null,
    passwordHash text not null,
    created timestamp not null,
    modified timestamp not null,
    passwordResetToken varchar(30),
    passwordResetTokenExpDate timestamp
);

create table jobs (
    id serial primary key not null,
    jobName varchar(50) not null,
    jobUrl text not null,
    frequency bigint not null,
    userid integer not null references users(id),
    active boolean not null,
    created timestamp not null,
    modified timestamp not null
);

create table jobLog(
    id serial primary key not null,
    userID integer not null references users(id),
    jobID integer not null references jobs(id),
    jobStatus varchar(30),
    logMessage text not null,
    timeChecked timestamp not null
);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO datacatdbuser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO datacatdbuser;
