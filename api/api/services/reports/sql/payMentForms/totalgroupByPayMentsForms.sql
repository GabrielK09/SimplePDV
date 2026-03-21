WITH resumo AS (
  SELECT 
    sp.specie AS especie,
    SUM(COALESCE(sp.amount_paid, 0)) AS total_paid

  FROM
    sales s

  INNER JOIN
    sale_pay_ment sp ON sp.specie_id = s.id
  
  WHERE
    s.status = 'Concluída' 
    AND s.created_at::DATE >= $1 
    AND s.created_at::DATE <= $2

  GROUP BY
    sp.specie
)

SELECT
  r.especie,
  r.total_paid

FROM
  resumo r;