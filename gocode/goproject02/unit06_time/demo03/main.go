package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// =========================================================
// 自定义 JSON 时间类型：格式 "2006-01-02 15:04:05"
// =========================================================
// type CustomTime struct {
// 	time.Time
// }

type CustomTime struct {
	time.Time
}

const customLayout = "2006-01-02 15:04:05"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format(customLayout) + "/自定义方法" + `"`), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	s = strings.TrimRight(s, `/`)

	if s == "null" || s == "" {
		return nil
	}
	t, err := time.ParseInLocation(customLayout, s, time.Local)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// =========================================================
// 自定义 JSON 时间类型：Unix 时间戳（秒）
// =========================================================
type TimestampTime struct {
	time.Time
}

func (t TimestampTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *TimestampTime) UnmarshalJSON(data []byte) error {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(ts, 0)
	return nil
}

// -------------------------------------------------------
// 1. JSON 默认格式（RFC 3339）
// -------------------------------------------------------
func demoJSONDefault() {
	fmt.Println("\n===== 1. JSON 默认格式（RFC 3339）=====")

	type Event struct {
		Name    string    `json:"name"`
		StartAt time.Time `json:"start_at"`
	}

	// 序列化
	e := Event{Name: "技术分享会", StartAt: time.Now()}
	data, _ := json.Marshal(e)
	fmt.Println("序列化:", string(data))

	// 反序列化
	jsonStr := `{"name":"团建","start_at":"2026-03-31T14:00:00+08:00"}`
	var e2 Event
	json.Unmarshal([]byte(jsonStr), &e2)
	fmt.Printf("反序列化: name=%s, time=%v\n", e2.Name, e2.StartAt.Format(customLayout))
}

// -------------------------------------------------------
// 2. JSON 自定义格式（"2006-01-02 15:04:05"）
// -------------------------------------------------------
func demoJSONCustomFormat() {
	fmt.Println("\n===== 2. JSON 自定义格式 =====")

	type Order struct {
		OrderNo   string     `json:"order_no"`
		CreatedAt CustomTime `json:"created_at"`
	}

	// 序列化
	o := Order{
		OrderNo:   "ORD-20260331-001",
		CreatedAt: CustomTime{time.Now()},
	}
	data, _ := json.Marshal(o)
	fmt.Println("序列化:", string(data))

	// 反序列化
	jsonStr := `{"order_no":"ORD-20260331-002","created_at":"2026-03-31 09:30:00"}`
	var o2 Order
	json.Unmarshal([]byte(jsonStr), &o2)
	fmt.Printf("反序列化: order_no=%s, time=%v\n", o2.OrderNo, o2.CreatedAt.Format(customLayout))
}

// -------------------------------------------------------
// 3. JSON 时间戳格式
// -------------------------------------------------------
func demoJSONTimestamp() {
	fmt.Println("\n===== 3. JSON 时间戳格式 =====")

	type Session struct {
		UserID    int           `json:"user_id"`
		LoginAt   TimestampTime `json:"login_at"`
		ExpiredAt TimestampTime `json:"expired_at"`
	}

	now := time.Now()
	s := Session{
		UserID:    1001,
		LoginAt:   TimestampTime{now},
		ExpiredAt: TimestampTime{now.Add(24 * time.Hour)},
	}
	data, _ := json.Marshal(s)
	fmt.Println("序列化:", string(data))

	jsonStr := `{"user_id":1002,"login_at":1743400200,"expired_at":1743486600}`
	var s2 Session
	json.Unmarshal([]byte(jsonStr), &s2)
	fmt.Printf("反序列化: user_id=%d, login=%s, expired=%s\n",
		s2.UserID,
		s2.LoginAt.Format(customLayout),
		s2.ExpiredAt.Format(customLayout),
	)
}

// -------------------------------------------------------
// 4. MySQL: time.Time ↔ DATETIME/TIMESTAMP（模拟）
// -------------------------------------------------------
func demoMySQLTime() {
	fmt.Println("\n===== 4. MySQL 时间转换（模拟）=====")

	// 模拟 MySQL DATETIME 字段读取
	// 实际项目中：dsn 加 parseTime=true&loc=Local，driver 会自动转换
	mysqlDatatime := "2006-01-02 14:40:00"
	t, _ := time.ParseInLocation(customLayout, mysqlDatatime, time.Local)
	fmt.Printf("MySQL DATETIME → time.Time: %v\n", t)

	// time.Time → MySQL DATETIME 格式字符串
	now := time.Now()
	formatted := now.Format("2006-01-02 15:04:05")
	fmt.Printf("time.Time → MySQL DATETIME: %s\n", formatted)

	// MySQL DATE 类型
	mysqlDate := "2026-03-31"
	t2, _ := time.ParseInLocation("2006-01-02", mysqlDate, time.Local)
	fmt.Printf("MySQL DATE → time.Time: %v\n", t2)

	// time.Time → MySQL DATE
	fmt.Printf("time.Time → MySQL DATE: %s\n", now.Format("2006-01-02"))
}

// -------------------------------------------------------
// 5. MySQL: INT 时间戳字段
// -------------------------------------------------------
func demoMySQLTimestamp() {
	fmt.Println("\n===== 5. MySQL 时间戳字段 =====")

	// 模拟从 MySQL INT 字段读出的时间戳
	var dbTimestamp int64 = 1743400200

	// INT → time.Time
	t := time.Unix(dbTimestamp, 0)
	fmt.Printf("MySQL INT(%d) → time.Time: %s\n", dbTimestamp, t.Format(customLayout))

	// time.Time → INT（写入 MySQL）
	now := time.Now()
	ts := now.Unix()
	fmt.Printf("time.Time → MySQL INT: %d\n", ts)

	// 毫秒时间戳
	tsMilli := now.UnixMilli()
	fmt.Printf("time.Time → MySQL BIGINT(毫秒): %d\n", tsMilli)

	// 毫秒时间戳 → time.Time
	t2 := time.UnixMilli(tsMilli)
	fmt.Printf("MySQL BIGINT(毫秒) → time.Time: %s\n", t2.Format(customLayout))
}

// -------------------------------------------------------
// 6. MySQL: sql.NullTime 处理 NULL 值
// -------------------------------------------------------
func demoMySQLNullTime() {
	fmt.Println("\n===== 6. MySQL NullTime 处理 NULL =====")

	// 模拟有值的情况
	validTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if validTime.Valid {
		fmt.Printf("有值: %s\n", validTime.Time.Format(customLayout))
	}

	// 模拟 NULL 的情况
	nullTime := sql.NullTime{
		Valid: false,
	}
	if nullTime.Valid {
		fmt.Printf("有值: %s\n", nullTime.Time.Format(customLayout))
	} else {
		fmt.Println("NULL: 该字段为空（如 deleted_at 未删除）")
	}

	// 使用指针方式处理（GORM 常用）
	var deletedAt *time.Time = nil
	if deletedAt == nil {
		fmt.Println("*time.Time 为 nil → MySQL NULL")
	}

	now := time.Now()
	deletedAt = &now
	fmt.Printf("*time.Time 有值 → MySQL DATETIME: %s\n", deletedAt.Format(customLayout))
}

// -------------------------------------------------------
// 7. Redis: 时间戳存取
// -------------------------------------------------------
func demoRedisTimestamp() {
	fmt.Println("\n===== 7. Redis 时间戳存取（模拟）=====")

	now := time.Now()

	// 存储：time.Time → string（时间戳）
	// 对应 rdb.Set(ctx, "user:1:login_at", now.Unix(), 0)
	redisValue := strconv.FormatInt(now.Unix(), 10)
	fmt.Printf("存入 Redis: key=user:1:login_at, value=%s\n", redisValue)

	// 读取：string → time.Time
	// 对应 ts, _ := rdb.Get(ctx, "user:1:login_at").Int64()
	ts, _ := strconv.ParseInt(redisValue, 10, 64)
	t := time.Unix(ts, 0)
	fmt.Printf("从 Redis 读出: %s\n", t.Format(customLayout))

	// 存储格式化字符串
	redisStr := now.Format(customLayout)
	fmt.Printf("存入 Redis（字符串）: %s\n", redisStr)

	t2, _ := time.ParseInLocation(customLayout, redisStr, time.Local)
	fmt.Printf("从 Redis 读出（字符串）: %s\n", t2.Format(customLayout))
}

// -------------------------------------------------------
// 8. Redis: Duration 与 TTL
// -------------------------------------------------------
func demoRedisTTL() {
	fmt.Println("\n===== 8. Redis Duration 与 TTL（模拟）=====")

	// 设置过期时间
	// rdb.Set(ctx, "session:abc", data, 30*time.Minute)
	ttl := 30 * time.Minute
	fmt.Printf("SET session:abc ... EX %v\n", ttl)
	fmt.Printf("  等价秒数: %d\n", int(ttl.Seconds()))

	// ExpireAt：设定到某一精确时间点
	expireAt := time.Now().Add(48 * time.Hour)
	fmt.Printf("EXPIREAT key %d (到 %s 过期)\n", expireAt.Unix(), expireAt.Format(customLayout))

	// 读取 TTL 后的处理
	// ttlResult, _ := rdb.TTL(ctx, "session:abc").Result()
	ttlResult := 25*time.Minute + 30*time.Second // 模拟返回值
	fmt.Printf("TTL 返回: %v\n", ttlResult)
	fmt.Printf("  剩余分钟: %.1f\n", ttlResult.Minutes())
	fmt.Printf("  到期时间: %s\n", time.Now().Add(ttlResult).Format(customLayout))
}

// -------------------------------------------------------
// 9. Redis: 有序集合 Score 存时间
// -------------------------------------------------------
func demoRedisZSet() {
	fmt.Println("\n===== 9. Redis 有序集合 Score（模拟）=====")

	type ZMember struct {
		Score  float64
		Member string
	}

	now := time.Now()

	// 写入：用时间戳作为 score
	// rdb.ZAdd(ctx, "delay_queue", redis.Z{Score: float64(t.Unix()), Member: "task_1"})
	tasks := []ZMember{
		{Score: float64(now.Add(1 * time.Minute).Unix()), Member: "task_send_email"},
		{Score: float64(now.Add(5 * time.Minute).Unix()), Member: "task_generate_report"},
		{Score: float64(now.Add(30 * time.Minute).Unix()), Member: "task_cleanup"},
	}

	fmt.Println("写入延迟队列:")
	for _, t := range tasks {
		execAt := time.Unix(int64(t.Score), 0)
		fmt.Printf("  ZADD delay_queue %.0f %s (执行时间: %s)\n",
			t.Score, t.Member, execAt.Format(customLayout))
	}

	// 读取：查询已到期的任务
	// ZRANGEBYSCORE delay_queue -inf {now_timestamp}
	queryTime := now.Add(3 * time.Minute)
	fmt.Printf("\n查询 %s 之前到期的任务:\n", queryTime.Format("15:04:05"))
	for _, t := range tasks {
		if t.Score <= float64(queryTime.Unix()) {
			fmt.Printf("  已到期: %s\n", t.Member)
		}
	}

	// Score → time.Time
	fmt.Println("\nScore 转回时间:")
	for _, t := range tasks {
		execTime := time.Unix(int64(t.Score), 0)
		remaining := execTime.Sub(now).Round(time.Second)
		fmt.Printf("  %s → %s (还剩 %v)\n", t.Member, execTime.Format("15:04:05"), remaining)
	}
}

// -------------------------------------------------------
// main
// -------------------------------------------------------
func main() {
	// JSON 相关
	demoJSONDefault()
	demoJSONCustomFormat()
	demoJSONTimestamp()

	// MySQL 相关
	demoMySQLTime()
	demoMySQLTimestamp()
	demoMySQLNullTime()

	// Redis 相关
	demoRedisTimestamp()
	demoRedisTTL()
	demoRedisZSet()

	fmt.Println("\n===== 全部示例执行完毕 =====")
}
