WITH resumo AS (
    SELECT
        s.load AS load,
        COALESCE(SUM(s.total_shopping), 0) AS totalShopping

    FROM
        shopping s

    WHERE
        s.status = 'Concluída'
        AND s.created_at::DATE >= $1 
        AND s.created_at::DATE <= $2

    GROUP BY
        s.load
)

SELECT
    r.load,
    r.totalShopping
FROM
    resumo r