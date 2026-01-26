CREATE TABLE shopping (
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    load INT,
    operation VARCHAR(120),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
)