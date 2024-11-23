// internal/game/types.go
package game

import "time"

// 游戏相关常量
const (
	MaxPlayers      = 4
	MinStartPlayers = 2
	TurnTimeout     = 30 * time.Second
	GameTimeout     = 10 * time.Minute
)

// GameStatus 游戏状态
type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

// PlayerStatus 玩家状态
type PlayerStatus string

const (
	PlayerStatusWaiting PlayerStatus = "waiting"
	PlayerStatusPlaying PlayerStatus = "playing"
	PlayerStatusOffline PlayerStatus = "offline"
)

// ActionType 动作类型
type ActionType string

const (
	ActionRollDice    ActionType = "rollDice"
	ActionBuyProperty ActionType = "buyProperty"
	ActionPayRent     ActionType = "payRent"
	ActionUpgrade     ActionType = "upgrade"
)

// TileType 地块类型
type TileType string

const (
	TileStart    TileType = "start"
	TileProperty TileType = "property"
	TileChance   TileType = "chance"
	TileFate     TileType = "fate"
	TileBank     TileType = "bank"
	TilePrison   TileType = "prison"
)

// 游戏配置常量
const (
	InitialEntranceFee  = 1000
	PassingGoRewardRate = 0.02 // 过路奖励为奖池的2%
	ChanceRewardRate    = 0.01 // 机会奖励为奖池的1%
)
