package util

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/yushengguo557/register/common"
)

// GenerateToken 生成 token
func GenerateToken(openid, sessionKey string) (token string, err error) {
	token, err = jwt.NewWithClaims(jwt.SigningMethodES256, common.Claims{
		Openid:     openid,
		SessionKey: sessionKey,
	}).SignedString("")
	// 使用密钥对 token 进行签名
	if err != nil {
		return "", fmt.Errorf("签名 %w", err)
	}
	return token, nil
}
