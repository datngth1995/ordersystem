-- Database: dordersystem

-- DROP DATABASE dordersystem;

CREATE DATABASE dordersystem
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    TEMPLATE template0;

COMMENT ON DATABASE dordersystem
    IS 'Database to store all user info, product info and order info';
