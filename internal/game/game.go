// internal/game/game.go
package game

import (
	"monopoly/pkg/utils"
	"sort"
	"sync"
	"time"
)

// Game 表示一局游戏
type Game struct {
	ID                 string             `json:"id"`
	Players            map[string]*Player `json:"players"`
	Status             GameStatus         `json:"status"`
	PrizePool          int                `json:"prizePool"`
	Map                *GameMap           `json:"map"`
	CurrentPlayerID    string             `json:"currentPlayerId"`
	CurrentTurnStarted time.Time          `json:"currentTurnStarted"`
	StartTime          time.Time          `json:"startTime"`
	Actions            []*GameAction      `json:"actions"`
	mutex              sync.RWMutex
}

// PlayerResult 表示玩家的最终游戏结果
type PlayerResult struct {
	PlayerID      string `json:"playerId"`
	Name          string `json:"name"`
	FinalCoins    int    `json:"finalCoins"`
	PropertyValue int    `json:"propertyValue"`
	PropertyCount int    `json:"propertyCount"`
	TotalAssets   int    `json:"totalAssets"`
}

// NewGame 创建新游戏
func NewGame(id string) *Game {
	return &Game{
		ID:      id,
		Players: make(map[string]*Player),
		Status:  StatusWaiting,
		Map:     NewDefaultMap(),
		Actions: make([]*GameAction, 0),
	}
}

// AddPlayer 添加玩家到游戏
func (g *Game) AddPlayer(player *Player) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if len(g.Players) >= MaxPlayers {
		return utils.ErrGameFull
	}

	if g.Status != StatusWaiting {
		return utils.ErrGameInProgress
	}

	if _, exists := g.Players[player.ID]; exists {
		return utils.ErrPlayerExists
	}

	g.Players[player.ID] = player
	return nil
}

// StartGame 开始游戏
func (g *Game) StartGame() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if err := g.canStartGame(); err != nil {
		return err
	}

	if err := g.collectEntranceFees(); err != nil {
		return err
	}

	g.initializeGame()
	return nil
}

// 游戏开始前的检查
func (g *Game) canStartGame() error {
	if len(g.Players) < MinStartPlayers {
		return utils.ErrNotEnoughPlayers
	}

	if g.Status != StatusWaiting {
		return utils.ErrInvalidGameState
	}

	return nil
}

// 收集入场费
func (g *Game) collectEntranceFees() error {
	for _, player := range g.Players {
		if player.Coins < InitialEntranceFee {
			return utils.ErrInsufficientFunds
		}
		player.Coins -= InitialEntranceFee
		g.PrizePool += InitialEntranceFee
	}
	return nil
}

// 初始化游戏状态
func (g *Game) initializeGame() {
	playerIDs := make([]string, 0, len(g.Players))
	for id := range g.Players {
		playerIDs = append(playerIDs, id)
	}
	sort.Strings(playerIDs)
	g.CurrentPlayerID = playerIDs[0]

	g.Status = StatusPlaying
	g.StartTime = time.Now()
	g.CurrentTurnStarted = time.Now()
}

// GetRemainingTime 获取游戏剩余时间（秒）
func (g *Game) GetRemainingTime() int {
	if g.Status != StatusPlaying {
		return 0
	}
	remaining := GameTimeout - time.Since(g.StartTime)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// GetTurnTimeLeft 获取当前回合剩余时间（秒）
func (g *Game) GetTurnTimeLeft() int {
	if g.Status != StatusPlaying {
		return 0
	}
	remaining := TurnTimeout - time.Since(g.CurrentTurnStarted)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// AddAction 添加游戏动作记录
func (g *Game) AddAction(action *GameAction) {
	g.Actions = append(g.Actions, action)
}

// ... 保持之前的代码不变 ...

// EndGame 结束游戏并分配奖池
func (g *Game) EndGame() error {
	if g.Status != StatusPlaying {
		return utils.ErrInvalidGameState
	}

	g.Status = StatusFinished

	// 计算玩家总资产
	type PlayerAsset struct {
		PlayerID string
		Total    int
	}
	assets := make([]PlayerAsset, 0, len(g.Players))

	// 计算每个玩家的总资产（现金 + 地产价值）
	for id, player := range g.Players {
		total := player.Coins
		// 加上地产价值
		for _, tile := range g.Map.Tiles {
			if tile.OwnerID == id {
				total += tile.Price * (tile.Level + 1)
			}
		}
		assets = append(assets, PlayerAsset{id, total})
	}

	// 按总资产排序
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Total > assets[j].Total
	})

	// 分配奖池
	prizeRatios := []float64{0.5, 0.3, 0.15, 0.05} // 奖池分配比例
	for i, asset := range assets {
		if i >= len(prizeRatios) {
			break
		}
		prize := int(float64(g.PrizePool) * prizeRatios[i])
		g.Players[asset.PlayerID].Coins += prize
	}

	// 记录游戏结束动作
	g.Actions = append(g.Actions, &GameAction{
		Type:      "gameEnd",
		Timestamp: time.Now(),
	})

	return nil
}

// GetFinalResults 获取游戏最终结果
func (g *Game) GetFinalResults() ([]PlayerResult, error) {
	if g.Status != StatusFinished {
		return nil, utils.ErrInvalidGameState
	}

	results := make([]PlayerResult, 0, len(g.Players))
	for id, player := range g.Players {
		propertyValue := 0
		propertyCount := 0
		for _, tile := range g.Map.Tiles {
			if tile.OwnerID == id {
				propertyValue += tile.Price * (tile.Level + 1)
				propertyCount++
			}
		}

		results = append(results, PlayerResult{
			PlayerID:      id,
			Name:          player.Name,
			FinalCoins:    player.Coins,
			PropertyValue: propertyValue,
			PropertyCount: propertyCount,
			TotalAssets:   player.Coins + propertyValue,
		})
	}

	// 按总资产排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalAssets > results[j].TotalAssets
	})

	return results, nil
}
