CREATE TABLE IF NOT EXISTS egw.product (
	id uuid NOT NULL,
	"name" varchar NOT NULL,
	short_description varchar NOT NULL,
	description varchar NOT NULL,
	price int4 NOT NULL,
	created_at time NULL,
	updated_at time NULL,
	CONSTRAINT product_pk PRIMARY KEY (id)
);