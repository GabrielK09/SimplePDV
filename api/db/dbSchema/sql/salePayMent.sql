CREATE TABLE sale_pay_ment (
    id SERIAL PRIMARY KEY,
    sale_id bigint NOT NULL,
    specie_id bigint NOT NULL,
    specie character varying NOT NULL,
    amount_paid FLOAT NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT pay_ment_forms_sales_id_foreign FOREIGN KEY (sale_id) REFERENCES public.sales(id),
    CONSTRAINT pay_ment_forms_pay_ment_forms_id_foreign FOREIGN KEY (specie_id) REFERENCES public.pay_ment_forms(id)
)

--pay_ment_forms