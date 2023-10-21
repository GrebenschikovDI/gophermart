CREATE TABLE if not exists users
(
    id            serial primary key,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null,
    created_at    timestamp    not null default now()
);

CREATE TABLE if not exists orders
(
    id          serial primary key,
    user_id     int references users (id) on delete cascade not null,
    status      varchar(255),
    uploaded_at timestamp                                   not null default now()
);

CREATE TABLE if not exists balance
(
    id           serial primary key,
    user_id      int references users (id) on delete cascade not null,
    order_id     int references orders (id)                  not null,
    amount       double precision                            not null default 0,
    processed_at timestamp                                   not null default now()
);

CREATE TABLE if not exists withdrawals
(
    id           serial primary key,
    user_id      int references users (id) on delete cascade not null,
    order_id     int references orders (id)                  not null,
    amount       double precision                            not null default 0,
    total        double precision                            not null default 0,
    processed_at timestamp                                   not null default now()
)