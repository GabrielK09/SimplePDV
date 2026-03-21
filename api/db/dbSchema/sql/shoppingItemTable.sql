CREATE TABLE shopping_itens (
    id SERIAL PRIMARY KEY,
    product_id INT NULL,
    name character varying NOT NULL,
    qtde_purchased INT NOT NULL,
    purchased_value FLOAT NOT NULL,
    shopping_id bigint NOT NULL,
    status character varying DEFAULT 'Associado',
    deleted_at timestamp without time zone,
    created_at timestamp without time zone,
  	updated_at timestamp without time zone,
    CONSTRAINT shopping_itens_shopping_id_foreign FOREIGN KEY (sale_id) REFERENCES public.shopping(id),
	CONSTRAINT shopping_itens_products_id_foreign FOREIGN key (product_id) REFERENCES public.products(id)
)