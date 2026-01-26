CREATE TABLE sales (
  id SERIAL PRIMARY KEY,
  customer character varying NOT NULL,
  specie character varying NOT NULL,
  sale_value FLOAT,
  status character DEFAULT 'Pendente',  
  created_at timestamp without time zone,
  updated_at timestamp without time zone
)