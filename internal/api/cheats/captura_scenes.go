package cheats

import (
	"github.com/ricochhet/dbmod/internal/config"
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

func ApplyCapturaScenes(resources, virtuals, inventory []byte, index int) ([]byte, error) {
	miscItems, err := jsonx.ResultAsArray(inventory, "MiscItems", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportResources := gjson.ParseBytes(resources).Map()
	exportVirtuals := gjson.ParseBytes(virtuals).Map()
	seen := make(map[string]struct{}, len(miscItems))
	combined := []string{}

	for _, item := range miscItems {
		seen[item.Get("ItemType").String()] = struct{}{}
		combined = append(combined, item.Raw)
	}

	parents := make(map[string]string)

	for name, node := range exportResources {
		if p := node.Get("parentName").String(); p != "" {
			parents[name] = p
		}
	}

	for name, node := range exportVirtuals {
		if _, ok := parents[name]; !ok {
			if p := node.Get("parentName").String(); p != "" {
				parents[name] = p
			}
		}
	}

	const photoboothTile = "/Lotus/Types/Items/MiscItems/PhotoboothTile"
	for name := range exportResources {
		if !config.ResourceInheritsFromMap(parents, name, photoboothTile) {
			continue
		}

		if _, exists := seen[name]; exists {
			continue
		}

		item, err := database.NewItem(name, 1)
		if err != nil {
			return nil, errorx.New("database.NewItem", err)
		}

		combined = append(combined, item)
		seen[name] = struct{}{}
	}

	newInventory, err := jsonx.SetSliceInRawBytes(inventory, "MiscItems", combined, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(miscItems), len(combined)-len(miscItems), len(combined))

	return newInventory, nil
}
