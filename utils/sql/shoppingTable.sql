CREATE TABLE shopping (
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    load INT,
    operation VARCHAR(120),
    date_of_purchase timestamp without time zone,
)