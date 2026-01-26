CREATE TABLE cash_registers (
  id SERIAL PRIMARY KEY,
  description character varying NOT NULL,
  customer character varying NOT NULL,
  specie character varying NOT NULL,
  input_value float,
  output_value float,
  total_balance float NOT NULL,
  sale_id integer,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  CONSTRAINT cash_registers_sales_id_foreign FOREIGN key (sale_id) REFERENCES public.sales(id)
)