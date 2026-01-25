CREATE TABLE sale_itens (
    id bigint NOT NULL,
    product_id INT NOT NULL,
	name character varying NOT NULL,
	qtde INT NOT NULL,
	sale_value FLOAT NOT NULL,
	date_of_movement timestamp without time zone,
	CONSTRAINT sale_item_pkey PRIMARY KEY (id)
)