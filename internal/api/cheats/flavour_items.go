package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/tidwall/gjson"
)

var flavourItemsBlacklist = map[string]struct{}{
	// Kavat/Kubrow color templates.
	"/Lotus/Types/Game/KubrowPet/Colors/KubrowPetColorKavatBase":      {},
	"/Lotus/Types/Game/KubrowPet/Colors/KubrowPetColorKavatSecondary": {},
	"/Lotus/Types/Game/KubrowPet/Colors/KubrowPetColorKavatTertiary":  {},
	// Glyphs with default/broken data.
	"/Lotus/Types/StoreItems/AvatarImages/FanChannel/AvatarImageDramakins":    {},
	"/Lotus/Types/StoreItems/AvatarImages/FanChannel/AvatarImageSenastra":     {},
	"/Lotus/Types/StoreItems/AvatarImages/FanChannel/AvatarImageDesRPG":       {},
	"/Lotus/Types/StoreItems/AvatarImages/FanChannel/AvatarImageKacchi":       {},
	"/Lotus/Types/StoreItems/AvatarImages/FanChannel/AvatarImageLovinDaTacos": {},
	"/Lotus/Types/StoreItems/AvatarImages/AvatarImageCreatorWgrates":          {},
	"/Lotus/Types/StoreItems/AvatarImages/ImageConquera2022B":                 {},
	"/Lotus/Types/StoreItems/AvatarImages/ImageConquera2022C":                 {},
	"/Lotus/Types/StoreItems/AvatarImages/ImageConquera2022D":                 {},
	// Duplicate color palettes.
	"/Lotus/Types/StoreItems/SuitCustomizations/ColourPickerItemD": {},
}

func ApplyFlavourItems(flavor, inventory []byte, index int) ([]byte, error) {
	flavourItems, err := jsonx.ResultAsArray(inventory, "FlavourItems", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportFlavour := gjson.ParseBytes(flavor).Map()
	seen := make(map[string]struct{}, len(flavourItems))
	combined := []string{}

	for _, item := range flavourItems {
		seen[item.Get("ItemType").String()] = struct{}{}
		combined = append(combined, item.Raw)
	}

	for uniqueName, item := range exportFlavour {
		if item.Get("name").String() == "" {
			flavourItemsBlacklist[uniqueName] = struct{}{}
			continue
		}

		if _, blacklisted := flavourItemsBlacklist[uniqueName]; blacklisted {
			continue
		}

		if _, exists := seen[uniqueName]; exists {
			continue
		}

		itemType, err := database.NewItemType(uniqueName)
		if err != nil {
			return nil, errorx.New("database.NewItemType", err)
		}

		combined = append(combined, itemType)
		seen[uniqueName] = struct{}{}
	}

	result := combined[:0]
	for _, raw := range combined {
		itemType := gjson.Get(raw, "ItemType").String()
		if _, blacklisted := flavourItemsBlacklist[itemType]; blacklisted {
			logx.Infof("Skipping blacklisted FlavourItem: %s\n", itemType)
			continue
		}

		result = append(result, raw)
	}

	newInventory, err := jsonx.SetSliceInRawBytes(inventory, "FlavourItems", result, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(flavourItems), len(result)-len(flavourItems), len(result))

	return newInventory, nil
}
