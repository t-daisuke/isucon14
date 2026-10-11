package main

// 空き椅子(有効で、未完了のライドに割り当てられていない椅子)と、その最新位置を返すSQL
// 呼び出し側で末尾に条件や並び順を足す
// 空きとみなすのは COMPLETED が椅子に通知済み(chair_sent_at)になってから。通知前に次のライドを渡すと椅子側でエラーになる
// ponytail: 使用中判定を毎回全ライドで行う。椅子・ライドが増えて遅くなったら椅子ごとの使用中フラグに変える
const freeChairsWithLatestLocationSQL = `
SELECT c.id, c.name, c.model, l.latitude, l.longitude
FROM chairs c
JOIN (
  SELECT chair_id, MAX(created_at) AS created_at
  FROM chair_locations
  GROUP BY chair_id
) m ON m.chair_id = c.id
JOIN chair_locations l ON l.chair_id = m.chair_id AND l.created_at = m.created_at
WHERE c.is_active = TRUE
  AND c.id NOT IN (
    SELECT r.chair_id
    FROM rides r
    WHERE r.chair_id IS NOT NULL
      AND NOT EXISTS (
        SELECT 1 FROM ride_statuses rs
        WHERE rs.ride_id = r.id AND rs.status = 'COMPLETED' AND rs.chair_sent_at IS NOT NULL
      )
  )
`

type freeChair struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Model     string `db:"model"`
	Latitude  int    `db:"latitude"`
	Longitude int    `db:"longitude"`
}
