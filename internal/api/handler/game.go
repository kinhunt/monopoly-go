// internal/api/handler/game.go
package handler

import (
	"encoding/json"
	"monopoly/internal/api/response"
	"monopoly/internal/game"
	"monopoly/internal/manager"
	"monopoly/internal/user"
	"monopoly/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GameHandler 游戏相关的HTTP请求处理器
type GameHandler struct {
	gameManager *manager.GameManager
	userManager *user.Manager
}

// NewGameHandler 创建新的游戏处理器
func NewGameHandler(gm *manager.GameManager, um *user.Manager) *GameHandler {
	return &GameHandler{
		gameManager: gm,
		userManager: um,
	}
}

// Create 创建新游戏
func (h *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameID string `json:"gameId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	newGame := h.gameManager.CreateGame(req.GameID)
	response.JSON(w, http.StatusCreated, response.Success(newGame))
}

// Get 获取游戏信息
func (h *GameHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	game, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(game))
}

// Join 加入游戏
func (h *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	user, err := h.userManager.GetUser(req.UserID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	player := game.NewPlayer(user.ID, user.Name, user.Coins)
	if err := g.AddPlayer(player); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(g))
}

// StartGame 开始游戏
func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	if err := g.StartGame(); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(g))
}

// RollDice 掷骰子
func (h *GameHandler) RollDice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	var req struct {
		PlayerID string `json:"playerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	action, err := g.RollDice(req.PlayerID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(action))
}

// BuyProperty 购买地产
func (h *GameHandler) BuyProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]
	// position在这里不需要，因为使用玩家当前位置
	// position, _ := strconv.Atoi(vars["position"])

	var req struct {
		PlayerID string `json:"playerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	action, err := g.BuyProperty(req.PlayerID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(action))
}

// UpgradeProperty 升级地产
func (h *GameHandler) UpgradeProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]
	position, _ := strconv.Atoi(vars["position"])

	var req struct {
		PlayerID string `json:"playerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	action, err := g.UpgradeProperty(req.PlayerID, position)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(action))
}

// EndTurn 结束回合
func (h *GameHandler) EndTurn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	var req struct {
		PlayerID string `json:"playerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	if err := g.NextTurn(); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(g))
}

// GetGameStatus 获取游戏状态
func (h *GameHandler) GetGameStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	// 构建游戏状态
	status := struct {
		GameID        string          `json:"gameId"`
		Status        game.GameStatus `json:"status"`
		PlayerCount   int             `json:"playerCount"`
		CurrentPlayer string          `json:"currentPlayer"`
		RemainingTime int             `json:"remainingTime"`
		TurnTimeLeft  int             `json:"turnTimeLeft"`
		PrizePool     int             `json:"prizePool"`
	}{
		GameID:        g.ID,
		Status:        g.Status,
		PlayerCount:   len(g.Players),
		CurrentPlayer: g.CurrentPlayerID,
		RemainingTime: g.GetRemainingTime(),
		TurnTimeLeft:  g.GetTurnTimeLeft(),
		PrizePool:     g.PrizePool,
	}

	response.JSON(w, http.StatusOK, response.Success(status))
}

// GetPlayerStatus 获取玩家状态
func (h *GameHandler) GetPlayerStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]
	playerID := vars["playerId"]

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	player, exists := g.Players[playerID]
	if !exists {
		response.JsonError(w, utils.ErrPlayerNotFound)
		return
	}

	status := player.GetStatus()
	response.JSON(w, http.StatusOK, response.Success(status))
}

// LeaveGame 离开游戏
func (h *GameHandler) LeaveGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]
	playerID := vars["playerId"]

	g, err := h.gameManager.GetGame(gameID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	// 只有在等待状态才能离开游戏
	if g.Status != game.StatusWaiting {
		response.JsonError(w, utils.ErrGameInProgress)
		return
	}

	delete(g.Players, playerID)
	response.JSON(w, http.StatusOK, response.Success(nil))
}
