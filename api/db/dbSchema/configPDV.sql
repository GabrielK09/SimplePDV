CREATE TABLE config_pdv (
    id SERIAL PRIMARY KEY,
    confirm_to_pinter BOOLEAN DEFAULT 'false',
    block_sale_negative_stock BOOLEAN DEFAULT 'false',
    reserve_stock BOOLEAN DEFAULT 'true'
)