package patches

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

func ApplyShipDecorations(resources, inventory []byte, index int) ([]byte, error) {
	decorations, err := jsonx.ArrayElementFieldValues(inventory, "ShipDecorations", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportResources := gjson.ParseBytes(resources).Map()

	// Build a map from prop-item type → canonical resource name.
	propToResource := make(map[string]string, len(exportResources))
	for uniqueName, item := range exportResources {
		if deco := item.Get("deco").String(); deco != "" {
			propToResource[deco] = uniqueName
		}
	}

	blacklisted := make(map[string]struct{})
	accumulated := make(map[string]int64)

	for _, deco := range decorations {
		itemType := deco.Get("ItemType").String()
		if uniqueName, found := propToResource[itemType]; found {
			blacklisted[itemType] = struct{}{}
			accumulated[uniqueName] += deco.Get("ItemCount").Int()
			logx.Infof("%d : %s\n", deco.Get("ItemCount").Int(), itemType)
		}
	}

	result := []string{}

	for _, deco := range decorations {
		itemType := deco.Get("ItemType").String()
		if _, skip := blacklisted[itemType]; skip {
			logx.Infof("Skipping blacklisted ShipDecoration: %s\n", itemType)
			continue
		}

		result = append(result, deco.Raw)
	}

	for uniqueName, count := range accumulated {
		item, err := database.NewItem(uniqueName, count)
		if err != nil {
			return nil, errorx.New("database.NewItem", err)
		}

		result = append(result, item)
	}

	newInventory, err := jsonx.SetArrayElementFieldArray(
		inventory,
		"ShipDecorations",
		result,
		index,
	)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(decorations), len(result)-len(decorations), len(result))

	return newInventory, nil
}
