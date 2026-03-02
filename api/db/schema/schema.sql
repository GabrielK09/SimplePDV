-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.
CREATE TABLE public.customers (
  id SERIAL PRIMARY KEY,
  name character varying NOT NULL,
  cpf_cnpj character varying NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.pay_ment_forms (
  id SERIAL PRIMARY KEY,
  specie character varying NOT NULL,
  pix_key character varying,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.products (
  id SERIAL PRIMARY KEY,
  name character varying NOT NULL,
  price double precision NOT NULL,
  commission double precision NOT NULL,
  qtde integer NOT NULL,
  returned integer DEFAULT 0,
  saled integer DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.sales (
  id SERIAL PRIMARY KEY,
  customer_id integer,
  customer character varying NOT NULL,
  sale_value double precision,
  status character varying DEFAULT 'Pendente'::character varying,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sales_customer_id_foreign FOREIGN KEY (customer_id) REFERENCES public.customers(id)
);
CREATE TABLE public.shopping (
  id SERIAL PRIMARY KEY,
  load integer,
  operation character varying,
  status character varying DEFAULT 'Pendente'::character varying,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.shopping_itens (
  id SERIAL PRIMARY KEY,
  product_id integer,
  name character varying NOT NULL,
  qtde_purchased integer NOT NULL,
  purchased_value double precision NOT NULL,
  shopping_id bigint NOT NULL,
  status character varying DEFAULT 'Associado'::character varying,
  deleted_at timestamp without time zone,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT shopping_itens_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT shopping_itens_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);
CREATE TABLE public.cash_registers (
  id SERIAL PRIMARY KEY,
  description character varying NOT NULL,
  customer_id integer NOT NULL,
  customer character varying NOT NULL,
  specie_id integer NOT NULL,
  specie character varying NOT NULL,
  input_value double precision,
  output_value double precision,
  total_balance double precision NOT NULL,
  sale_id integer,
  shopping_id integer,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  
  CONSTRAINT cash_registers_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT cash_registers_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT cash_registers_customer_id_foreign FOREIGN KEY (customer_id) REFERENCES public.customers(id),
  CONSTRAINT cash_registers_specie_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);
CREATE TABLE public.sale_itens (
  id SERIAL PRIMARY KEY,
  product_id integer NOT NULL,
  name character varying NOT NULL,
  qtde integer NOT NULL,
  sale_value double precision NOT NULL,
  sale_id bigint NOT NULL,
  status character varying DEFAULT 'Pendente'::character varying,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sales_itens_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT sales_itens_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);
CREATE TABLE public.sale_pay_ment (
  id SERIAL PRIMARY KEY,
  sale_id bigint NOT NULL,
  specie_id bigint NOT NULL,
  specie character varying NOT NULL,
  amount_paid double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT pay_ment_forms_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);
CREATE TABLE public.config_pdv (
  id SERIAL PRIMARY KEY,
  confirm_to_pinter BOOLEAN DEFAULT 'false',
  block_sale_negative_stock BOOLEAN DEFAULT 'false',
  reserve_stock BOOLEAN DEFAULT 'true'
);