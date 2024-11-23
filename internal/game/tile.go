// internal/game/tile.go
package game

import (
	"monopoly/pkg/utils"
)

// Tile 表示地图上的一个地块
type Tile struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Type      TileType `json:"type"`
	Price     int      `json:"price,omitempty"`     // 购买价格(仅地产类型)
	OwnerID   string   `json:"ownerId,omitempty"`   // 所有者ID(仅地产类型)
	Level     int      `json:"level,omitempty"`     // 升级等级(仅地产类型)
	RentPrice []int    `json:"rentPrice,omitempty"` // 各等级过路费(仅地产类型)
}

// NewTile 创建新的地块
func NewTile(id int, name string, tileType TileType) *Tile {
	return &Tile{
		ID:    id,
		Name:  name,
		Type:  tileType,
		Level: 0,
	}
}

// NewPropertyTile 创建新的地产类型地块
func NewPropertyTile(id int, name string, price int, rentPrices []int) *Tile {
	return &Tile{
		ID:        id,
		Name:      name,
		Type:      TileProperty,
		Price:     price,
		RentPrice: rentPrices,
		Level:     0,
	}
}

// GetRent 获取当前等级的过路费
func (t *Tile) GetRent() (int, error) {
	if t.Type != TileProperty {
		return 0, utils.ErrNotProperty
	}

	if t.OwnerID == "" {
		return 0, utils.ErrPropertyNotOwned
	}

	if t.Level >= len(t.RentPrice) {
		return 0, utils.ErrInvalidPropertyLevel
	}

	return t.RentPrice[t.Level], nil
}

// CanBeUpgraded 检查地产是否可以升级
func (t *Tile) CanBeUpgraded() bool {
	return t.Type == TileProperty && t.Level < len(t.RentPrice)-1
}

// GetUpgradeCost 获取升级费用
func (t *Tile) GetUpgradeCost() (int, error) {
	if t.Type != TileProperty {
		return 0, utils.ErrNotProperty
	}

	if !t.CanBeUpgraded() {
		return 0, utils.ErrMaxLevel
	}

	return t.Price / 2, nil
}

// Upgrade 升级地产
func (t *Tile) Upgrade() error {
	if t.Type != TileProperty {
		return utils.ErrNotProperty
	}

	if !t.CanBeUpgraded() {
		return utils.ErrMaxLevel
	}

	t.Level++
	return nil
}

// SetOwner 设置地产所有者
func (t *Tile) SetOwner(playerID string) error {
	if t.Type != TileProperty {
		return utils.ErrNotProperty
	}

	if t.OwnerID != "" {
		return utils.ErrPropertyOwned
	}

	t.OwnerID = playerID
	return nil
}

// GetValue 获取地产当前价值
func (t *Tile) GetValue() (int, error) {
	if t.Type != TileProperty {
		return 0, utils.ErrNotProperty
	}

	// 地产价值 = 原价 + (升级等级 * 升级费用)
	upgradeValue := t.Level * (t.Price / 2)
	return t.Price + upgradeValue, nil
}

// IsSpecialTile 检查是否为特殊地块
func (t *Tile) IsSpecialTile() bool {
	return t.Type == TileChance || t.Type == TileFate ||
		t.Type == TileBank || t.Type == TilePrison
}

// Clone 创建地块的深拷贝
func (t *Tile) Clone() *Tile {
	rentPriceCopy := make([]int, len(t.RentPrice))
	copy(rentPriceCopy, t.RentPrice)

	return &Tile{
		ID:        t.ID,
		Name:      t.Name,
		Type:      t.Type,
		Price:     t.Price,
		OwnerID:   t.OwnerID,
		Level:     t.Level,
		RentPrice: rentPriceCopy,
	}
}

// Reset 重置地块状态
func (t *Tile) Reset() {
	t.OwnerID = ""
	t.Level = 0
}
