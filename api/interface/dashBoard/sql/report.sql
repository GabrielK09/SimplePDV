WITH resumo AS (
    SELECT
        COALESCE(SUM(s.sale_value), 0.0) AS total_saled,
        COALESCE(SUM(s.sale_value * COALESCE(p.commission / 100, 1)), 0.0) AS commission

    FROM 
        sales s

    INNER JOIN
        sale_itens si ON si.sale_id = s.id

    INNER JOIN
        products p ON p.id = si.product_id

    WHERE
        s.status = 'Concluída' 
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE < $2

), best_customers AS (
    SELECT
        s.customer AS best_customer,
        COUNT(si.product_id) AS amount_saled

    FROM    
        sales s

    INNER JOIN
        sale_itens si ON si.sale_id = s.id

    WHERE
        s.status = 'Concluída' 
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE < $2

    GROUP BY
        s.customer

    ORDER BY
        COALESCE(SUM(s.sale_value), 0.0) DESC 
    LIMIT 1
)
    
SELECT
    r.total_saled,
    r.commission,
    COALESCE(mc.best_customer, '') AS best_customer,
    COALESCE(mc.amount_saled, 0) AS amount_saled


FROM
    resumo r

LEFT JOIN
    best_customers mc ON true;