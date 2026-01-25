CREATE TABLE products (
    id bigint NOT NULL,
	name character varying NOT NULL,
	un VARCHAR(4) DEFAULT 'UN',
	qtde INT NOT NULL,           
	returned INT DEFAULT 0,
	saled INT DEFAULT 0,
	date_of_purchase timestamp without time zone,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)