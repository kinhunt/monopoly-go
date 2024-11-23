// pkg/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID 生成唯一ID
func GenerateID() string {
	// 使用时间戳作为前缀
	timestamp := time.Now().UnixNano()

	// 生成8字节的随机数
	b := make([]byte, 8)
	rand.Read(b)

	// 组合时间戳和随机数
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(b))
}

// GenerateGameID 生成游戏ID
func GenerateGameID() string {
	return fmt.Sprintf("game-%s", GenerateID())
}

// GenerateUserID 生成用户ID
func GenerateUserID() string {
	return fmt.Sprintf("user-%s", GenerateID())
}

// GenerateTransactionID 生成交易ID
func GenerateTransactionID() string {
	return fmt.Sprintf("tx-%s", GenerateID())
}
