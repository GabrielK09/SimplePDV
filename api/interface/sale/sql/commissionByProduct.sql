SELECT
    p.name AS "Produto",
    sl.sale_value AS "Valor da venda",
    p.commission AS "Comissão do produto",
    COALESCE(SUM(sl.sale_value * COALESCE(p.commission / 100, 1)), 0.0) AS "Comissão gerada pelo produto"

FROM
    sale_itens sl

INNER JOIN
    products p ON p.id = sl.product_id

WHERE
    sl.sale_id = $1
    AND p.commission > 0

GROUP BY
    p.name,
    sl.sale_value,
    p.commission