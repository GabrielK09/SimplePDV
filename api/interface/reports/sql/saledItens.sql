WITH resumo AS (
    SELECT
        si.sale_id AS sale_id,
        si.product_id AS product_id,
        si.name AS produto,
        si.sale_value AS item_sale_value,
        si.qtde AS qtde

    FROM 
        sales s

    INNER JOIN
        sale_itens si ON si.sale_id = s.id

    WHERE
        s.status = 'Concluída' 
        AND si.status = 'Concluída'
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE < $2

)

SELECT
    sale_id,
    product_id,
    produto,
    item_sale_value,
    qtde

FROM
    resumo;