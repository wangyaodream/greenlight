package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/wangyaodream/greenlight/internal/validator"
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

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
    v.Check(tokenPlaintext != "", "token", "must be provided")
    v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
    DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
    token, err := generateToken(userID, ttl, scope)
    if err != nil {
        return nil, err
    }

    err = m.Insert(token)
    return token, err
}

func (m TokenModel) Insert(token *Token) error {
    query := `
        INSERT INTO tokens (hash, user_id, expiry, scope)
        VALUES ($1, $2, $3, $4)
    `

    args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

    _, err := m.DB.Exec(query, args...)
    return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
    query := `
        DELETE FROM tokens
        WHERE scope = $1 AND user_id = $2
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := m.DB.ExecContext(ctx, query, scope, userID)
    return err
}
