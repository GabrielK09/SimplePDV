CREATE SCHEMA public;

CREATE TABLE public.users (
  id SERIAL PRIMARY KEY,
  name character varying NOT NULL,
  cpf character varying NOT NULL,
  login character varying NOT NULL unique,
  password character varying NOT NULL,
  is_admin boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.customers (
  id SERIAL PRIMARY KEY,
  name character varying NOT NULL,
  cpf_cnpj character varying NOT NULL,
  deleted_at timestamp without time zone,
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
  reserved_qtde integer NULL DEFAULT 0,
  future_qtde integer NULL DEFAULT 0,
  use_grid BOOLEAN DEFAULT 'true',
  deleted_at timestamp without time zone,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TYPE sizes AS ENUM ('PP', 'P', 'M', 'G', 'GG', 'XG', 'XGG', 'EG', 'EGG', 'O')

CREATE TABLE public.product_grids (
  id SERIAL PRIMARY KEY,
  product_id integer NOT NULL,
  size sizes NOT NULL,
  grid_qtde integer NOT NULL,
  deleted_at timestamp without time zone,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT products_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);
ALTER TABLE public.product_grids
ADD CONSTRAINT unique_product_grid_size
UNIQUE (size, product_id);

-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 
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
CREATE TABLE public.sale_itens_grid (
  id SERIAL PRIMARY KEY,
  product_id integer NULL,
  sale_id integer NULL,
  product_grid_id integer NULL,
  size_saled sizes NULL,
  grid_qtde integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sale_itens_grid_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT sale_itens_grid_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id),
  CONSTRAINT sale_itens_grid_product_grids_id_foreign FOREIGN KEY (product_grid_id) REFERENCES public.product_grids(id)
);
ALTER TABLE sale_itens_grid
ADD CONSTRAINT unique_sale_grid
UNIQUE (sale_id, size_saled, product_grid_id);
-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 

CREATE TABLE public.sale_pay_ment (
  id SERIAL PRIMARY KEY,
  sale_id bigint NOT NULL,
  specie_id bigint NOT NULL,
  specie character varying NOT NULL,
  amount_paid double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sale_pay_ment_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT sale_pay_ment_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);
-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 
CREATE TABLE public.shopping (
  id SERIAL PRIMARY KEY,
  load integer unique,
  operation character varying,
  status character varying DEFAULT 'Pendente'::character varying,
  total_shopping double precision,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE public.shopping_itens (
  id SERIAL PRIMARY KEY,
  shopping_id integer,
  product_id integer,
  name character varying NOT NULL,
  qtde_purchased integer NOT NULL,
  purchased_value double precision NOT NULL,
  status character varying DEFAULT 'Pendente'::character varying,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT shopping_itens_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT shopping_itens_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);
CREATE TABLE public.shopping_itens_grid (
  id SERIAL PRIMARY KEY,
  product_id integer NULL,
  shopping_id integer NULL,
  product_grid_id integer NULL,
  size_saled sizes NULL,
  grid_qtde integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT shopping_itens_grid_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping_id(id),
  CONSTRAINT shopping_itens_grid_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id),
  CONSTRAINT shopping_itens_grid_product_grids_id_foreign FOREIGN KEY (product_grid_id) REFERENCES public.product_grids(id)
);
ALTER TABLE shopping_itens_grid
ADD CONSTRAINT unique_shopping_grid
UNIQUE (shopping_id, size_saled, product_grid_id);

ALTER TABLE shopping_itens_grid
ADD CONSTRAINT shopping_itens_grid_shopping_id_foreign
REFERENCES public.shopping_id(id)
-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 
CREATE TABLE public.cash_registers (
  id SERIAL PRIMARY KEY,
  sale_id integer,
  shopping_id integer,
  description character varying NOT NULL,
  customer_id integer NOT NULL,
  customer character varying NOT NULL,
  specie_id integer NOT NULL,
  specie character varying NOT NULL,
  input_value double precision,
  output_value double precision,
  total_balance double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  
  CONSTRAINT cash_registers_sale_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT cash_registers_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT cash_registers_customer_id_foreign FOREIGN KEY (customer_id) REFERENCES public.customers(id),
  CONSTRAINT cash_registers_specie_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);
CREATE TABLE public.shopping_pay_ment (
  id SERIAL PRIMARY KEY,
  shopping_id bigint NOT NULL,
  specie_id bigint NOT NULL,
  specie character varying NOT NULL,
  amount_paid double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT shopping_pay_ment_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id),
  CONSTRAINT shopping_pay_ment_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id)
);