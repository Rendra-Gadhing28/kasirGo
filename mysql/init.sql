-- Database Initialization for KasirPro
CREATE DATABASE IF NOT EXISTS kasirpro CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE kasirpro;

-- Least privilege user for application runtime
CREATE USER IF NOT EXISTS 'kasirpro_app'@'%' IDENTIFIED BY 'kasirpro_secret_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON kasirpro.* TO 'kasirpro_app'@'%';

FLUSH PRIVILEGES;
