package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

var missionBlacklist = map[string]struct{}{
	"PvpNode0":     {},
	"PvpNode9":     {},
	"PvpNode10":    {},
	"MercuryHUB":   {},
	"EarthHUB":     {},
	"TradeHUB1":    {},
	"SaturnHUB":    {},
	"EventNode763": {}, // The Index: Endurance (unused).
	"PlutoHUB":     {},
	"ZarimanHub":   {},
	"SolNode234":   {}, // Dormizone.
}

func ApplyMissions(regions, inventory []byte, index int) ([]byte, error) {
	missions, err := jsonx.ArrayElementFieldValues(inventory, "Missions", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportRegions := gjson.ParseBytes(regions).Map()
	seen := make(map[string]struct{}, len(missions))
	combined := []string{}

	for _, item := range missions {
		seen[item.Get("Tag").String()] = struct{}{}
		combined = append(combined, item.Raw)
	}

	for uniqueName := range exportRegions {
		if _, exists := seen[uniqueName]; exists {
			continue
		}

		if _, blacklisted := missionBlacklist[uniqueName]; blacklisted {
			logx.Infof("Skipping blacklisted Mission: %s\n", uniqueName)
			continue
		}

		node, err := database.NewNode(uniqueName, 1, 1)
		if err != nil {
			return nil, errorx.New("database.NewNode", err)
		}

		combined = append(combined, node)
		seen[uniqueName] = struct{}{}
	}

	result := combined[:0]
	for _, raw := range combined {
		tag := gjson.Get(raw, "Tag").String()
		if _, blacklisted := missionBlacklist[tag]; blacklisted {
			logx.Infof("Skipping blacklisted Mission: %s\n", tag)
			continue
		}

		result = append(result, raw)
	}

	newInventory, err := jsonx.SetArrayElementFieldArray(inventory, "Missions", result, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(missions), len(result)-len(missions), len(result))

	return newInventory, nil
}
