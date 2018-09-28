package utils

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pinyin "github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"
)

/*
StringArrayContainer function
检测数组里面是否包含某个字符串
*/
func StringArrayContainer(array []interface{}, found interface{}) bool {
	for _, v := range array {
		if v.(string) == found.(string) {
			return true
		}
	}

	return false
}

/*
GenerateRangeNum is to generate a random num between min and max
e.g. GenerateRangeNum(0, 9) => num of [0,9]
*/
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

/*
GenerateVerifyCode 以时间作为随机因子生成验证码
*/
func GenerateVerifyCode(count int) string {
	if count < 4 {
		panic("count must getter than or equal to 4")
	}

	min := int(math.Pow(10.0, float64(count-1)))
	max := int(math.Pow(10.0, float64(count))) - 1
	code := GenerateRangeNum(min, max)

	return strconv.Itoa(code)
}

// GeneratePassword is a func to generate password
func GeneratePassword(password string) (string, error) {
	pwdByte := []byte(password)
	bcryptPassword, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptPassword[:]), nil
}

// CheckPassword check the password
// return a nil on success
func CheckPassword(plainPassword string, password string) error {
	hashedPassword := []byte(password)
	pwd := []byte(plainPassword)
	return bcrypt.CompareHashAndPassword(hashedPassword, pwd)
}

// GenerateToken 生成Token
func GenerateToken(user string) (string, error) {
	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	// // Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = user + fmt.Sprintf("_%d", time.Now().Unix())
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("c13be55b40cf9dacb8231156ff28d41e65c8b48b"))
	if err != nil {
		return "", err
	}
	return t, nil
}

// DecodeToken is a func get token from string
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		return []byte("c13be55b40cf9dacb8231156ff28d41e65c8b48b"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// GetDefaultAvatar 获取默认头像
func GetDefaultAvatar() {

}

// GetNowUnixTime 获取当前的Unix时间(13位)
func GetNowUnixTime() int64 {
	now := time.Now()
	return now.UnixNano() / 1e6
}

func MicroTimeFormat(microUnix int64, format string) string {
	t := time.Unix(microUnix/1e3, 0)
	return t.Format(format)
}

// Time2MicroUnix 将时间转化为 13位 Unix 时间
func Time2MicroUnix(time *time.Time) int64 {
	return time.UnixNano() / 1e6
}

// Hans2Pinyin 将汉字转换为拼音
func Hans2Pinyin(hans string, sep string) string {
	pinyinArgs := pinyin.NewArgs()
	pinyinNameArr := pinyin.Pinyin(hans, pinyinArgs)
	pinyinName := []string{}
	for _, v := range pinyinNameArr {
		pinyinName = append(pinyinName, v[0])
	}

	return strings.Join(pinyinName, sep)
}

func Unicode2Hans(unicode string) (string, error) {
	var context string
	sUnicode := strings.Split(unicode, "\\u")
	for _, v := range sUnicode {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			return "", err
		}
		context += fmt.Sprintf("%c", temp)
	}

	return context, nil
}

func ISOTime2Time(t ...string) (time.Time, error) {
	defaultLayout := "2006-01-02 15:04:05.000 -0700 UTC"
	// defaultLayout := time.RFC1123

	if len(t) <= 0 {
		// panic("Need a time string")
		return time.Now(), errors.New("Need a time string")
	}

	if len(t) > 1 {
		defaultLayout = t[1]
	}

	// 2018-09-06 09:55:42.405 +0000 UTC
	return time.Parse(defaultLayout, t[0])
}

func ISOTime2MicroUnix(t ...string) int64 {
	time, err := ISOTime2Time(t...)

	if err != nil {
		return 0
	}

	return Time2MicroUnix(&time)
}

func IsZero(v interface{}) bool {
	t := reflect.TypeOf(v)

	if !t.Comparable() {
		return false
	}

	// val := reflect.ValueOf(v)
	return reflect.DeepEqual(v, reflect.Zero(t).Interface())
}
