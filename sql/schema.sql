CREATE SCHEMA gaea;

CREATE TABLE gaea.user (
    user_name text PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    role text NOT NULL,
    password text,
    dip_id text,
    passport text,
    section text,
    updated_at timestamp,
    update_token text
);

CREATE TABLE gaea.sale (
    sale_id uuid PRIMARY KEY,
    sale_type text,
    open_date timestamp,
    close_date timestamp,
    status text
);

CREATE TABLE gaea.inventory (
    inventory_id uuid PRIMARY KEY,
    sale_id uuid REFERENCES gaea.sale (sale_id),
    updated_at timestamp,
    supplier_id text NOT NULL,
    name text NOT NULL,
    description text,
    abv text,
    size text,
    year text,
    nonmem_price integer,
    mem_price integer NOT NULL,
    types text[],
    origin text[],
    changelog text[]
);


CREATE TABLE gaea.order (
    order_id uuid PRIMARY KEY,
    sale_id uuid REFERENCES gaea.sale (sale_id),
    status text,
    status_date timestamp,
    user_name text REFERENCES gaea.user (user_name),
    sale_type text
);

CREATE TABLE gaea.orderitem (
    orderitem_id uuid PRIMARY KEY,
    order_id uuid REFERENCES gaea.order (order_id),
    inventory_id uuid REFERENCES gaea.inventory (inventory_id),
    qty integer,
    updated_at timestamp,
    user_name text REFERENCES gaea.user (user_name)
);

INSERT INTO gaea.user (
    user_name, 
    first_name, 
    last_name,
    email,
    role,
    password) VALUES (
    'burkebt',
    'Bryan',
    'Burke',
    'btburke@fastmail.com',
    'superadmin',
    '16384$8$1$84c73e785d4d9a45df5923cf1663af04$59a8f646c5e13714cf0fe2ee832af9aa03ae32779c9d9157fb65e9ab98cc1bfd'
);