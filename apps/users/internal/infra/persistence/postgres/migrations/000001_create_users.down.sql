-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS identity_verifications;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;