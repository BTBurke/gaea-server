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
    updated_at timestamp NOT NULL,
    update_token text,
    last_login timestamp NOT NULL,
    member_exp timestamp NOT NULL,
    member_type text,
    stripe_token text
);

CREATE TABLE gaea.sale (
    sale_id serial PRIMARY KEY,
    sale_type text NOT NULL,
    open_date timestamp NOT NULL,
    close_date timestamp NOT NULL,
    status text NOT NULL,
    title text NOT NULL,
    salescopy text NOT NULL,
    require_final boolean NOT NULL
);

CREATE TABLE gaea.inventory (
    inventory_id serial PRIMARY KEY,
    sale_id serial REFERENCES gaea.sale (sale_id),
    updated_at timestamp NOT NULL,
    supplier_id text NOT NULL,
    name text NOT NULL,
    description text,
    abv text,
    size text,
    year text,
    nonmem_price numeric(7,2) NOT NULL,
    mem_price numeric(7,2) NOT NULL,
    types text NOT NULL,
    origin text,
    in_stock boolean,
    changelog text,
    use_case_pricing boolean,
    case_size integer,
    split_case_penalty_per_item_pct integer,
    currency varchar(3)
);


CREATE TABLE gaea.order (
    order_id serial PRIMARY KEY,
    sale_id serial REFERENCES gaea.sale (sale_id),
    status text,
    status_date timestamp,
    user_name text REFERENCES gaea.user (user_name),
    sale_type text
);

CREATE TABLE gaea.orderitem (
    orderitem_id serial PRIMARY KEY,
    order_id serial REFERENCES gaea.order (order_id),
    inventory_id serial REFERENCES gaea.inventory (inventory_id),
    qty integer NOT NULL,
    updated_at timestamp,
    user_name text REFERENCES gaea.user (user_name)
);

CREATE TABLE gaea.transaction (
    transaction_id serial PRIMARY KEY,
    sale_id serial REFERENCES gaea.sale (sale_id),
    order_id serial REFERENCES gaea.order (order_id),
    user_name text REFERENCES gaea.user (user_name),
    from0 text NOT NULL,
    to1 text NOT NULL,
    description text NOT NULL,
    amount numeric(7,2) NOT NULL,
    type text NOT NULL,
    status text NOT NULL,
    track text,
    notes text,
    pay_date timestamp,
    updated_at timestamp,
    authorized_by1 text REFERENCES gaea.user (user_name),
    authorized_by2 text REFERENCES gaea.user (user_name)
);

CREATE TABLE gaea.announcement (
  announcement_id serial PRIMARY KEY,
  title text NOT NULL,
  markdown text NOT NULL,
  show_at timestamp NOT NULL,
  show_until timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

INSERT INTO gaea.user (
    user_name,
    first_name,
    last_name,
    email,
    role,
    password,
    dip_id,
    passport,
    section,
    updated_at,
    update_token,
    last_login,
    member_exp,
    member_type,
    stripe_token
    ) VALUES (
    'burkebt',
    'Bryan',
    'Burke',
    'btburke@fastmail.com',
    'superadmin',
    '16384$8$1$84c73e785d4d9a45df5923cf1663af04$59a8f646c5e13714cf0fe2ee832af9aa03ae32779c9d9157fb65e9ab98cc1bfd',
    '',
    '910267372',
    'NIV',
    '2015-06-06T00:00:00Z',
    '',
    '1900-01-01T00:00:00Z',
    '2016-04-01T00:00:00Z',
    'family',
    ''
);
