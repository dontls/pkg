package ginx

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type jwtClaims struct {
	UserClaims
	jwt.StandardClaims
}

// JWT 定义
type jwtToken struct {
	SigningKey []byte
}

var (
	ErrTokenMalformed    error = errors.New("token not event")
	ErrTokenExpired      error = errors.New("token expired")
	ErrTokenNotValidYet  error = errors.New("token not valid yet")
	ErrTokenInValidation error = errors.New("token invalidation")
)

func signKey() string {
	return "ginfast"
}

var JWTToken = &jwtToken{
	[]byte(signKey()),
}

// createToken 创建Token
func (j *jwtToken) createToken(claims jwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *jwtToken) Parse(tokenString string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInValidation
			}
		}
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInValidation
}

// RefreshToken 刷新token
func (j *jwtToken) Refresh(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.createToken(*claims)
	}
	return "", ErrTokenInValidation
}

// JWTToken 生成Token
func (j *jwtToken) Create(v *UserClaims) (string, error) {
	claims := jwtClaims{}
	claims.StandardClaims = jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
		ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
		Issuer:    string(j.SigningKey),            // 签名的发行者
	}
	claims.UserClaims = *v
	return j.createToken(claims)
}

// JWTAuth 鉴权
func JWTAuth(r *gin.RouterGroup) *gin.RouterGroup {
	r.Use(func(c *gin.Context) {
		// jwt鉴权取头部信息 x-token
		// 登录时回返回token信息
		// 前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("Authorization")
		ctx := JSON(c)
		if !strings.Contains(token, "Bearer ") {
			ctx.JSONWriteMsg(StatusForbidden, ErrTokenExpired)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := JWTToken.Parse(token)
		if err != nil {
			// if err == middleware.TokenValidationErrorExpired {
			// 	c.Abort()
			// 	return
			// }
			ctx.JSONWriteMsg(StatusLoginExpired, err)
			c.Abort()
			return
		}
		c.Set("claims", claims.UserClaims)
		c.Next()
	})
	return r
}
