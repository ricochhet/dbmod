package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

type CodexScans struct {
	MaxScans int64
}

func (c *CodexScans) Apply(custom, codex, enemies, stats []byte, index int) ([]byte, error) {
	scans, err := jsonx.ArrayElementFieldValues(stats, "Scans", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportCustomScans := gjson.ParseBytes(custom).Array()
	exportCodex := gjson.ParseBytes(codex).Get("objects").Map()
	exportEnemies := gjson.ParseBytes(enemies).Get("avatars").Map()
	seen := make(map[string]struct{}, len(scans))
	combined := []string{}

	for _, scan := range scans {
		t := scan.Get("type").String()

		changed, err := database.NewScan(t, max(scan.Get("scans").Int(), c.MaxScans))
		if err != nil {
			return nil, errorx.New("database.NewScan", err)
		}

		if scan.Raw != changed {
			logx.Infof("Old: %v, New: %v\n", scan.Raw, changed)
		}

		combined = append(combined, changed)
		seen[t] = struct{}{}
	}

	addScan := func(t string) error {
		if _, exists := seen[t]; exists {
			return nil
		}

		s, err := database.NewScan(t, c.MaxScans)
		if err != nil {
			return errorx.New("database.NewScan", err)
		}

		combined = append(combined, s)
		seen[t] = struct{}{}

		return nil
	}

	for _, scan := range exportCustomScans {
		if err := addScan(scan.String()); err != nil {
			return nil, err
		}
	}

	for t := range exportCodex {
		if err := addScan(t); err != nil {
			return nil, err
		}
	}

	for t := range exportEnemies {
		if err := addScan(t); err != nil {
			return nil, err
		}
	}

	newStats, err := jsonx.SetArrayElementFieldArray(stats, "Scans", combined, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(scans), len(combined)-len(scans), len(combined))

	return newStats, nil
}
