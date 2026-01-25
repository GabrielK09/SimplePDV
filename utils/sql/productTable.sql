CREATE TABLE products (
    id bigint NOT NULL,
	name character varying NOT NULL,
	un VARCHAR(4),
	qtde INT NOT NULL,           
	returned INT,
	saled INT,      
	date_of_purchase timestamp without time zone,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)