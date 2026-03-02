CREATE TABLE pay_ment_forms (
    id SERIAL PRIMARY KEY,
    specie character varying NOT NULL,
    pix_key character varying NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
)