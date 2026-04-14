//nolint:tagliatelle // match SNS database structure.
package database

import (
	"encoding/json"
	"fmt"

	"github.com/ricochhet/dbmod/pkg/cryptox"
	"github.com/ricochhet/dbmod/pkg/errorx"
)

type Accolades struct {
	Heirloom bool `json:"Heirloom"`
}

type Item struct {
	ItemType  string `json:"ItemType"`
	ItemCount int64  `json:"ItemCount,omitempty"`
}

type XPInfo struct {
	ItemType string `json:"ItemType"`
	XP       int64  `json:"XP"`
}

type Node struct {
	Tag       string `json:"Tag"`
	Completes int64  `json:"Completes,omitempty"`
	Tier      int64  `json:"Tier,omitempty"`
}

type Scan struct {
	Type  string `json:"type"`
	Scans int64  `json:"scans"`
}

type WeaponSkin struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	ItemType string `json:"ItemType"`
}

type ChallengeProgress struct {
	Progress  int64    `json:"Progress"`
	Name      string   `json:"Name"`
	Completed []string `json:"Completed"`
}

type Enemy struct {
	Type       string `json:"type"`
	Kills      int64  `json:"kills"`
	Assists    int64  `json:"assists"`
	Headshots  int64  `json:"headshots"`
	Captures   int64  `json:"captures,omitempty"`
	Executions int64  `json:"executions,omitempty"`
	Deaths     int64  `json:"deaths,omitempty"`
}

func NewAccolades(heirloom bool) (string, error) {
	b, err := json.Marshal(Accolades{Heirloom: heirloom})
	return string(b), errorx.WithFrame(err)
}

func NewItem(itemType string, itemCount int64) (string, error) {
	b, err := json.Marshal(Item{ItemType: itemType, ItemCount: itemCount})
	return string(b), errorx.WithFrame(err)
}

func NewXPInfo(itemType string, xp int64) (string, error) {
	b, err := json.Marshal(XPInfo{ItemType: itemType, XP: xp})
	return string(b), errorx.WithFrame(err)
}

func NewItemType(itemType string) (string, error) {
	b, err := json.Marshal(Item{ItemType: itemType})
	return string(b), errorx.WithFrame(err)
}

func NewNode(tag string, completes, tier int64) (string, error) {
	b, err := json.Marshal(Node{Tag: tag, Completes: completes, Tier: tier})
	return string(b), errorx.WithFrame(err)
}

func NewScan(t string, scans int64) (string, error) {
	b, err := json.Marshal(Scan{Type: t, Scans: scans})
	return string(b), errorx.WithFrame(err)
}

func NewWeaponSkin(itemType string) (string, error) {
	ws := WeaponSkin{ItemType: itemType}
	ws.ID.Oid = fmt.Sprintf("cb70cb70cb70cb70%08x", cryptox.CatBreadHash(itemType))
	b, err := json.Marshal(ws)

	return string(b), errorx.WithFrame(err)
}

func NewChallengeProgress(progress int64, name string) (string, error) {
	b, err := json.Marshal(ChallengeProgress{Progress: progress, Name: name})
	return string(b), errorx.WithFrame(err)
}

func NewEnemy(
	itemType string,
	kills, assists, headshots, captures, executions, deaths int64,
) (string, error) {
	b, err := json.Marshal(Enemy{
		Type:       itemType,
		Kills:      kills,
		Assists:    assists,
		Headshots:  headshots,
		Captures:   captures,
		Executions: executions,
		Deaths:     deaths,
	})

	return string(b), errorx.WithFrame(err)
}
