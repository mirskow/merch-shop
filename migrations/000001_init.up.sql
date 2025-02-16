CREATE TABLE users (
    id serial PRIMARY KEY,
    username varchar(255) not null unique,
    coin int not null DEFAULT 1000,
    password_hash varchar(255) not null
);

CREATE TABLE merch (
    id serial PRIMARY KEY,
    item varchar(255) not null unique,
    cost int not null CHECK (cost > 0)
);

CREATE TABLE purchases (
    id serial PRIMARY KEY,
    user_id int references users(id) on delete cascade not null,
    item_id int references merch(id) on delete cascade not null
);

CREATE TABLE transactions (
    id serial PRIMARY KEY,
    from_user int references users(id) on delete cascade not null,
    to_user int references users(id) on delete cascade not null,
    amount int not null check(amount > 0)
);

INSERT INTO merch (item, cost) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);

INSERT INTO users (username, coin, password_hash) VALUES ('test', 1000, 'test');
