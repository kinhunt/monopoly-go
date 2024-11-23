// internal/game/player.go
package game

import (
	"time"
)

// Player 表示游戏中的玩家
type Player struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Coins      int          `json:"coins"`
	Position   int          `json:"position"`
	Status     PlayerStatus `json:"status"`
	HasRolled  bool         `json:"hasRolled"`  // 是否已经掷过骰子
	InPrison   bool         `json:"inPrison"`   // 是否在监狱中
	PrisonDays int          `json:"prisonDays"` // 剩余监禁天数
	JoinTime   time.Time    `json:"joinTime"`   // 加入游戏的时间
}

// NewPlayer 创建新玩家
func NewPlayer(id string, name string, coins int) *Player {
	return &Player{
		ID:         id,
		Name:       name,
		Coins:      coins,
		Position:   0,
		Status:     PlayerStatusWaiting,
		HasRolled:  false,
		InPrison:   false,
		PrisonDays: 0,
		JoinTime:   time.Now(),
	}
}

// Reset 重置玩家状态
func (p *Player) Reset() {
	p.Position = 0
	p.Status = PlayerStatusWaiting
	p.HasRolled = false
	p.InPrison = false
	p.PrisonDays = 0
}

// CanMove 检查玩家是否可以移动
func (p *Player) CanMove() bool {
	return !p.HasRolled && !p.InPrison && p.Status == PlayerStatusPlaying
}

// EnterPrison 进入监狱
func (p *Player) EnterPrison() {
	p.InPrison = true
	p.PrisonDays = 1
}

// ExitPrison 离开监狱
func (p *Player) ExitPrison() {
	p.InPrison = false
	p.PrisonDays = 0
}

// UpdateStatus 更新玩家状态
func (p *Player) UpdateStatus(status PlayerStatus) {
	p.Status = status
}

// AddCoins 增加金币
func (p *Player) AddCoins(amount int) {
	p.Coins += amount
}

// DeductCoins 扣除金币
func (p *Player) DeductCoins(amount int) bool {
	if p.Coins < amount {
		return false
	}
	p.Coins -= amount
	return true
}

// MoveTo 移动到指定位置
func (p *Player) MoveTo(position int) {
	p.Position = position
	p.HasRolled = true
}

// EndTurn 结束回合
func (p *Player) EndTurn() {
	p.HasRolled = false
}

// GetStatus 获取玩家状态视图
func (p *Player) GetStatus() PlayerStatusView {
	return PlayerStatusView{
		ID:         p.ID,
		Name:       p.Name,
		Coins:      p.Coins,
		Position:   p.Position,
		Status:     p.Status,
		InPrison:   p.InPrison,
		PrisonDays: p.PrisonDays,
		HasRolled:  p.HasRolled,
	}
}

// PlayerStatusView 玩家状态视图
type PlayerStatusView struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Coins      int          `json:"coins"`
	Position   int          `json:"position"`
	Status     PlayerStatus `json:"status"`
	InPrison   bool         `json:"inPrison"`
	PrisonDays int          `json:"prisonDays"`
	HasRolled  bool         `json:"hasRolled"`
}

// Clone 创建玩家的深拷贝
func (p *Player) Clone() *Player {
	return &Player{
		ID:         p.ID,
		Name:       p.Name,
		Coins:      p.Coins,
		Position:   p.Position,
		Status:     p.Status,
		HasRolled:  p.HasRolled,
		InPrison:   p.InPrison,
		PrisonDays: p.PrisonDays,
		JoinTime:   p.JoinTime,
	}
}

// IsBankrupt 检查玩家是否破产
func (p *Player) IsBankrupt() bool {
	return p.Coins <= 0
}

// IsActive 检查玩家是否处于活跃状态
func (p *Player) IsActive() bool {
	return p.Status == PlayerStatusPlaying && !p.IsBankrupt()
}
