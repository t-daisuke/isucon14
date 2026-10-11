USE isuride;

-- 初期データ投入後に、既存の位置情報から椅子の総移動距離を一度だけ計算して埋める
INSERT INTO chair_total_distances (chair_id, total_distance, updated_at)
SELECT chair_id,
       IFNULL(SUM(distance), 0),
       MAX(created_at)
FROM (
  SELECT chair_id,
         created_at,
         ABS(latitude - LAG(latitude) OVER w) + ABS(longitude - LAG(longitude) OVER w) AS distance
  FROM chair_locations
  WINDOW w AS (PARTITION BY chair_id ORDER BY created_at)
) t
GROUP BY chair_id;
