package database

import (
	"github.com/ricochhet/dbmod/pkg/errorx"
	"github.com/ricochhet/dbmod/pkg/jsonx"
	"github.com/tidwall/gjson"
)

type InventoryItems struct {
	LongGuns     []gjson.Result
	Pistols      []gjson.Result
	Melee        []gjson.Result
	Hoverboards  []gjson.Result
	OperatorAmps []gjson.Result
	MoaPets      []gjson.Result
	KubrowPets   []gjson.Result
}

func NewInventoryItems(inventory []byte, index int) (*InventoryItems, error) {
	fetch := func(field string) ([]gjson.Result, error) {
		r, err := jsonx.ArrayElementFieldValues(inventory, field, index)
		if err != nil {
			return nil, errorx.New("jsonx.ResultAsArray ("+field+")", err)
		}

		return r, nil
	}

	longGuns, err := fetch("LongGuns")
	if err != nil {
		return nil, err
	}

	pistols, err := fetch("Pistols")
	if err != nil {
		return nil, err
	}

	melee, err := fetch("Melee")
	if err != nil {
		return nil, err
	}

	hoverboards, err := fetch("Hoverboards")
	if err != nil {
		return nil, err
	}

	operatorAmps, err := fetch("OperatorAmps")
	if err != nil {
		return nil, err
	}

	moaPets, err := fetch("MoaPets")
	if err != nil {
		return nil, err
	}

	kubrowPets, err := fetch("KubrowPets")
	if err != nil {
		return nil, err
	}

	return &InventoryItems{
		LongGuns:     longGuns,
		Pistols:      pistols,
		Melee:        melee,
		Hoverboards:  hoverboards,
		OperatorAmps: operatorAmps,
		MoaPets:      moaPets,
		KubrowPets:   kubrowPets,
	}, nil
}
