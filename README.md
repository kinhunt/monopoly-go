# 大富翁游戏服务器设计文档

## 1. 项目概述

这是一个多人在线大富翁游戏的服务器实现，采用 Go 语言开发。游戏支持多个房间同时进行，每个房间支持2-4名玩家参与。游戏采用回合制，通过 HTTP API 提供服务。

### 1.1 核心特性
- 多房间并发支持
- 实时游戏状态管理
- 完整的游戏币经济系统
- 地产交易和升级系统
- 机会与命运卡片系统

## 2. 系统架构

### 2.1 项目结构
```
monopoly/
├── cmd/server/          # 服务器入口
├── internal/            # 内部包
│   ├── game/           # 游戏核心逻辑
│   ├── manager/        # 游戏管理器
│   ├── user/           # 用户管理
│   └── api/            # HTTP API 实现
└── pkg/                # 公共工具包
```

### 2.2 核心模块
1. **游戏核心(game)**
   - 游戏状态管理
   - 回合控制
   - 地图系统
   - 玩家动作处理

2. **用户系统(user)**
   - 用户账户管理
   - 游戏币管理
   - 交易记录

3. **游戏管理器(manager)**
   - 房间管理
   - 游戏生命周期控制
   - 并发控制

4. **API接口(api)**
   - REST API 实现
   - 请求处理
   - 响应封装

## 3. 游戏设计

### 3.1 游戏参数配置
```go
const (
    MaxPlayers     = 4            // 最大玩家数
    MinStartPlayers = 2           // 最小开始人数
    TurnTimeout    = 30 * time.Second  // 回合超时时间
    GameTimeout    = 10 * time.Minute  // 游戏超时时间
    InitialEntranceFee = 1000     // 入场费
)
```

### 3.2 游戏状态
```go
type GameStatus string

const (
    StatusWaiting  GameStatus = "waiting"   // 等待开始
    StatusPlaying  GameStatus = "playing"   // 游戏中
    StatusFinished GameStatus = "finished"  // 已结束
)
```

### 3.3 地图系统
- 环形地图，默认20个格子
- 地块类型：起点、地产、机会、命运、银行、监狱
- 地产可升级，有不同等级的过路费

### 3.4 奖池机制
- 所有玩家入场费进入奖池
- 过路费、地产交易费用进入奖池
- 游戏结束时按排名分配奖池

## 4. 核心流程

### 4.1 游戏创建流程
1. 创建游戏房间
2. 玩家加入
3. 收取入场费
4. 开始游戏

### 4.2 回合流程
1. 掷骰子移动
2. 触发格子效果
3. 执行玩家操作（购买/升级）
4. 结束回合

### 4.3 游戏结束条件
1. 达到时间限制
2. 某玩家资产达到目标
3. 某玩家破产

## 5. 数据模型

### 5.1 游戏对象
```go
type Game struct {
    ID                string
    Players           map[string]*Player
    Status            GameStatus
    PrizePool         int
    Map               *GameMap
    CurrentPlayerID   string
    CurrentTurnStarted time.Time
    StartTime         time.Time
    Actions           []*GameAction
}
```

### 5.2 玩家对象
```go
type Player struct {
    ID         string
    Name       string
    Coins      int
    Position   int
    Status     PlayerStatus
    HasRolled  bool
    InPrison   bool
    PrisonDays int
}
```

## 6. API 接口

### 6.1 基础端点
```
POST   /api/users              # 创建用户
POST   /api/games              # 创建游戏
POST   /api/games/{id}/join    # 加入游戏
GET    /api/games/{id}/status  # 获取状态
```

### 6.2 游戏操作端点
```
POST   /api/games/{id}/roll          # 掷骰子
POST   /api/games/{id}/property/buy  # 购买地产
POST   /api/games/{id}/end-turn      # 结束回合
```

## 7. 扩展建议

### 7.1 可扩展方向
1. **持久化存储**
   - 添加数据库支持
   - 实现用户数据持久化
   - 游戏记录存档

2. **实时通信**
   - 集成 WebSocket
   - 实现实时状态推送
   - 添加聊天功能

3. **游戏玩法**
   - 增加特殊卡片效果
   - 添加任务系统
   - 实现成就系统

4. **多人交互**
   - 添加交易系统
   - 实现联盟机制
   - 增加多人互动效果

### 7.2 性能优化
1. 使用缓存优化状态查询
2. 实现游戏状态快照
3. 优化并发处理机制

## 8. 开发指南

### 8.1 环境要求
- Go 1.16+
- 支持 HTTP/1.1
- gorilla/mux 路由器

### 8.2 本地开发
```bash
# 获取代码
git clone [repository-url]
cd monopoly

# 安装依赖
go mod tidy

# 运行服务器
go run cmd/server/main.go
```

### 8.3 测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/game
```

## 9. 注意事项

1. **并发安全**
   - 所有状态修改需要加锁
   - 使用互斥锁保护共享资源
   - 避免死锁情况

2. **错误处理**
   - 统一使用 utils 包中定义的错误
   - 适当的错误封装和转换
   - 完整的错误日志

3. **接口设计**
   - 遵循 RESTful 设计原则
   - 统一的响应格式
   - 合理的状态码使用

4. **代码规范**
   - 遵循 Go 标准代码规范
   - 添加必要的注释
   - 保持模块独立性


# 大富翁游戏 API 测试用例

## 1. 用户管理

### 1.1 创建用户
```bash
# 创建第一个用户
curl -X POST http://localhost:8080/api/users \
-H "Content-Type: application/json" \
-d '{
    "id": "user1",
    "name": "Player One",
    "coins": 5000
}'

# 创建第二个用户
curl -X POST http://localhost:8080/api/users \
-H "Content-Type: application/json" \
-d '{
    "id": "user2",
    "name": "Player Two",
    "coins": 5000
}'
```

### 1.2 查询用户
```bash
# 查询用户信息
curl http://localhost:8080/api/users/user1

# 查询用户余额
curl http://localhost:8080/api/users/user1/balance
```

### 1.3 添加游戏币
```bash
# 为用户添加游戏币
curl -X POST http://localhost:8080/api/users/user1/coins \
-H "Content-Type: application/json" \
-d '{
    "amount": 1000
}'
```

### 1.4 查询用户交易记录
```bash
# 查询用户交易历史
curl http://localhost:8080/api/users/user1/transactions
```

## 2. 游戏管理

### 2.1 创建游戏
```bash
# 创建新游戏
curl -X POST http://localhost:8080/api/games \
-H "Content-Type: application/json" \
-d '{
    "gameId": "game1"
}'
```

### 2.2 加入游戏
```bash
# 玩家一加入游戏
curl -X POST http://localhost:8080/api/games/game1/join \
-H "Content-Type: application/json" \
-d '{
    "userId": "user1"
}'

# 玩家二加入游戏
curl -X POST http://localhost:8080/api/games/game1/join \
-H "Content-Type: application/json" \
-d '{
    "userId": "user2"
}'
```

### 2.3 查询游戏状态
```bash
# 获取游戏信息
curl http://localhost:8080/api/games/game1

# 获取游戏状态
curl http://localhost:8080/api/games/game1/status
```

### 2.4 开始游戏
```bash
# 开始游戏
curl -X POST http://localhost:8080/api/games/game1/start \
-H "Content-Type: application/json"
```

## 3. 游戏操作

### 3.1 掷骰子
```bash
# 玩家一掷骰子
curl -X POST http://localhost:8080/api/games/game1/roll \
-H "Content-Type: application/json" \
-d '{
    "playerId": "user1"
}'
```

### 3.2 购买地产
```bash
# 玩家购买当前位置的地产
curl -X POST http://localhost:8080/api/games/game1/properties/1/buy \
-H "Content-Type: application/json" \
-d '{
    "playerId": "user1"
}'
```

### 3.3 升级地产
```bash
# 升级指定位置的地产
curl -X POST http://localhost:8080/api/games/game1/properties/1/upgrade \
-H "Content-Type: application/json" \
-d '{
    "playerId": "user1"
}'
```

### 3.4 结束回合
```bash
# 结束当前玩家的回合
curl -X POST http://localhost:8080/api/games/game1/end-turn \
-H "Content-Type: application/json" \
-d '{
    "playerId": "user1"
}'
```

### 3.5 查询玩家状态
```bash
# 获取玩家在游戏中的状态
curl http://localhost:8080/api/games/game1/players/user1
```

## 4. 完整游戏流程测试脚本
```bash
#!/bin/bash

# 1. 创建用户
echo "Creating users..."
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"id":"user1","name":"Player One","coins":5000}'
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"id":"user2","name":"Player Two","coins":5000}'

sleep 1

# 2. 创建游戏
echo "Creating game..."
curl -X POST http://localhost:8080/api/games -H "Content-Type: application/json" -d '{"gameId":"game1"}'

sleep 1

# 3. 玩家加入游戏
echo "Players joining game..."
curl -X POST http://localhost:8080/api/games/game1/join -H "Content-Type: application/json" -d '{"userId":"user1"}'
curl -X POST http://localhost:8080/api/games/game1/join -H "Content-Type: application/json" -d '{"userId":"user2"}'

sleep 1

# 4. 开始游戏
echo "Starting game..."
curl -X POST http://localhost:8080/api/games/game1/start -H "Content-Type: application/json"

sleep 1

# 5. 模拟几个回合
for i in {1..3}; do
    echo "Round $i"
    
    # 玩家1回合
    curl -X POST http://localhost:8080/api/games/game1/roll -H "Content-Type: application/json" -d '{"playerId":"user1"}'
    sleep 1
    curl -X POST http://localhost:8080/api/games/game1/properties/1/buy -H "Content-Type: application/json" -d '{"playerId":"user1"}'
    sleep 1
    curl -X POST http://localhost:8080/api/games/game1/end-turn -H "Content-Type: application/json" -d '{"playerId":"user1"}'
    sleep 1
    
    # 玩家2回合
    curl -X POST http://localhost:8080/api/games/game1/roll -H "Content-Type: application/json" -d '{"playerId":"user2"}'
    sleep 1
    curl -X POST http://localhost:8080/api/games/game1/properties/2/buy -H "Content-Type: application/json" -d '{"playerId":"user2"}'
    sleep 1
    curl -X POST http://localhost:8080/api/games/game1/end-turn -H "Content-Type: application/json" -d '{"playerId":"user2"}'
    sleep 1
done

# 6. 查看游戏状态
echo "Checking game status..."
curl http://localhost:8080/api/games/game1/status
```

## 5. 测试响应示例
```json
// 创建用户响应
{
    "success": true,
    "data": {
        "id": "user1",
        "name": "Player One",
        "coins": 5000,
        "createAt": "2024-11-22T10:00:00Z"
    }
}

// 游戏状态响应
{
    "success": true,
    "data": {
        "gameId": "game1",
        "status": "playing",
        "playerCount": 2,
        "currentPlayer": "user1",
        "remainingTime": 550,
        "turnTimeLeft": 25,
        "prizePool": 2000
    }
}
```

这些测试用例涵盖了所有主要功能，可以用来验证API的正确性。需要我解释任何具体的测试用例吗？
