CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    customer character varying NOT NULL,
    cpf_cnpj character varying NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
)