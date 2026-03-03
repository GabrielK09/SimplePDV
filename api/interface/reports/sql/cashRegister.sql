WITH resumo AS (
    SELECT
        cr.description AS descricao,
        cr.customer AS cliente,
        cr.specie AS especie,
        cr.input_value AS valorEntrada,
        cr.output_value AS valoraSaida
        
    FROM
        cash_registers cr

    WHERE
        cr.created_at::DATE >= $1 
        AND cr.created_at::DATE < $2

), tot AS (
    SELECT
        'TOTAL'::text AS descricao,
        ''::text AS cliente,
        ''::text AS especie,
        0 AS valorEntrada,
        0 AS valorSaida,
        SUM(COALESCE(cr.input_value, 0)) AS totalEntrada,
        SUM(COALESCE(cr.output_value, 0)) AS totalSaida
    FROM
        cash_registers cr
    WHERE
        cr.created_at::DATE >= $1 
        AND cr.created_at::DATE < $2
)

SELECT
    r.descricao,
    r.cliente,
    r.especie,
    r.valorEntrada,
    r.valoraSaida,
    0 AS totalEntrada,
    0 AS totalSaida

FROM
    resumo r

UNION ALL

SELECT  
    t.descricao,
    t.cliente,
    t.especie,
    t.valorEntrada,
    t.valorSaida,
    t.totalEntrada,
    t.totalSaida

FROM    
    tot t;