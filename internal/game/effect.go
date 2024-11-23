// internal/game/effects.go
package game

import (
	"math/rand"
	"monopoly/pkg/utils"
	"time"
)

// handleTileEffect 处理地块效果
func (g *Game) handleTileEffect(player *Player) error {
	tile := g.Map.Tiles[player.Position]
	switch tile.Type {
	case TileProperty:
		if tile.OwnerID != "" && tile.OwnerID != player.ID {
			return g.handleRentPayment(player, tile)
		}
	case TileChance:
		return g.handleChanceCard(player)
	case TileFate:
		return g.handleFateCard(player)
	case TilePrison:
		return g.handlePrison(player)
	}
	return nil
}

// handleRentPayment 处理租金支付
func (g *Game) handleRentPayment(player *Player, tile *Tile) error {
	owner := g.Players[tile.OwnerID]
	if owner == nil {
		return utils.ErrPlayerNotFound
	}

	rent := tile.RentPrice[tile.Level]
	if player.Coins < rent {
		return utils.ErrInsufficientFunds
	}

	player.Coins -= rent
	owner.Coins += rent

	g.Actions = append(g.Actions, &GameAction{
		Type:      ActionPayRent,
		PlayerID:  player.ID,
		Position:  player.Position,
		Amount:    rent,
		Timestamp: time.Now(),
	})

	return nil
}

// handleChanceCard 处理机会卡片
func (g *Game) handleChanceCard(player *Player) error {
	effect := rand.Intn(100)
	var amount int

	switch {
	case effect < 30: // 30%概率获得奖池的1%
		amount = int(float64(g.PrizePool) * ChanceRewardRate)
		player.Coins += amount
		g.PrizePool -= amount

		g.Actions = append(g.Actions, &GameAction{
			Type:      "chanceReward",
			PlayerID:  player.ID,
			Amount:    amount,
			Timestamp: time.Now(),
		})

	case effect < 60: // 30%概率失去100金币
		amount = 100
		if player.Coins < amount {
			amount = player.Coins
		}
		player.Coins -= amount
		g.PrizePool += amount

		g.Actions = append(g.Actions, &GameAction{
			Type:      "chancePenalty",
			PlayerID:  player.ID,
			Amount:    -amount,
			Timestamp: time.Now(),
		})

	case effect < 80: // 20%概率传送到随机位置
		oldPosition := player.Position
		player.Position = rand.Intn(len(g.Map.Tiles))

		g.Actions = append(g.Actions, &GameAction{
			Type:      "chanceTeleport",
			PlayerID:  player.ID,
			Position:  player.Position,
			Timestamp: time.Now(),
		})

		// 如果经过起点，给予奖励
		if player.Position < oldPosition {
			g.handlePassingGo(player)
		}
	}

	return nil
}

// handleFateCard 处理命运卡片
func (g *Game) handleFateCard(player *Player) error {
	effect := rand.Intn(100)
	var amount int

	switch {
	case effect < 30: // 30%概率获得其他玩家每人100金币
		amount = 0
		for _, p := range g.Players {
			if p.ID != player.ID {
				deduct := 100
				if p.Coins < deduct {
					deduct = p.Coins
				}
				p.Coins -= deduct
				amount += deduct
			}
		}
		player.Coins += amount

		g.Actions = append(g.Actions, &GameAction{
			Type:      "fateCollect",
			PlayerID:  player.ID,
			Amount:    amount,
			Timestamp: time.Now(),
		})

	case effect < 60: // 30%概率支付所有地产10%的维护费
		amount = 0
		for _, tile := range g.Map.Tiles {
			if tile.OwnerID == player.ID {
				fee := tile.Price / 10
				amount += fee
			}
		}
		if amount > player.Coins {
			amount = player.Coins
		}
		player.Coins -= amount
		g.PrizePool += amount

		g.Actions = append(g.Actions, &GameAction{
			Type:      "fateMaintenance",
			PlayerID:  player.ID,
			Amount:    -amount,
			Timestamp: time.Now(),
		})

	case effect < 80: // 20%概率获得一个随机空地产
		for _, tile := range g.Map.Tiles {
			if tile.Type == TileProperty && tile.OwnerID == "" {
				tile.OwnerID = player.ID

				g.Actions = append(g.Actions, &GameAction{
					Type:      "fateProperty",
					PlayerID:  player.ID,
					Position:  tile.ID,
					Timestamp: time.Now(),
				})

				break
			}
		}
	}

	return nil
}

// handlePrison 处理监狱
func (g *Game) handlePrison(player *Player) error {
	player.InPrison = true
	player.PrisonDays = 1

	g.Actions = append(g.Actions, &GameAction{
		Type:      "prison",
		PlayerID:  player.ID,
		Timestamp: time.Now(),
	})

	return nil
}
