// pkg/utils/errors.go
package utils

import (
	"errors"
	"net/http"
)

// 基础错误
var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// 游戏状态相关错误
var (
	ErrGameNotFound     = errors.New("game not found")
	ErrGameFull         = errors.New("game is full")
	ErrGameInProgress   = errors.New("game is already in progress")
	ErrGameFinished     = errors.New("game is already finished")
	ErrInvalidGameState = errors.New("invalid game state")
	ErrNotEnoughPlayers = errors.New("not enough players to start game")
)

// 玩家相关错误
var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrPlayerExists   = errors.New("player already exists")
	ErrNotYourTurn    = errors.New("not your turn")
	ErrAlreadyRolled  = errors.New("already rolled dice this turn")
	ErrInPrison       = errors.New("player is in prison")
)

// 地产相关错误
var (
	ErrInvalidPosition      = errors.New("invalid position")
	ErrNotProperty          = errors.New("tile is not a property")
	ErrPropertyOwned        = errors.New("property is already owned")
	ErrNotOwner             = errors.New("player is not the owner of this property")
	ErrMaxLevel             = errors.New("property is already at maximum level")
	ErrCannotAfford         = errors.New("cannot afford this action")
	ErrPropertyNotOwned     = errors.New("property not owned")
	ErrInvalidPropertyLevel = errors.New("invalid property level")
)

// 游戏操作相关错误
var (
	ErrActionNotAllowed = errors.New("action not allowed")
	ErrInvalidAction    = errors.New("invalid action")
	ErrTimeout          = errors.New("operation timeout")
)

// 用户相关错误
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrInvalidUsername = errors.New("invalid username")
)

// 错误检查函数
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrGameNotFound) ||
		errors.Is(err, ErrPlayerNotFound) ||
		errors.Is(err, ErrUserNotFound)
}

func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrInvalidPosition) ||
		errors.Is(err, ErrInvalidAction) ||
		errors.Is(err, ErrInvalidUserID) ||
		errors.Is(err, ErrInvalidUsername)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrNotYourTurn) ||
		errors.Is(err, ErrNotOwner)
}

func IsInsufficientFunds(err error) bool {
	return errors.Is(err, ErrInsufficientFunds) ||
		errors.Is(err, ErrCannotAfford)
}

func IsGameStateError(err error) bool {
	return errors.Is(err, ErrGameInProgress) ||
		errors.Is(err, ErrGameFinished) ||
		errors.Is(err, ErrInvalidGameState)
}

func IsPropertyError(err error) bool {
	return errors.Is(err, ErrNotProperty) ||
		errors.Is(err, ErrPropertyOwned) ||
		errors.Is(err, ErrMaxLevel)
}

// ErrorResponse 用于API响应的错误信息结构
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse 创建新的错误响应
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Code:    getErrorCode(err),
		Message: err.Error(),
	}
}

// getErrorCode 根据错误类型返回错误代码
func getErrorCode(err error) string {
	switch {
	case IsNotFound(err):
		return "NOT_FOUND"
	case IsInvalidInput(err):
		return "INVALID_INPUT"
	case IsUnauthorized(err):
		return "UNAUTHORIZED"
	case IsInsufficientFunds(err):
		return "INSUFFICIENT_FUNDS"
	case IsGameStateError(err):
		return "GAME_STATE_ERROR"
	case IsPropertyError(err):
		return "PROPERTY_ERROR"
	default:
		return "INTERNAL_ERROR"
	}
}

// HTTPStatusFromError 根据错误类型返回对应的HTTP状态码
func HTTPStatusFromError(err error) int {
	switch {
	case IsNotFound(err):
		return http.StatusNotFound
	case IsInvalidInput(err):
		return http.StatusBadRequest
	case IsUnauthorized(err):
		return http.StatusUnauthorized
	case IsInsufficientFunds(err):
		return http.StatusPaymentRequired
	case IsGameStateError(err):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
