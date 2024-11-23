// internal/user/user.go
package user

import (
	"monopoly/pkg/utils"
	"sync"
	"time"
)

// User 用户基本信息
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Coins    int       `json:"coins"`
	CreateAt time.Time `json:"createAt"`
}

// Transaction 交易记录
type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"` // "add" or "deduct"
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// UserGame 用户游戏记录
type UserGame struct {
	GameID    string    `json:"gameId"`
	UserID    string    `json:"userId"`
	JoinTime  time.Time `json:"joinTime"`
	EndTime   time.Time `json:"endTime,omitempty"`
	FinalRank int       `json:"finalRank,omitempty"`
}

// Manager 用户管理器
type Manager struct {
	users        map[string]*User
	transactions map[string][]Transaction
	userGames    map[string][]UserGame
	mutex        sync.RWMutex
}

// NewManager 创建新的用户管理器
func NewManager() *Manager {
	return &Manager{
		users:        make(map[string]*User),
		transactions: make(map[string][]Transaction),
		userGames:    make(map[string][]UserGame),
	}
}

// CreateUser 创建新用户
func (m *Manager) CreateUser(user *User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[user.ID]; exists {
		return utils.ErrUserExists
	}

	user.CreateAt = time.Now()
	m.users[user.ID] = user
	m.transactions[user.ID] = make([]Transaction, 0)
	m.userGames[user.ID] = make([]UserGame, 0)

	return nil
}

// GetUser 获取用户信息
func (m *Manager) GetUser(id string) (*User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return nil, utils.ErrUserNotFound
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (m *Manager) UpdateUser(user *User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[user.ID]; !exists {
		return utils.ErrUserNotFound
	}

	m.users[user.ID] = user
	return nil
}

// DeleteUser 删除用户
func (m *Manager) DeleteUser(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[id]; !exists {
		return utils.ErrUserNotFound
	}

	delete(m.users, id)
	delete(m.transactions, id)
	delete(m.userGames, id)

	return nil
}

// AddCoins 添加游戏币
func (m *Manager) AddCoins(userID string, amount int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, exists := m.users[userID]
	if !exists {
		return utils.ErrUserNotFound
	}

	user.Coins += amount

	// 记录交易
	transaction := Transaction{
		ID:        utils.GenerateID(), // 假设utils包中有这个函数
		UserID:    userID,
		Type:      "add",
		Amount:    amount,
		Timestamp: time.Now(),
	}
	m.transactions[userID] = append(m.transactions[userID], transaction)

	return nil
}

// DeductCoins 扣除游戏币
func (m *Manager) DeductCoins(userID string, amount int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user, exists := m.users[userID]
	if !exists {
		return utils.ErrUserNotFound
	}

	if user.Coins < amount {
		return utils.ErrInsufficientFunds
	}

	user.Coins -= amount

	// 记录交易
	transaction := Transaction{
		ID:        utils.GenerateID(),
		UserID:    userID,
		Type:      "deduct",
		Amount:    amount,
		Timestamp: time.Now(),
	}
	m.transactions[userID] = append(m.transactions[userID], transaction)

	return nil
}

// GetTransactions 获取用户交易记录
func (m *Manager) GetTransactions(userID string) []Transaction {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if transactions, exists := m.transactions[userID]; exists {
		return transactions
	}
	return []Transaction{}
}

// RecordGameParticipation 记录用户参与游戏
func (m *Manager) RecordGameParticipation(userID string, gameID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[userID]; !exists {
		return utils.ErrUserNotFound
	}

	userGame := UserGame{
		GameID:   gameID,
		UserID:   userID,
		JoinTime: time.Now(),
	}
	m.userGames[userID] = append(m.userGames[userID], userGame)

	return nil
}

// GetUserGames 获取用户参与的游戏记录
func (m *Manager) GetUserGames(userID string) []UserGame {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if games, exists := m.userGames[userID]; exists {
		return games
	}
	return []UserGame{}
}

// UpdateGameResult 更新用户游戏结果
func (m *Manager) UpdateGameResult(userID string, gameID string, rank int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	games := m.userGames[userID]
	for i := range games {
		if games[i].GameID == gameID {
			games[i].EndTime = time.Now()
			games[i].FinalRank = rank
			return nil
		}
	}

	return utils.ErrGameNotFound
}
