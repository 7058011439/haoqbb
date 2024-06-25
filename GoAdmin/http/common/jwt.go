package common

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/config"
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type MyJwt struct {
	jwt.StandardClaims
	Data interface{}
}

func NewToken(data interface{}, expires time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyJwt{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expires).Unix(),
		},
	})
	return token.SignedString(config.HttpJWTKey())
}

var tokenList = Stl.NewDoubleMap()

type UserType = int

const (
	UserTypeAdmin UserType = iota
)

func marshal(iType int, id int64) int64 {
	// 这里玩家和管理员都用的同一套验证机制，然后通过位移的方式将类型和id组合成新id。当id超过0x0000FFFFFFFFFFFF(42E*65535)的时候会出错,理论上到地球毁灭不会溢出的
	return int64(iType)<<48 + id
}

func unMarshal(data int64) (iType int, id int64) {
	iType = int(data >> 48)        // 取出高 16 位作为 iType
	id = data & 0x0000FFFFFFFFFFFF // 取出低 48 位作为 id
	return iType, id
}

func UpdateToken(eType UserType, id int64, token string) {
	key := marshal(eType, id)
	if oldToken := tokenList.GetValue(key); oldToken != nil {
		tokenList.RemoveByValue(oldToken)
	}
	if token != "" {
		tokenList.Add(key, token)
	}
}

func isTokenValid(token string, eType UserType) bool {
	return true // todo
	if key := tokenList.GetKey(token); key != nil {
		iType, _ := unMarshal(key.(int64))
		return iType == eType
	}
	return false
}

func ParseToken(token string, eType UserType) (interface{}, error) {
	if isTokenValid(token, eType) {
		parsedToken, err := jwt.ParseWithClaims(token, &MyJwt{}, func(token *jwt.Token) (interface{}, error) {
			return config.HttpJWTKey(), nil
		})
		if err != nil {
			// 如果解析或验证 token 时发生错误，直接返回这个错误
			return nil, fmt.Errorf("%v, %w", ResponseTokenError, err)
		}

		// 如果 token 的签名验证成功，但其他声明无效（如已过期）
		if !parsedToken.Valid {
			return nil, fmt.Errorf("%v, 签名验证成功，但其他声明无效", ResponseTokenError)
		}

		// 提取自定义的 claims
		claims, ok := parsedToken.Claims.(*MyJwt)
		if !ok {
			// 理论上永远不会来
			return nil, fmt.Errorf("%v, 数据类型错误", ResponseTokenError)
		}

		// 返回 Data 字段
		return claims.Data, nil
	}

	return nil, fmt.Errorf("token 无效, 未知token")
}

func CheckAdminToken(c *gin.Context) {
	ret := Http.NewResult(c)
	token := c.Request.Header.Get("token")
	if token == "" {
		token = c.Query("token")
	}

	if data, err := ParseToken(token, UserTypeAdmin); err == nil {
		SetClaim(c, data)
		path := fmt.Sprintf("%-10s%s", c.Request.Method, c.FullPath())
		if err := checkPermission(c, path); err != nil {
			ret.Abort(err.Error(), nil)
		} else {
			insertOperateLog(c, path)
		}
	} else {
		ret.Abort(err.Error(), nil)
	}
}

func getClaim(c *gin.Context) interface{} {
	ret, _ := c.Get(CLAIMS)
	return ret
}

func SetClaim(c *gin.Context, data interface{}) {
	c.Set(CLAIMS, data)
}

func GetInt64(c *gin.Context, key string) int64 {
	if data, ok := getClaim(c).(map[string]interface{}); ok {
		return int64(data[key].(float64))
	} else {
		Log.ErrorLog("当前请求token解析异常")
		return 0
	}
}

func GetString(c *gin.Context, key string) string {
	if data, ok := getClaim(c).(map[string]interface{}); ok {
		return data[key].(string)
	} else {
		Log.ErrorLog("当前请求token解析异常")
		return ""
	}
}

func GetAdminId(c *gin.Context) int64 {
	return GetInt64(c, TokenKeyId)
}

func GetAdminRoleId(c *gin.Context) int64 {
	return GetInt64(c, TokenKeyRoleId)
}

func GetCurrAdmin(c *gin.Context) *admin.User {
	return admin.GetAdmin(GetAdminId(c))
}
