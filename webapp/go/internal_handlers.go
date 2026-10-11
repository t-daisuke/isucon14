package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sync"
)

// ponytail: 1インスタンス前提の全体ロック。複数インスタンスにするなら DB 側の行ロックに変える
var matchingMu sync.Mutex

// このAPIをインスタンス内から一定間隔で叩かせることで、椅子とライドをマッチングさせる
func internalGetMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	matchingMu.Lock()
	defer matchingMu.Unlock()

	// 一度の呼び出しで待っているライドをまとめて処理する (最大10件)
	for range 10 {
		ride := &Ride{}
		if err := db.GetContext(ctx, ride, `SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at LIMIT 1`); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		// 配車地点に一番近い空き椅子を選ぶ
		matched := &freeChair{}
		err := db.GetContext(ctx, matched, freeChairsWithLatestLocationSQL+`
ORDER BY ABS(l.latitude - ?) + ABS(l.longitude - ?)
LIMIT 1`, ride.PickupLatitude, ride.PickupLongitude)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", matched.ID, ride.ID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
