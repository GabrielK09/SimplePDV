CREATE TABLE sale_itens (
	id SERIAL PRIMARY KEY,
	product_id INT NOT NULL,
	name character varying NOT NULL,
	qtde INT NOT NULL,
	sale_value FLOAT NOT NULL,
	sale_id bigint NOT NULL,
	status character varying DEFAULT 'Pendente',
	created_at timestamp without time zone,
  	updated_at timestamp without time zone,
	CONSTRAINT sales_itens_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
	CONSTRAINT sales_itens_products_id_foreign FOREIGN key (product_id) REFERENCES public.products(id)
);

DROP TABLE sale_itens;