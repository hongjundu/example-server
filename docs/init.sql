CREATE USER "test"@"localhost" IDENTIFIED BY "123456";
CREATE USER "test"@"%" IDENTIFIED BY "123456";
create database test;
grant all privileges on test.* to "test"@"%";
flush privileges;