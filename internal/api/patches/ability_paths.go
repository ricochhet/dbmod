package patches

import (
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
	"github.com/ricochhet/dbmod/pkg/strx"
	"github.com/tidwall/gjson"
)

func ApplyU42AbilityPaths(warframesU41, warframesU42, inventory []byte, index int) ([]byte, error) {
	var err error

	full, err := jsonx.ArrayElement(inventory, index)
	if err != nil {
		return nil, errorx.WithFrame(err)
	}

	newInventory := full.Raw
	matches := jsonx.Search(newInventory, jsonx.SearchCriteria{
		Key: "AbilityOverride",
		Fields: []jsonx.FieldCriteria{
			{Key: "Ability"},
			{Key: "Index"},
		},
		MatchAll: true,
	})

	abilitiesU41 := GetAbilityPaths(warframesU41)
	abilitiesU42 := GetAbilityPaths(warframesU42)
	replacements := MigrateAbilityPaths(abilitiesU41, abilitiesU42)

	changed := 0

	for _, m := range matches {
		abilityPath := m.Path + ".Ability"

		current := jsonx.Search(newInventory, jsonx.SearchCriteria{
			Key:   "Ability",
			Value: "",
		})

		currentAbilityName := ""

		for _, a := range current {
			if a.Path == abilityPath {
				currentAbilityName = a.Value.String()
				break
			}
		}

		newAbilityName := strx.ReplaceMap(currentAbilityName, replacements)
		if newAbilityName == currentAbilityName {
			logx.Infof("Skipping ability: %s (already updated)\n", currentAbilityName)
			continue
		}

		newInventory, err = jsonx.Edit(newInventory, abilityPath, newAbilityName)
		if err != nil {
			return nil, errorx.New("jsonx.Edit", err)
		}

		changed++
	}

	logx.Infof("Updated: %d\n", changed)

	return jsonx.SetArrayElementRaw(inventory, []byte(newInventory), index)
}

func MigrateAbilityPaths(u41, u42 map[string][]string) map[string]string {
	result := make(map[string]string)

	for uniqueName, abilitiesU41 := range u41 {
		abilitiesU42, ok := u42[uniqueName]
		if !ok {
			continue
		}

		n := min(len(abilitiesU42), len(abilitiesU41))
		for i := range n {
			pathU41 := abilitiesU41[i]
			pathU42 := abilitiesU42[i]

			if pathU41 != "" && pathU42 != "" && pathU41 != pathU42 {
				result[pathU41] = pathU42
			}
		}
	}

	return result
}

func GetAbilityPaths(resources []byte) map[string][]string {
	result := map[string][]string{}

	exportWarframes := gjson.ParseBytes(resources).Map()
	for uniqueName, item := range exportWarframes {
		abilities := item.Get("abilities").Array()
		if len(abilities) == 0 {
			continue
		}

		list := make([]string, 0, len(abilities))
		for _, ability := range abilities {
			name := ability.Get("uniqueName").String()
			if name == "" {
				continue
			}

			list = append(list, name)
		}

		if len(list) > 0 {
			result[uniqueName] = list
		}
	}

	return result
}
