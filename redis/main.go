package main

//redis 6使用v8，redis 7使用v9
import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func init() {
	//连接redis
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", //无密码
		DB:       0,  //默认DB
	})
}

func main() {
	//ctx:=context.TODO()也行，都是空的context
	ctx := context.Background()

	//0代表不过期
	// err := client.Set(ctx, "GoRedisString", "123456", 4*time.Minute).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.SetNX(ctx, "GoRedisStringNumber", 5, 4*time.Minute).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// vals, err := client.MGet(ctx, "GoRedisString", "GoRedisStringNumber").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key 不存在")
	// } else if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(vals)

	// valFloat, err := client.IncrByFloat(ctx, "GoRedisStringNumber", 2.2).Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(valFloat)

	//使用原生命令
	// result, err := client.Do(ctx, "getset", "gorediskey", "goredisval2").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result.(string))

	//这里0可不代表不过期，用0的话直接就过期了
	// err = client.Expire(ctx, "goredisnums", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// client.HSet(ctx, "GoRedisHash", "username1", "zhangsan", "username2", "lisi", "username3", "wangwu")
	// username, _ := client.HGetAll(ctx, "GoRedisHahs").Result()
	// fmt.Println(username)
	// for field, value := range client.HGetAll(ctx, "GoRedisHash").Val() {
	// 	fmt.Println(field, value)
	// }
	// count, _ := client.HIncrByFloat(ctx, "GoRedisHash", "wjh", 666).Result()
	// fmt.Println(count)
	// exist, _ := client.HExists(ctx, "GoRedisHash", "username4").Result()
	// fmt.Println(exist)

	//同理Hlen,Hkeys,HMGet,HMSet,HSetNX,HDel

	// values := []interface{}{1, 2, 3, 4, "中"}
	// client.LPush(ctx, "GoRedisList", values...) //...代表传入不定长参数
	// client.LPushX(ctx, "GoRedisList", 11)
	// val, _ := client.RPop(ctx, "GoRedisList").Result()
	// fmt.Println(val)
	// len, _ := client.LLen(ctx, "GoRedisList").Result()
	// fmt.Println(len)
	// vals, _ := client.LRange(ctx, "GoRedisList", 0, -1).Result()
	// fmt.Println(vals)
	// client.LInsert(ctx, "GoRedisList", "after", 2, -1)
	//同理Lrem,Ltrim,Lindex

	// client.SAdd(ctx, "GoRedisSet1", 100, 100, 200, 300, 400, 10000)
	// exist, _ := client.SIsMember(ctx, "GoRedisSet1", 100).Result()
	// fmt.Println(exist)
	// card, _ := client.SCard(ctx, "GoRedisSet1").Result()
	// fmt.Println(card)
	// client.SRem(ctx, "GoRedisSet1", 200)
	// members, _ := client.SMembers(ctx, "GoRedisSet1").Result()
	// fmt.Println(members)
	// client.SAdd(ctx, "GoRedisSet2", 400, 1000, 600, 700)
	// members, _ = client.SMembers(ctx, "GoRedisSet2").Result()
	// fmt.Println(members)
	// members, _ = client.SInter(ctx, "GoRedisSet1", "GoRedisSet2").Result()
	// fmt.Println(members)
	//同理spop,spopN

	client.ZAdd(ctx, "GoRedisZset", redis.Z{Score: 2.5, Member: "zhangsan"}, redis.Z{Score: 3.3, Member: "lisi"}, redis.Z{Score: 6, Member: "wangwu"})
	count, _ := client.ZCount(ctx, "GoRedisZset", "1", "5").Result()
	fmt.Println(count)
	client.ZIncrBy(ctx, "GoRedisZset", 2, "zhangsan")
	result, _ := client.ZRevRange(ctx, "GoRedisZset", 0, -1).Result()
	fmt.Println(result)
	op := redis.ZRangeBy{Min: "2", Max: "10", Offset: 1, Count: 3}
	result2 := client.ZRangeByScoreWithScores(ctx, "GoRedisZset", &op)
	fmt.Println(result2.Val())
	//同理zcard,zrangebylex,zrem,zremrangebyscore,zremrangebyrank,
	//   zremrangebylex,zrank,zrevrank
	defer client.Close()

}
