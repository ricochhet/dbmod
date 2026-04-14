package config

import "github.com/tidwall/gjson"

func ResourceGetParent(resources, virtuals []byte, name string) string {
	path := name + ".parentName"
	if v := gjson.GetBytes(resources, path); v.Exists() {
		return v.String()
	}

	return gjson.GetBytes(virtuals, path).String()
}

func ResourceInheritsFrom(resources, virtuals []byte, name, targetName string) bool {
	for parent := ResourceGetParent(resources, virtuals, name); parent != ""; parent = ResourceGetParent(resources, virtuals, parent) {
		if parent == targetName {
			return true
		}
	}

	return false
}

func ResourceInheritsFromMap(parents map[string]string, name, targetName string) bool {
	for {
		parent, ok := parents[name]
		if !ok || parent == "" {
			return false
		}

		if parent == targetName {
			return true
		}

		name = parent
	}
}
