package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)


const (
    ScopeActivation = "activation"
)

type Token struct {
    Plaintext string
    Hash []byte
    UserID int64
    Expiry time.Time
    Scope string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
    // 创建一个Token实例
    token := &Token{
        UserID: userID,
        Expiry: time.Now().Add(ttl),
        Scope: scope,
    }

    // 初始化一个空的字节切片
    randomBytes := make([]byte, 16)

    _, err := rand.Read(randomBytes)
    if err != nil {
        return nil, err
    }

    // 使用base32编码生成一个随机的字符串
    token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

    // 使用SHA-256算法生成一个哈希值
    hash := sha256.Sum256([]byte(token.Plaintext))
    // 将哈希值赋值给Token实例的Hash字段
    token.Hash = hash[:]

    return token, nil
}
