package cheats

import (
	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/ricochhet/dbmod/pkg/strx"
	"github.com/tidwall/gjson"
)

var weaponSkinFilter = struct {
	ItemTypes               map[string]struct{}
	EndsWith                []string
	Icons                   map[string]struct{}
	Descriptions            map[string]struct{}
	WhitelistedDescriptions map[string]struct{}
}{
	ItemTypes: map[string]struct{}{
		// Internal / base items.
		"/Lotus/Upgrades/Skins/Effects/BaseFootsteps":                    {},
		"/Lotus/Upgrades/Skins/Operator/AnimationSets/BaseOperatorAnims": {},
		// Unused auxiliary cosmetics.
		"/Lotus/Upgrades/Skins/Halos/PrototypeRaidHalo": {},
		// Unreleased / unfinished cosmetics.
		"/Lotus/Upgrades/Skins/Operator/Hair/HairAdultNightwave":  {},
		"/Lotus/Upgrades/Skins/Operator/Hair/HairAdultNightwaveB": {},
		"/Lotus/Upgrades/Skins/Promo/ChangYou/CYSingleStaffSkin":  {},
		// Default (base) cosmetics.
		"/Lotus/Types/Game/InfestedKavatPet/Patterns/InfestedCritterPatternDefault":     {},
		"/Lotus/Types/Game/InfestedPredatorPet/Patterns/InfestedPredatorPatternDefault": {},
		"/Lotus/Upgrades/Skins/Excalibur/ExcaliburPrimeAlabasterSkin":                   {},
		"/Lotus/Upgrades/Skins/Saryn/WF1999SarynHelmet":                                 {},
	},
	EndsWith: []string{"WingsRight", "WingsStaticRight"},
	Icons: map[string]struct{}{
		"/Lotus/Interface/Icons/StoreIcons/Resources/CraftingComponents/GenericWarframeHelmet.png": {},
	},
	Descriptions: map[string]struct{}{
		"/Lotus/Language/Items/GenericSuitCustomizationDesc": {},
	},
	WhitelistedDescriptions: map[string]struct{}{
		"/Lotus/Language/Items/GenericOperatorHairDescription": {},
		"/Lotus/Language/Operator/DrifterBeardDesc":            {},
	},
}

func ApplyWeaponSkins(customs, inventory []byte, index int) ([]byte, error) {
	skins, err := jsonx.ArrayElementFieldValues(inventory, "WeaponSkins", index)
	if err != nil {
		return nil, errorx.New("jsonx.ResultAsArray", err)
	}

	exportCustoms := gjson.ParseBytes(customs).Map()
	seen := make(map[string]struct{}, len(skins))
	combined := []string{}

	for _, skin := range skins {
		itemType := skin.Get("ItemType").String()
		if _, blacklisted := weaponSkinFilter.ItemTypes[itemType]; blacklisted {
			logx.Infof("Skipping blacklisted WeaponSkin: %s\n", itemType)
			continue
		}

		seen[itemType] = struct{}{}

		combined = append(combined, skin.Raw)
	}

	for uniqueName, item := range exportCustoms {
		if item.Get("name").String() == "" {
			if _, whitelisted := weaponSkinFilter.WhitelistedDescriptions[item.Get("description").String()]; !whitelisted {
				weaponSkinFilter.ItemTypes[uniqueName] = struct{}{}
				continue
			}
		}

		if _, blacklisted := weaponSkinFilter.Descriptions[item.Get("description").String()]; blacklisted {
			weaponSkinFilter.ItemTypes[uniqueName] = struct{}{}
			continue
		}

		if _, blacklisted := weaponSkinFilter.Icons[item.Get("icon").String()]; blacklisted {
			weaponSkinFilter.ItemTypes[uniqueName] = struct{}{}
			continue
		}

		if strx.EndsWithAny(uniqueName, weaponSkinFilter.EndsWith) {
			weaponSkinFilter.ItemTypes[uniqueName] = struct{}{}
			continue
		}

		if _, exists := seen[uniqueName]; exists {
			continue
		}

		seen[uniqueName] = struct{}{}

		weaponSkin, err := database.NewWeaponSkin(uniqueName)
		if err != nil {
			return nil, errorx.New("database.NewWeaponSkin", err)
		}

		combined = append(combined, weaponSkin)
	}

	result := combined[:0]
	for _, raw := range combined {
		itemType := gjson.Get(raw, "ItemType").String()
		if _, blacklisted := weaponSkinFilter.ItemTypes[itemType]; blacklisted {
			logx.Infof("Skipping blacklisted WeaponSkin: %s\n", itemType)
			continue
		}

		result = append(result, raw)
	}

	newInventory, err := jsonx.SetArrayElementFieldArray(inventory, "WeaponSkins", result, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetSliceInRawBytes", err)
	}

	logx.Infof("Original: %d, Added: %d, Final: %d\n",
		len(skins), len(result)-len(skins), len(result))

	return newInventory, nil
}
