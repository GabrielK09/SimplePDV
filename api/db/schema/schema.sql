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

CREATE INDEX idx_customers_name ON customers(name);

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
  use_grid BOOLEAN DEFAULT false,
  deleted_at timestamp without time zone,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_productS_active_name ON products(name) WHERE deleted_at IS NULL;

CREATE TYPE sizes AS ENUM ('PP', 'P', 'M', 'G', 'GG', 'XG', 'XGG', 'EG', 'EGG', 'O');

CREATE TABLE public.product_grids (
  id SERIAL PRIMARY KEY,
  product_id integer NOT NULL,
  size sizes NOT NULL,
  grid_qtde integer NOT NULL,
  deleted_at timestamp without time zone,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_products_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);
ALTER TABLE public.product_grids
ADD CONSTRAINT unique_product_grid_size
UNIQUE (size, product_id);

CREATE INDEX idx_product_grid_product_id ON product_grids(product_id);

-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 
CREATE TYPE sale_shopping_status AS ENUM('Pendente', 'Concluída', 'Cancelada');
CREATE TYPE sale_shopping_itens_status AS ENUM('Pendente', 'Concluída', 'Cancelada');

CREATE TABLE public.sales(
  id SERIAL PRIMARY KEY,
  customer_id integer NOT NULL,
  customer varchar NOT NULL,
  sale_value NUMERIC(12,2) NOT NULL,
  status sale_shopping_status DEFAULT 'Pendente',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_sales_customer_id_foreign FOREIGN KEY (customer_id) REFERENCES public.customers(id) 
);

CREATE INDEX idx_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_status_created_at ON sales(status, created_at DESC);
CREATE INDEX idx_sales_status ON sales(status);
CREATE INDEX idx_sales_created_at ON sales(created_at DESC);

CREATE TABLE public.sale_itens(
  id SERIAL PRIMARY KEY,
  product_id integer NOT NULL,
  name character varying NOT NULL,
  qtde integer NOT NULL,
  sale_value double precision NOT NULL,
  sale_id bigint NOT NULL,
  status sale_shopping_itens_status DEFAULT 'Pendente',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_sales_itens_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT fk_sales_itens_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);

CREATE INDEX idx_sale_itens_sale_id ON sale_itens(sale_id);
CREATE INDEX idx_sale_itens_product_id ON sale_itens(product_id);
CREATE INDEX idx_sale_itens_sale_id_status ON sale_itens(sale_id, status);

CREATE TABLE public.sale_itens_grid(
  id SERIAL PRIMARY KEY,
  product_id integer NULL,
  sale_id integer NULL,
  status sale_shopping_itens_status DEFAULT 'Pendente',
  product_grid_id integer NULL,
  size_saled sizes NULL,
  grid_qtde integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_sale_itens_grid_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT fk_sale_itens_grid_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id),
  CONSTRAINT fk_sale_itens_grid_product_grids_id_foreign FOREIGN KEY (product_grid_id) REFERENCES public.product_grids(id)
);

ALTER TABLE sale_itens_grid
ADD CONSTRAINT unique_sale_grid
UNIQUE (sale_id, size_saled, product_grid_id);

CREATE INDEX idx_sale_itens_grid_product_id ON sale_itens_grid(product_id);
CREATE INDEX idx_sale_itens_grid_product_grid_id ON sale_itens_grid(product_grid_id);
CREATE TABLE idx_sale_itens_grid_status ON sale_itens_grid(status);

CREATE TABLE public.sale_pay_ment (
  id SERIAL PRIMARY KEY,
  sale_id bigint NOT NULL,
  specie_id bigint NOT NULL,
  specie character varying NOT NULL,
  amount_paid double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_sale_pay_ment_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT fk_sale_pay_ment_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);

CREATE INDEX idx_sale_payment_sale_id ON sale_pay_ment(sale_id);
CREATE INDEX idx_sale_payment_specie_id ON sale_pay_ment(specie_id);

-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 

CREATE TABLE public.shopping (
  id SERIAL PRIMARY KEY,
  load integer unique,
  operation character varying,
  status sale_shopping_status DEFAULT 'Pendente',
  total_shopping NUMERIC(12,2) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_shopping_created_at ON shopping(created_at DESC);

CREATE TABLE public.shopping_itens (
  id SERIAL PRIMARY KEY,
  shopping_id integer,
  product_id integer,
  name character varying NOT NULL,
  qtde_purchased integer NOT NULL,
  purchased_value double precision NOT NULL,
  status sale_shopping_itens_status DEFAULT 'Pendente',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_shopping_itens_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT fk_shopping_itens_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id)
);

CREATE INDEX idx_shopping_itens_shopping_id ON shopping_itens(shopping_id);
CREATE INDEX idx_shopping_itens_product_id ON shopping_itens(product_id);
CREATE INDEX idx_shopping_itens_shopping_idstatus ON shopping_itens(shopping_id, status);

CREATE TABLE public.shopping_itens_grid(
  id SERIAL PRIMARY KEY,
  product_id integer NULL,
  shopping_id integer NULL,
  status sale_shopping_itens_status DEFAULT 'Pendente',
  product_grid_id integer NULL,
  size_saled sizes NULL,
  grid_qtde integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_shopping_itens_grid_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT fk_shopping_itens_grid_products_id_foreign FOREIGN KEY (product_id) REFERENCES public.products(id),
  CONSTRAINT fk_shopping_itens_grid_product_grids_id_foreign FOREIGN KEY (product_grid_id) REFERENCES public.product_grids(id)
);

CREATE INDEX idx_shopping_itens_grid_product_id ON shopping_itens_grid(product_id);
CREATE INDEX idx_shopping_itens_grid_product_grid_id ON shopping_itens_grid(product_grid_id);
CREATE TABLE idx_shopping_itens_grid_status ON shopping_itens_grid(status);

ALTER TABLE shopping_itens_grid
ADD CONSTRAINT unique_shopping_grid
UNIQUE (shopping_id, size_saled, product_grid_id);

CREATE TABLE public.shopping_pay_ment(
  id SERIAL PRIMARY KEY,
  shopping_id bigint NOT NULL,
  specie_id bigint NOT NULL,
  specie character varying NOT NULL,
  amount_paid double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT fk_shopping_pay_ment_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id),
  CONSTRAINT fk_shopping_pay_ment_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id)
);

CREATE INDEX idx_shopping_payment_shopping_id ON shopping_pay_ment(shopping_id);
CREATE INDEX idx_shopping_payment_specie_id ON shopping_pay_ment(specie_id);

-- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## -- ## 
CREATE TABLE public.cash_registers (
  id SERIAL PRIMARY KEY,
  sale_id integer NULL,
  shopping_id integer NULL,
  description character varying NOT NULL,
  customer_id integer,
  customer varchar NULL,
  specie_id integer NOT NULL,
  specie character varying NOT NULL,
  input_value NUMERIC(12,2) NOT NULL,
  output_value NUMERIC(12,2) NOT NULL,
  total_balance NUMERIC(12,2) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  
  CONSTRAINT fk_cash_registers_sale_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
  CONSTRAINT fk_cash_registers_shopping_id_foreign FOREIGN KEY (shopping_id) REFERENCES public.shopping(id),
  CONSTRAINT fk_cash_registers_customer_id_foreign FOREIGN KEY (customer_id) REFERENCES public.customers(id),
  CONSTRAINT fk_cash_registers_specie_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
);

CREATE INDEX idx_cash_registers_customer_id ON cash_registers(customer_id);
CREATE INDEX idx_cash_registers_sale_id ON cash_registers(sale_id);
CREATE INDEX idx_cash_registers_shopping_id ON cash_registers(shopping_id);
CREATE INDEX idx_cash_registers_specie_id ON cash_registers(specie_id);
CREATE INDEX idx_cash_registers_created_at ON cash_registers(created_at DESC);