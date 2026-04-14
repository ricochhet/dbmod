package patches

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/ricochhet/dbmod/pkg/strx"
	"github.com/tidwall/gjson"
)

var (
	xpInfoBlacklist = map[string]struct{}{}

	xpInfoWhitelist = map[string]struct{}{
		"/Lotus/Powersuits/Khora/Kavat/KhoraKavatPowerSuit":      {},
		"/Lotus/Powersuits/Khora/Kavat/KhoraPrimeKavatPowerSuit": {},
	}

	xpEarningParts = map[string]struct{}{
		"LWPT_BLADE":       {},
		"LWPT_GUN_BARREL":  {},
		"LWPT_AMP_OCULUS":  {},
		"LWPT_MOA_HEAD":    {},
		"LWPT_ZANUKA_HEAD": {},
		"LWPT_HB_DECK":     {},
	}

	xpInfoEndsWithBlacklist = []string{"PetWeapon"}
)

func ApplyXPInfo(weapons, warframes, sentinels, inventory []byte, index int) ([]byte, error) {
	xpInfo, err := jsonx.ResultAsArray(inventory, "XPInfo", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	inventoryItems, err := database.NewInventoryItems(inventory, index)
	if err != nil {
		return nil, errorx.New("database.NewInventoryItems", err)
	}

	exportWeapons := gjson.ParseBytes(weapons).Map()
	exportWarframes := gjson.ParseBytes(warframes).Map()
	exportSentinels := gjson.ParseBytes(sentinels).Map()

	seen := make(map[string]struct{}, len(xpInfo))
	combined := []string{}

	for _, item := range xpInfo {
		seen[item.Get("ItemType").String()] = struct{}{}
		combined = append(combined, item.Raw)
	}

	for _, item := range exportWarframes {
		for _, exalted := range item.Get("exalted").Array() {
			t := exalted.String()
			if _, whitelisted := xpInfoWhitelist[t]; !whitelisted {
				xpInfoBlacklist[t] = struct{}{}
			}
		}
	}

	combined = append(combined, collectXPInfo(seen, *inventoryItems, exportWeapons)...)

	result := []string{}

	seen2 := make(map[string]struct{}, len(combined))
	for _, raw := range combined {
		itemType := gjson.Get(raw, "ItemType").String()
		if _, blacklisted := xpInfoBlacklist[itemType]; blacklisted {
			logx.Infof("Skipping blacklisted XPInfo: %s\n", itemType)
			continue
		}

		if strx.EndsWithAny(itemType, xpInfoEndsWithBlacklist) {
			logx.Infof("Skipping blacklisted XPInfo: %s\n", itemType)
			continue
		}

		if exportWeapons[itemType].Index == 0 &&
			exportWarframes[itemType].Index == 0 &&
			exportSentinels[itemType].Index == 0 {
			logx.Infof("Unknown ItemType: %s\n", itemType)
		}

		if _, exists := seen2[itemType]; exists {
			continue
		}

		seen2[itemType] = struct{}{}

		result = append(result, raw)
	}

	newInventory, err := jsonx.SetSliceInRawBytes(inventory, "XPInfo", result, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(xpInfo), len(result)-len(xpInfo), len(result))

	return newInventory, nil
}

func collectXPInfo(
	seen map[string]struct{},
	items database.InventoryItems,
	weapons map[string]gjson.Result,
) []string {
	slots := [][]gjson.Result{
		items.LongGuns,
		items.Pistols,
		items.Melee,
		items.Hoverboards,
		items.OperatorAmps,
		items.MoaPets,
		items.KubrowPets,
	}

	var result []string //nolint:prealloc // really?
	for _, slot := range slots {
		result = append(result, xpInfoFromSlot(seen, slot, weapons)...)
	}

	return result
}

func xpInfoFromSlot(
	seen map[string]struct{},
	items []gjson.Result,
	weapons map[string]gjson.Result,
) []string {
	var result []string

	for _, item := range items {
		xp := item.Get("XP").Int()
		itemType := item.Get("ItemType").String()

		for _, modularPart := range item.Get("ModularParts").Array() {
			uniqueName := modularPart.String()
			part := weapons[uniqueName]
			partType := part.Get("partType").String()

			if partType == "" {
				continue
			}

			if _, earnsXP := xpEarningParts[partType]; !earnsXP {
				continue
			}

			xpInfoBlacklist[itemType] = struct{}{}

			if _, exists := seen[uniqueName]; exists {
				continue
			}

			entry, err := database.NewXPInfo(uniqueName, xp)
			if err != nil {
				continue
			}

			result = append(result, entry)
			seen[uniqueName] = struct{}{}
		}
	}

	return result
}
