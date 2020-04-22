package helper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"one/models"
)

// GenToken 生成JWT
func GenToken(claims *models.Claim, expiresAt int64, issuer string) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		// 过期时间
		ExpiresAt: expiresAt,
		// 签发人
		Issuer: issuer,
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(models.JwtSecret)
}

// 解析token
func ParseToken(tokenString string) (*models.Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claim{}, func(token *jwt.Token) (i interface{}, err error) {
		return models.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验 token
	if claims, ok := token.Claims.(*models.Claim); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Token 无效 ")
}

// 从jwt获取用户信息
func GetJwtClaimsByContext(c *gin.Context) *models.Claim {
	jwtClaims, hasClaims := c.Get(models.JwtClaimsKey)
	if !hasClaims {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    "用户信息丢失",
			"data":   "",
		})
		return &models.Claim{}
	}
	claims, ok := jwtClaims.(*models.Claim)
	if !ok {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"msg":    "用户信息断言失败",
			"data":   "",
		})
		return &models.Claim{}
	}
	return claims
}
