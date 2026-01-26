CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	name character varying NOT NULL,
	price FLOAT NOT NULL,
	qtde INT NOT NULL,           
	returned INT DEFAULT 0,
	saled INT DEFAULT 0,
	created_at timestamp without time zone,
  	updated_at timestamp without time zone
)