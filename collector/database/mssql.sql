-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS collector WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

-- Create a table
CREATE TABLE IF NOT EXISTS collector.logs (
id int PRIMARY KEY,
level text,
msg text,
timestamp BIGINT);

-- Insert some data
INSERT INTO collector.logs
(id,level, msg, timestamp )
VALUES (1, 'Warning', 'MSG cassandra',1677207144000);
INSERT INTO collector.logs
(id,level, last_update_timestamp)
VALUES ('1234', 'Imformation', toTimeStamp(now()));

--  MSSQL
/opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P "Khong123"

CREATE DATABASE TestDB;

SELECT Name from sys.databases;

go

USE TestDB;

CREATE TABLE logs (id INT, level NVARCHAR(50), msg TEXT, time BIGINT);
INSERT INTO logs VALUES (1, "Error", "MSSQL Log 1",1678106049000 ) ;

--
set https_proxy=http://192.168.5.8:3128
set http_proxy=http://192.168.5.8:3128
go env -w GOPROXY=https://proxy.golang.org,direct

