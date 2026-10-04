-- 一次执行：创建库 + 建表（Go 服务不会在连接时自动 CREATE DATABASE）
-- 用法（按你的账号密码改）：
--   Bash / cmd.exe：
--     mysql -h 127.0.0.1 -u wjhaccount -p < scripts/init-database.sql
--   PowerShell（不支持上面的 < 重定向，用管道）：
--     Get-Content .\scripts\init-database.sql -Raw | mysql -h 127.0.0.1 -u wjhaccount -p
--   或借 cmd 跑输入重定向（整条用单引号，避免 PowerShell 解析 <）：
--     cmd /c 'mysql -h 127.0.0.1 -u wjhaccount -p < scripts\init-database.sql'

CREATE DATABASE IF NOT EXISTS mysqlWithtotp
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE mysqlWithtotp;

CREATE TABLE IF NOT EXISTS user_info
(
    id bigint not null auto_increment,
    username varchar(255) not null default '',
    password_hash varchar(255) not null default '',
    totp_secret varchar(255) not null default '',
    totp_pending_secret varchar(255) not null default '',
    totp_enabled tinyint(1) not null default 0,
    UNIQUE KEY idx_username (username),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
