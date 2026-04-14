package cheats

import (
	"fmt"

	"github.com/ricochhet/dbmod/internal/database"
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/ricochhet/dbmod/pkg/logx"
)

type Accolades struct {
	Staff     bool
	Founder   int64
	Guide     int64
	Moderator bool
	Partner   bool
	Heirloom  bool
	Counselor bool
}

func (a *Accolades) Apply(inventory []byte, index int) ([]byte, error) {
	newAccolades, err := database.NewAccolades(a.Heirloom)
	if err != nil {
		return nil, errorx.New("database.NewAccolades", err)
	}

	result, err := jsonx.SetFieldInRawBytes(inventory, "Accolades", newAccolades, index)
	if err != nil {
		return nil, errorx.New("jsonx.SetFieldInRawBytes", err)
	}

	fields := map[string]any{
		"Staff":     a.Staff,
		"Founder":   a.Founder,
		"Guide":     a.Guide,
		"Moderator": a.Moderator,
		"Partner":   a.Partner,
		"Counselor": a.Counselor,
	}

	newInventory := string(result)

	for k, v := range fields {
		logx.Infof("Added: %s: %v\n", k, v)

		newInventory, err = jsonx.SetFieldInBytes(newInventory, k, v, index)
		if err != nil {
			return nil, errorx.New(fmt.Sprintf("jsonx.SetFieldInBytes (%s)", k), err)
		}
	}

	return []byte(newInventory), nil
}
