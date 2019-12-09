CREATE USER "test"@"localhost" IDENTIFIED BY "123456";
CREATE USER "test"@"%" IDENTIFIED BY "123456";
create database test;
grant all privileges on test.* to "test"@"%";
flush privileges;


CREATE TABLE casbin_rule (     p_type VARCHAR(100),     v0 VARCHAR(100),     v1 VARCHAR(100),     v2 VARCHAR(100) );
INSERT INTO casbin_rule VALUES('p', 'user', 'data', 'read');
INSERT INTO casbin_rule VALUES('p', 'admin', 'data', 'read');
INSERT INTO casbin_rule VALUES('p', 'admin', 'data', 'write');
INSERT INTO casbin_rule(p_type, v0, v1)  VALUES('g', 'duhj', 'user');
INSERT INTO casbin_rule(p_type, v0, v1)  VALUES('g', 'admin', 'admin');