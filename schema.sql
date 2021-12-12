DROP TABLE IF EXISTS urlandlinks;
CREATE TABLE urlandlinks
(
    url varchar(255) not null unique,
    link varchar(255) not null unique
);

CREATE INDEX links ON urlandlinks (link);