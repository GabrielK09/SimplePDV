CREATE TABLE sales (
  id SERIAL PRIMARY KEY,
  customer_id integer null,  
  customer character varying NOT NULL,
  sale_value FLOAT,
  status character varying DEFAULT 'Pendente',  
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  CONSTRAINT sales_customer_id_foreign FOREIGN key (customer_id) REFERENCES public.customers(id)
);