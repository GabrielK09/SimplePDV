WITH resumo AS (
    SELECT
        si.shopping_id AS shopping_id,
        si.product_id AS product_id,
        si.name AS produto,
        si.purchased_value AS purchased_value,
        si.qtde_purchased AS qtde_purchased

    FROM 
        shopping s

    INNER JOIN
        shopping_itens si ON si.shopping_id = s.id

    WHERE
        s.status = 'Concluída'
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE <= $2
)

SELECT
    shopping_id,
    product_id,
    produto,
    purchased_value,
    qtde_purchased

FROM
    resumo;