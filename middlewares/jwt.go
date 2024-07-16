package middlewares

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"gm/global"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

type CustomClaims struct {
	ID             int
	Username       string
	Name           string
	RoleChannelIds []int
	jwt.StandardClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			rsp.ErrorJSON(c, ecode.UserNotLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			rsp.ErrorJSON(c, ecode.TokenFormatError)
			c.Abort()
			return
		}
		j := NewJwt()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(parts[1])
		if err != nil {
			if err == TokenExpired {
				rsp.ErrorJSON(c, ecode.TokenExpired)
				c.Abort()
				return
			}

			rsp.ErrorJSON(c, ecode.UserNotLogin)
			c.Abort()
			return
		}

		c.Set("uid", claims.ID)
		c.Set("name", claims.Name)
		c.Set("roleChannelIds", claims.RoleChannelIds)
		c.Next()
	}
}

type Jwt struct {
	Key []byte `json:"key"`
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJwt() *Jwt {
	return &Jwt{Key: []byte(global.ServerConfig.JwtInfo.SigningKey)}
}

// 创建token
func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}

// 解析token
func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// 更新token
func (j *Jwt) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
