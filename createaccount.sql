-- select relname,n_dead_tup,n_live_tup from pg_stat_all_tables where n_dead_tup <> 0 and n_live_tup <>0 order by 2 desc;

CREATE DATABASE tupletest;

CREATE TABLE accounts ( id SERIAL PRIMARY KEY,  username VARCHAR (50) UNIQUE NOT NULL,  created_at TIMESTAMP);
INSERT INTO accounts (username) VALUES ('teste');
INSERT INTO accounts (username) VALUES ('teste2');
INSERT INTO accounts (username) VALUES ('teste3');
INSERT INTO accounts (username) VALUES ('teste4');

DELETE FROM accounts where id = 0;