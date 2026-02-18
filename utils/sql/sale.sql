CREATE TABLE sales (
  id SERIAL PRIMARY KEY,
  customer character varying NOT NULL,
  sale_value FLOAT,
  status character varying DEFAULT 'Pendente',  
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);