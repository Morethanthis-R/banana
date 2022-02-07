package middleware

import (
	"banana/pkg/ecode"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/metadata"
	"time"
)

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request header.
	authorizationKey string = "Authorization"
)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(200, gin.H{
				"message": "token为空",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil || claims == nil {
			c.JSON(200, gin.H{
				"message": "解析token错误",
			})
			c.Abort()
			return
		}
		cache := NewCache()
		key := fmt.Sprintf("account:cookie:id:%d", claims.UserId)
		res,err := cache.Get(c, key).Result()
		if err != nil {
			c.JSON(200, gin.H{
				"message": "检验失败",
			})
			c.Abort()
			return
		}
		if res != claims.Id{
			c.JSON(200, gin.H{
				"message": "该token不合法",
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(200, gin.H{
				"message": "token已过期",
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
		c.Set("x-md-global-uid", claims.UserId)
	}
}
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var token string
			tr, ok := transport.FromServerContext(ctx)
			if ok {
				authorization := tr.RequestHeader().Get(authorizationKey)
				if authorization != "" {
					token = authorization
				}
			}
			md, ok := metadata.FromIncomingContext(ctx)
			if ok && len(md.Get("Authorization")) > 0 {
				token = md.Get("Authorization")[0]
			}
			if token == "" {
				return nil, ecode.AUTH_FAIL.SetMessage("token为空")
			}
			//ctx = context.WithValue(ctx, "Authorization", token)
			// 解析token
			claims, err := ParseToken(token)
			fmt.Println(err)
			if err != nil {
				return nil, ecode.AUTH_FAIL.SetMessage("token解析错误")
			} else if time.Now().Unix() > claims.ExpiresAt {
				return nil, ecode.AUTH_FAIL.SetMessage("token已过期")
			}
			ctx = context.WithValue(ctx,"claims",claims)
			cache := NewCache()
			key := fmt.Sprintf("account:cookie:id:%d", claims.UserId)
			res,err := cache.Get(ctx, key).Result()
			if err != nil {
				return nil,ecode.REDIS_ERR.SetMessage("redis返回错误")
			}
			if res == ""{
				return nil,ecode.REDIS_ERR.SetMessage("该token已登出")
			}
			ctx = context.WithValue(ctx, "x-md-global-uid", claims.UserId)
			return handler(ctx, req)
		}
	}
}

func SetCookie(ctx context.Context, token string, role, expireSec int) string{

	res := fmt.Sprintf("p-token=%s; _urole=%d; _expires=%d", token, role, expireSec)
	return res
	//ctx.SetCookie("_token",token,expireSec,"","47.107.95.82",false,false)
	//ctx.SetCookie("_role",token,role,"","47.107.95.82",false,false)
}

func NewCache() *redis.Client {
	var options = &redis.Options{
		Addr:        "47.107.95.82:6379",
		Password:    "55882664",
		DB:          0,
	}
	client := redis.NewClient(options)
	if client == nil {
		return nil
	}
	client.AddHook(redisotel.TracingHook{})
	return client

}
