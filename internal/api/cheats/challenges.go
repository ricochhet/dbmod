package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

func ApplyChallenges(achievements, inventory []byte, index int) ([]byte, error) {
	challengeProgress, err := jsonx.ArrayElementFieldValues(inventory, "ChallengeProgress", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportAchievements := gjson.ParseBytes(achievements).Map()
	seen := make(map[string]struct{}, len(challengeProgress))
	combined := []string{}

	for _, entry := range challengeProgress {
		seen[entry.Get("Name").String()] = struct{}{}
		combined = append(combined, entry.Raw)
	}

	for uniqueName, item := range exportAchievements {
		if _, ok := seen[uniqueName]; ok {
			continue
		}

		requiredCount := item.Get("requiredCount").Int()
		if !item.Get("requiredCount").Exists() {
			requiredCount = 1
		}

		challenge, err := database.NewChallengeProgress(requiredCount, uniqueName)
		if err != nil {
			return nil, errorx.New("database.NewChallengeProgress", err)
		}

		combined = append(combined, challenge)
	}

	newInventory, err := jsonx.SetArrayElementFieldArray(
		inventory,
		"ChallengeProgress",
		combined,
		index,
	)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(challengeProgress), len(combined)-len(challengeProgress), len(combined))

	return newInventory, nil
}
