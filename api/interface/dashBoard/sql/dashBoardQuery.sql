WITH summary_sale AS (
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
        AND s.created_at::DATE <= $2

), summary_shopping AS (
    SELECT
        COALESCE(SUM(s.total_shopping), 0.0) AS total_shopping,
        COALESCE(SUM(si.qtde_purchased), 0) AS amount_shopping_itens,
        COUNT(s.id) AS amount_shopping

    FROM 
        shopping s

    INNER JOIN
        shopping_itens si ON si.shopping_id = s.id

    INNER JOIN
        products p ON p.id = si.product_id

    WHERE
        s.status = 'Concluída' 
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE <= $2

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
        AND s.created_at::DATE <= $2

    GROUP BY
        s.customer

    ORDER BY
        COALESCE(SUM(s.sale_value), 0.0) DESC 
    LIMIT 1
)
    
SELECT
    sr.total_saled,
    sr.commission,
    COALESCE(mc.amount_saled, 0) AS amount_saled,
    COALESCE(mc.best_customer, '') AS best_customer,

    s_shopping.total_shopping,
    COALESCE(s_shopping.amount_shopping_itens, 0) AS amount_shopping_itens,
    COALESCE(s_shopping.amount_shopping, 0) AS amount_shopping

FROM
    summary_sale sr

LEFT JOIN
    best_customers mc ON true

LEFT JOIN
    summary_shopping s_shopping ON true;