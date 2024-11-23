// internal/manager/manager.go
package manager

import (
	"monopoly/internal/game"
	"monopoly/pkg/utils"
	"sync"
)

type GameManager struct {
	games map[string]*game.Game
	mutex sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*game.Game),
	}
}

func (gm *GameManager) CreateGame(id string) *game.Game {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	newGame := game.NewGame(id)
	gm.games[id] = newGame
	return newGame
}

func (gm *GameManager) GetGame(id string) (*game.Game, error) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()

	game, exists := gm.games[id]
	if !exists {
		return nil, utils.ErrNotFound
	}

	return game, nil
}
