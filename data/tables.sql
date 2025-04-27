create table users (
                       id integer primary key autoincrement,
                       email varchar(100) not null unique,
                       password varchar(256) not null,
                       display_name varchar(100) not null unique,
                       profile_img_x64 text
);