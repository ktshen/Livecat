package redigogo

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var RedisPort = "8869"
var Password = "livestream"

func AllKey(RedisIP string, CHANNEL string) ([]string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	c.Do("AUTH", Password)
	allkeys, err := redis.Strings(c.Do("KEYS", "*"))
	return allkeys, err
}

func Get(RedisIP string, CHANNEL string, KEY_NAME string) (string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.String(c.Do("GET", KEY_NAME))

	return content, err
}

func Set(RedisIP string, CHANNEL string, KEY_NAME string, VAL string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("SET", KEY_NAME, VAL)

	return err
}

func Set_WithExpire(RedisIP string, CHANNEL string, KEY_NAME string, VAL string, TIME int) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("SET", KEY_NAME, VAL)
	CheckError(err)
	_, err = c.Do("EXPIRE", KEY_NAME, TIME)

	return err
}

func Del(RedisIP string, CHANNEL string, KEY_NAME string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("DEL", KEY_NAME)

	return err
}

//if not find, return is false
func Exists(RedisIP string, CHANNEL string, KEY_NAME string) (bool, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	EXIST, err := redis.Bool(c.Do("EXISTS", KEY_NAME))

	return EXIST, err
}

func Lpush(RedisIP string, CHANNEL string, LISTNAME string, VALUE []byte) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("LPUSH", LISTNAME, VALUE)

	return err
}

func LpushString(RedisIP string, CHANNEL string, LISTNAME string, VALUE string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("LPUSH", LISTNAME, VALUE)

	return err
}

func LpushInt(RedisIP string, CHANNEL string, LISTNAME string, VALUE int) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("LPUSH", LISTNAME, VALUE)

	return err
}

func Rpush(RedisIP string, CHANNEL string, LISTNAME string, VALUE string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("RPUSH", LISTNAME, VALUE)

	return err
}

func Lpop(RedisIP string, CHANNEL string, LISTNAME string) (string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.String(c.Do("LPOP", LISTNAME))

	return content, err
}

func Rpop(RedisIP string, CHANNEL string, LISTNAME string) (string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.String(c.Do("RPOP", LISTNAME))

	return content, err
}

func Lrange(RedisIP string, CHANNEL string, LISTNAME string, START int, END int) ([]string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.Strings(c.Do("LRANGE", LISTNAME, START, END))

	return content, err
}

func Llen(RedisIP string, CHANNEL string, LISTNAME string) (int, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	num, err := redis.Int(c.Do("LLEN", LISTNAME))

	return num, err
}

func Sadd(RedisIP string, CHANNEL string, SETNAME string, VALUE string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = redis.String(c.Do("SADD", SETNAME, VALUE))

	return err
}

func Smembers(RedisIP string, CHANNEL string, SETNAME string) ([]string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.Strings(c.Do("SMEMBERS", SETNAME))

	return content, err
}

func Srem(RedisIP string, CHANNEL string, KEY_NAME string, MEMBER_NAME string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("SREM", KEY_NAME, MEMBER_NAME)
	CheckError(err)

	return err
}

func Hkeys(RedisIP string, CHANNEL string, NAME string) ([]string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	allHkey, err := redis.Strings(c.Do("HKEYS", NAME))

	return allHkey, err
}

func Hget(RedisIP string, CHANNEL string, NAME string, KEY_NAME string) (string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	allHcontent, err := redis.String(c.Do("HGET", NAME, KEY_NAME))

	return allHcontent, err
}

func Hset(RedisIP string, CHANNEL string, NAME string, KEY_NAME string, VAL string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("HSET", NAME, KEY_NAME, VAL)

	return err
}

func Hmget(RedisIP string, CHANNEL string, KEY string, FIELD string) ([]string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	content, err := redis.Strings(c.Do("HMGET", KEY, FIELD))

	return content, err
}

func Flushdb(RedisIP string, CHANNEL string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	_, err = c.Do("FLUSHDB")

	return err
}

func Info(RedisIP string) (string, error) {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()
	c.Do("AUTH", Password)
	CheckError(err)

	allcontent, err := redis.String(c.Do("INFO"))

	return allcontent, err
}

// it could be better
func SetBlacklist(RedisIP string, CHANNEL string, URL string) error {
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH", Password)
	_, err = c.Do("SELECT", CHANNEL)
	CheckError(err)

	START := time.Now().Format("20060102150405")
	_, err = c.Do("SET", URL, START)
	// CheckError(err)
	// _, err = c.Do("EXPIRE",URL,"43200")
	return err
}

/*
func Redis_DelAllUser()error{
	RedisIPPORT := fmt.Sprintf("%s:%s", RedisIP, RedisPort)
	c, err := redis.Dial("tcp", RedisIPPORT)
	CheckError(err)
	defer c.Close()

	c.Do("AUTH",Password)
	user := Redis_AllUser()
	var err2 error = nil
	for i:=0 ; i<len(user) ; i++{
		_, err2 = c.Do("DEL", user[i])
		CheckError(err2)
		if err2 != nil{
			return err2
		}
	}
	return err2
}
*/

func CheckError(err error) {
	if err != nil {
		log.Println("Error: ", err)
		// os.Exit(0)
	}
}
