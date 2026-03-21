CREATE TABLE shopping (
    id SERIAL PRIMARY KEY,
    load INT,
    operation VARCHAR(120),
    status character varying DEFAULT 'Pendente',
    created_at timestamp without time zone,
    updated_at timestamp without time zone
)