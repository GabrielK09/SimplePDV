CREATE TABLE cash_registers (
  id SERIAL PRIMARY KEY,
  description character varying NOT NULL,
  customer character varying NOT NULL,
  specie_id integer NOT NULL,
  specie character varying NOT NULL,
  input_value float,
  output_value float,
  total_balance float NOT NULL,
  sale_id integer null,
  shopping_id integer null,
  created_at timestamp without time zone,
  updated_at timestamp without time zone,
  CONSTRAINT cash_registers_sales_id_foreign FOREIGN key (sale_id) REFERENCES public.sales(id),
  CONSTRAINT cash_registers_shopping_id_foreign FOREIGN key (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT cash_registers_specie_id_foreign FOREIGN key (specie_id) REFERENCES public.pay_ment_forms(id)
)