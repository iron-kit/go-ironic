package ironic

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pinyin "github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"

	"go/ast"
	"reflect"
	"sync"
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

func getTypeName(t reflect.Type) string {
	// name := ""
	if t.Kind() == reflect.Ptr {
		return t.Elem().String()
	}

	return t.String()
}

func IsZero(v interface{}) bool {
	vv := reflect.ValueOf(v)

	return isZero(vv)
}

var typeTime = reflect.TypeOf(time.Time{})

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		if vt == typeTime {
			return v.Interface().(time.Time).IsZero()
		}
		for i := 0; i < v.NumField(); i++ {
			if vt.Field(i).PkgPath != "" && !vt.Field(i).Anonymous {
				continue // Private field
			}
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

type StructInfo struct {
	Type         reflect.Type
	StructFields []*StructField
}

type StructField struct {
	Name        string
	Index       int
	IsNormal    bool
	IsIgnored   bool
	IsInline    bool
	Tag         reflect.StructTag
	TagMap      map[string]string
	Struct      reflect.StructField
	Zero        reflect.Value
	InlineIndex []int
}

var structMap sync.Map

func parseTagConfig(tags reflect.StructTag) map[string]string {
	conf := map[string]string{}
	for index, str := range []string{tags.Get("bson"), tags.Get("monger")} {
		// bson
		if index == 0 {
			tags := strings.Split(str, ",")

			for _, tag := range tags {
				// if t == "" {
				// 	continue
				// }
				switch tag {
				case "inline":
					conf["INLINE"] = "true"
				case "omitempty":
					conf["OMITEMPTY"] = "true"
				default:
					conf["COLUMN"] = tag
				}
			}

			continue
		}
		tags := strings.Split(str, ",")
		for _, value := range tags {
			v := strings.Split(value, "=")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "COLUMN" {
				continue
			}
			if len(v) >= 2 {
				conf[k] = strings.Join(v[1:], "=")
			} else {
				conf[k] = k
			}
		}
	}

	return conf
}

func getStructInfo(d interface{}) *StructInfo {
	var structInfo StructInfo
	if d == nil {
		return &structInfo
	}

	reflectV := reflect.ValueOf(d)
	reflectType := reflectV.Type()

	if reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	if reflectType.Kind() != reflect.Struct {
		return &structInfo
	}

	if v, found := structMap.Load(reflectType); found && v != nil {
		return v.(*StructInfo)
	}

	structInfo.Type = reflectType

	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {

			field := &StructField{
				Struct:   fieldStruct,
				Name:     fieldStruct.Name,
				Tag:      fieldStruct.Tag,
				TagMap:   parseTagConfig(fieldStruct.Tag),
				Zero:     reflect.New(fieldStruct.Type).Elem(),
				Index:    i,
				IsInline: false,
			}

			// hidden
			if _, found := field.TagMap["-"]; found {
				field.IsIgnored = true
			} else if v, foundInline := field.TagMap["INLINE"]; foundInline && v == "true" {
				// the field is inline
				inlineFieldStruct := getStructInfo(reflect.New(fieldStruct.Type).Interface())

				for _, inlineField := range inlineFieldStruct.StructFields {
					inlineField.IsInline = true
					// inlineField.Index = []int{i, field.Index[0]}
					inlineField.InlineIndex = []int{i, inlineField.Index}
					structInfo.StructFields = append(structInfo.StructFields, inlineField)
				}
				continue
			} else {

				indirectType := fieldStruct.Type
				for indirectType.Kind() == reflect.Ptr {
					indirectType = indirectType.Elem()
				}

				fieldValue := reflect.New(indirectType).Interface()
				if _, isTime := fieldValue.(*time.Time); isTime {
					field.IsNormal = false
				} else {
					switch fieldStruct.Type.Kind() {
					case reflect.Slice:
						field.IsNormal = false
					case reflect.Struct:
						fallthrough
					case reflect.Ptr:
						field.IsNormal = false
					default:
						field.IsNormal = true
					}
				}
			}
			structInfo.StructFields = append(structInfo.StructFields, field)
		}
	}

	structMap.Store(reflectType, &structInfo)
	return &structInfo
}
