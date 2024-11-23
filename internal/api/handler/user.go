// internal/api/handler/user.go
package handler

import (
	"encoding/json"
	"monopoly/internal/api/response"
	"monopoly/internal/user"
	"monopoly/pkg/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// UserHandler 用户相关的HTTP请求处理器
type UserHandler struct {
	userManager *user.Manager
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(um *user.Manager) *UserHandler {
	return &UserHandler{
		userManager: um,
	}
}

// Create 创建新用户
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Coins int    `json:"coins"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 验证输入
	if req.ID == "" || req.Name == "" {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 验证初始金币数量
	if req.Coins < 0 {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 创建新用户
	newUser := &user.User{
		ID:    req.ID,
		Name:  req.Name,
		Coins: req.Coins,
	}

	if err := h.userManager.CreateUser(newUser); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.Success(newUser))
}

// Get 获取用户信息
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	user, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(user))
}

// Update 更新用户信息
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	if req.Name == "" {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	user, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	user.Name = req.Name
	if err := h.userManager.UpdateUser(user); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(user))
}

// AddCoins 为用户添加游戏币
func (h *UserHandler) AddCoins(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 验证金币数量
	if req.Amount <= 0 {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 验证用户是否存在
	if _, err := h.userManager.GetUser(userID); err != nil {
		response.JsonError(w, err)
		return
	}

	// 添加金币
	if err := h.userManager.AddCoins(userID, req.Amount); err != nil {
		response.JsonError(w, err)
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(updatedUser))
}

// DeductCoins 扣除用户游戏币
func (h *UserHandler) DeductCoins(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var req struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	// 验证金币数量
	if req.Amount <= 0 {
		response.JsonError(w, utils.ErrInvalidInput)
		return
	}

	if err := h.userManager.DeductCoins(userID, req.Amount); err != nil {
		response.JsonError(w, err)
		return
	}

	// 重新获取更新后的用户信息
	updatedUser, _ := h.userManager.GetUser(userID)
	response.JSON(w, http.StatusOK, response.Success(updatedUser))
}

// Delete 删除用户
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if err := h.userManager.DeleteUser(userID); err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(nil))
}

// GetUserTransactions 获取用户交易记录
func (h *UserHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	user, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	transactions := h.userManager.GetTransactions(userID)
	response.JSON(w, http.StatusOK, response.Success(map[string]interface{}{
		"user":         user,
		"transactions": transactions,
	}))
}

// GetUserGames 获取用户参与的游戏列表
func (h *UserHandler) GetUserGames(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	user, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	games := h.userManager.GetUserGames(userID)
	response.JSON(w, http.StatusOK, response.Success(map[string]interface{}{
		"user":  user,
		"games": games,
	}))
}

// CheckUserBalance 检查用户余额
func (h *UserHandler) CheckUserBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	user, err := h.userManager.GetUser(userID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.Success(map[string]interface{}{
		"userId": user.ID,
		"coins":  user.Coins,
	}))
}
