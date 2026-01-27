CREATE TABLE pay_ment_forms (
    id SERIAL PRIMARY KEY,
    sale_id bigint NOT NULL,
    specie character varying NOT NULL,
    amount_paid FLOAT NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT pay_ment_forms_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id)
)