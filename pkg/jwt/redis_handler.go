package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"muxi_auditor/config"
	"strings"
	"time"
)

const BASENAME = "muxiAuthor:users:ssid:"

// RedisJWTHandler 实现了处理 JWT 的接口，并使用 Redis 进行支持
type RedisJWTHandler struct {
	cmd redis.Cmdable // Redis 命令接口，用于与 Redis 进行交互
	jwt *JWT
}

// NewRedisJWTHandler 创建并返回一个新的 RedisJWTHandler 实例
func NewRedisJWTHandler(cmd *redis.Client, conf *config.JWTConfig) *RedisJWTHandler {
	return &RedisJWTHandler{
		cmd: cmd, //redis实体
		jwt: NewJWT(time.Duration(conf.Timeout), conf.SecretKey),
	}
}

// ClearToken 清除客户端的 JWT ，并在 Redis 中记录已过期的会话
func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	// 要求客户端设置为空
	ctx.Header("JWT-Token", "")
	// 在 Redis 中记录登出的会话
	uc := ctx.MustGet("user").(UserClaims)
	return r.cmd.Set(ctx, BASENAME+uc.ID, "expired", uc.ExpiresAt.Time.Sub(time.Now())).Err()
}

// ExtractToken 从请求中提取并返回解析完成的结构体
func (r *RedisJWTHandler) ParseToken(ctx *gin.Context) (UserClaims, error) {
	//获取请求头
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return UserClaims{}, errors.New("Authorization请求头缺失")
	}

	//分割获取token
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return UserClaims{}, errors.New("请求头格式错误!")
	}
	//解析Token
	uc, err := r.jwt.ParseToken(segs[1])
	if err != nil {
		return UserClaims{}, err
	}
	//检查是否被列入黑名单
	ok, err := r.CheckSession(ctx, uc.ID)
	if err != nil || ok {
		return UserClaims{}, errors.New("session检验：失败")
	}

	return uc, nil
}

// SetJWTToken 生成并设置用户的 JWT
func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid uint, name string) error {
	tokenStr, err := r.jwt.SetJWTToken(uid, name)
	if err != nil {
		return err
	}
	ctx.Header("JWT-Token", tokenStr)
	return nil
}

// CheckSession 检查给定 ssid 的会话是否有效
func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) (bool, error) {
	val, err := r.cmd.Exists(ctx, BASENAME+ssid).Result()
	return val > 0, err
}
