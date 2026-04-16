package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

type EnemyStats struct {
	Kills     int64
	Assists   int64
	Headshots int64
}

func (e *EnemyStats) Apply(enemies, stats []byte, index int) ([]byte, error) {
	enemyStats, err := jsonx.ArrayElementFieldValues(stats, "Enemies", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportEnemies := gjson.ParseBytes(enemies).Get("avatars").Map()
	seen := make(map[string]struct{}, len(enemyStats))
	combined := []string{}

	for _, enemy := range enemyStats {
		itemType := enemy.Get("type").String()
		seen[itemType] = struct{}{}

		changed, err := database.NewEnemy(itemType,
			max(enemy.Get("kills").Int(), e.Kills),
			max(enemy.Get("assists").Int(), e.Assists),
			max(enemy.Get("headshots").Int(), e.Headshots),
			max(0, enemy.Get("captures").Int()),
			max(0, enemy.Get("executions").Int()),
			max(0, enemy.Get("deaths").Int()),
		)
		if err != nil {
			return nil, errorx.New("database.NewEnemy", err)
		}

		if enemy.Raw != changed {
			logx.Infof("Old: %v, New: %v\n", enemy.Raw, changed)
		}

		combined = append(combined, changed)
	}

	for uniqueName := range exportEnemies {
		if _, exists := seen[uniqueName]; exists {
			continue
		}

		enemy, err := database.NewEnemy(uniqueName, e.Kills, e.Assists, e.Headshots, 0, 0, 0)
		if err != nil {
			return nil, errorx.New("database.NewEnemy", err)
		}

		combined = append(combined, enemy)
		seen[uniqueName] = struct{}{}
	}

	newStats, err := jsonx.SetArrayElementFieldArray(stats, "Enemies", combined, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(enemyStats), len(combined)-len(enemyStats), len(combined))

	return newStats, nil
}
