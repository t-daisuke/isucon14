USE isuride;

-- 初期データ投入後に、既存の位置情報から椅子の最新位置を一度だけ埋める
INSERT INTO chair_latest_locations (chair_id, latitude, longitude, created_at)
SELECT l.chair_id, l.latitude, l.longitude, l.created_at
FROM chair_locations l
JOIN (
  SELECT chair_id, MAX(created_at) AS created_at
  FROM chair_locations
  GROUP BY chair_id
) m ON m.chair_id = l.chair_id AND m.created_at = l.created_at;
