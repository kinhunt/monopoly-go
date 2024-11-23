// internal/game/map.go
package game

import (
	"monopoly/pkg/utils"
)

type GameMap struct {
	Tiles []*Tile `json:"tiles"`
}

func NewDefaultMap() *GameMap {
	tiles := []*Tile{
		{
			ID:   0,
			Name: "起点",
			Type: TileStart,
		},
		{
			ID:        1,
			Name:      "第一大道",
			Type:      TileProperty,
			Price:     200,
			RentPrice: []int{20, 40, 80},
		},
		{
			ID:   2,
			Name: "机会",
			Type: TileChance,
		},
		{
			ID:        3,
			Name:      "第二大道",
			Type:      TileProperty,
			Price:     300,
			RentPrice: []int{30, 60, 120},
		},
		{
			ID:   4,
			Name: "银行",
			Type: TileBank,
		},
	}

	return &GameMap{Tiles: tiles}
}

func (m *GameMap) GetTile(position int) (*Tile, error) {
	if position < 0 || position >= len(m.Tiles) {
		return nil, utils.ErrInvalidInput
	}
	return m.Tiles[position], nil
}
