package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

type ShipDecorations struct {
	MaxCount int64
}

func (s *ShipDecorations) Apply(resources, inventory []byte, index int) ([]byte, error) {
	decorations, err := jsonx.ResultAsArray(inventory, "ShipDecorations", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportResources := gjson.ParseBytes(resources).Map()
	seen := make(map[string]struct{}, len(decorations))
	combined := []string{}

	for _, deco := range decorations {
		itemType := deco.Get("ItemType").String()
		seen[itemType] = struct{}{}

		changed, err := database.NewItem(itemType, max(deco.Get("ItemCount").Int(), s.MaxCount))
		if err != nil {
			return nil, errorx.New("database.NewItem", err)
		}

		if deco.Raw != changed {
			logx.Infof("Old: %v, New: %v\n", deco.Raw, changed)
		}

		combined = append(combined, changed)
	}

	for uniqueName, item := range exportResources {
		if item.Get("productCategory").String() != "ShipDecorations" {
			continue
		}

		if _, exists := seen[uniqueName]; exists {
			continue
		}

		entry, err := database.NewItem(uniqueName, s.MaxCount)
		if err != nil {
			return nil, errorx.New("database.NewItem", err)
		}

		combined = append(combined, entry)
		seen[uniqueName] = struct{}{}
	}

	newInventory, err := jsonx.SetSliceInRawBytes(inventory, "ShipDecorations", combined, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(decorations), len(combined)-len(decorations), len(combined))

	return newInventory, nil
}
