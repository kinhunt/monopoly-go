// internal/game/action.go
package game

import (
	"math/rand"
	"monopoly/pkg/utils"
	"time"
)

// GameAction 表示游戏动作
type GameAction struct {
	Type      ActionType `json:"type"`
	PlayerID  string     `json:"playerId"`
	Position  int        `json:"position,omitempty"`
	Amount    int        `json:"amount,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

// RollDice 掷骰子并移动玩家
func (g *Game) RollDice(playerID string) (*GameAction, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if err := g.validateGameState(playerID); err != nil {
		return nil, err
	}

	player := g.Players[playerID]
	if player.HasRolled {
		return nil, utils.ErrAlreadyRolled
	}

	// 处理监狱状态
	if err := g.handlePrisonState(player); err != nil {
		return nil, err
	}

	// 执行移动
	steps := rand.Intn(6) + 1
	oldPosition := player.Position
	newPosition := (oldPosition + steps) % len(g.Map.Tiles)
	player.Position = newPosition
	player.HasRolled = true

	action := &GameAction{
		Type:      ActionRollDice,
		PlayerID:  playerID,
		Position:  newPosition,
		Timestamp: time.Now(),
	}

	// 处理过起点奖励
	if newPosition < oldPosition {
		g.handlePassingGo(player)
	}

	// 处理新位置效果
	if err := g.handleTileEffect(player); err != nil {
		return action, err // 返回动作但同时返回错误
	}

	g.Actions = append(g.Actions, action)
	return action, nil
}

// BuyProperty 购买地产
func (g *Game) BuyProperty(playerID string) (*GameAction, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if err := g.validateGameState(playerID); err != nil {
		return nil, err
	}

	player := g.Players[playerID]
	tile := g.Map.Tiles[player.Position]

	if err := g.validatePropertyPurchase(player, tile); err != nil {
		return nil, err
	}

	// 执行购买
	player.Coins -= tile.Price
	tile.OwnerID = playerID
	g.PrizePool += tile.Price

	action := &GameAction{
		Type:      ActionBuyProperty,
		PlayerID:  playerID,
		Position:  player.Position,
		Amount:    tile.Price,
		Timestamp: time.Now(),
	}

	g.Actions = append(g.Actions, action)
	return action, nil
}

// UpgradeProperty 升级地产
func (g *Game) UpgradeProperty(playerID string, position int) (*GameAction, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if err := g.validateGameState(playerID); err != nil {
		return nil, err
	}

	player := g.Players[playerID]
	tile := g.Map.Tiles[position]

	if err := g.validatePropertyUpgrade(player, tile); err != nil {
		return nil, err
	}

	upgradeCost := tile.Price / 2
	player.Coins -= upgradeCost
	tile.Level++
	g.PrizePool += upgradeCost

	action := &GameAction{
		Type:      ActionUpgrade,
		PlayerID:  playerID,
		Position:  position,
		Amount:    upgradeCost,
		Timestamp: time.Now(),
	}

	g.Actions = append(g.Actions, action)
	return action, nil
}

// NextTurn 进入下一个回合
func (g *Game) NextTurn() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Status != StatusPlaying {
		return utils.ErrInvalidGameState
	}

	// 重置当前玩家状态
	if currentPlayer := g.Players[g.CurrentPlayerID]; currentPlayer != nil {
		currentPlayer.HasRolled = false
	}

	// 获取排序后的玩家ID列表
	playerIDs := g.getOrderedPlayerIDs()
	g.CurrentPlayerID = g.getNextPlayerID(playerIDs)
	g.CurrentTurnStarted = time.Now()

	// 检查游戏是否应该结束
	if time.Since(g.StartTime) >= GameTimeout {
		return g.EndGame()
	}

	return nil
}

// 以下是辅助方法

// validateGameState 验证游戏状态和玩家回合
func (g *Game) validateGameState(playerID string) error {
	if g.Status != StatusPlaying {
		return utils.ErrInvalidGameState
	}

	if g.CurrentPlayerID != playerID {
		return utils.ErrNotYourTurn
	}

	if _, exists := g.Players[playerID]; !exists {
		return utils.ErrPlayerNotFound
	}

	return nil
}

// validatePropertyPurchase 验证地产购买条件
func (g *Game) validatePropertyPurchase(player *Player, tile *Tile) error {
	if tile.Type != TileProperty {
		return utils.ErrNotProperty
	}

	if tile.OwnerID != "" {
		return utils.ErrPropertyOwned
	}

	if player.Coins < tile.Price {
		return utils.ErrInsufficientFunds
	}

	return nil
}

// validatePropertyUpgrade 验证地产升级条件
func (g *Game) validatePropertyUpgrade(player *Player, tile *Tile) error {
	if tile.Type != TileProperty {
		return utils.ErrNotProperty
	}

	if tile.OwnerID != player.ID {
		return utils.ErrNotOwner
	}

	if tile.Level >= len(tile.RentPrice)-1 {
		return utils.ErrMaxLevel
	}

	upgradeCost := tile.Price / 2
	if player.Coins < upgradeCost {
		return utils.ErrInsufficientFunds
	}

	return nil
}

// handlePrisonState 处理玩家的监狱状态
func (g *Game) handlePrisonState(player *Player) error {
	if player.InPrison {
		if player.PrisonDays > 0 {
			player.PrisonDays--
			player.HasRolled = true
			return utils.ErrInPrison
		}
		player.InPrison = false
		player.PrisonDays = 0
	}
	return nil
}

// handlePassingGo 处理经过起点奖励
func (g *Game) handlePassingGo(player *Player) {
	passingGoReward := int(float64(g.PrizePool) * PassingGoRewardRate)
	player.Coins += passingGoReward
	g.PrizePool -= passingGoReward
}

// getOrderedPlayerIDs 获取排序后的玩家ID列表
func (g *Game) getOrderedPlayerIDs() []string {
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	return playerIDs
}

// getNextPlayerID 获取下一个玩家ID
func (g *Game) getNextPlayerID(playerIDs []string) string {
	for i, id := range playerIDs {
		if id == g.CurrentPlayerID {
			if i == len(playerIDs)-1 {
				return playerIDs[0] // 回到第一个玩家
			}
			return playerIDs[i+1]
		}
	}
	return playerIDs[0] // 如果找不到当前玩家（不应该发生），返回第一个玩家
}
