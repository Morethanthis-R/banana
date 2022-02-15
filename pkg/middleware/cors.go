package middleware

import (
	"github.com/gin-gonic/gin"
	h1 "net/http"
	"regexp"
)

//func Cors() middleware.Middleware {
//	return func(handler middleware.Handler) middleware.Handler {
//		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
//
//			tr, _ := transport.FromServerContext(ctx)
//			method := tr.RequestHeader().Get("Method")
//			origin := tr.RequestHeader().Get("Origin")
//			var filterHost = [...]string{"http://localhost.*"}
//			// filterHost 做过滤器，防止不合法的域名访问
//			var isAccess = false
//			for _, v := range(filterHost) {
//				match, _ := regexp.MatchString(v, origin)
//				if match {
//					isAccess = true
//				}
//			}
//			if isAccess {
//				// 核心处理方式
//				tr.ReplyHeader().Set("Access-Control-Allow-Origin",origin)
//				tr.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
//				tr.ReplyHeader().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
//				tr.ReplyHeader().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
//				tr.ReplyHeader().Set("Access-Control-Allow-Credentials", "true")
//				tr.ReplyHeader().Set("content-type", "application/json")
//			}
//			if method =="OPTIONS" {
//				tr.ReplyHeader().Set("Access-Control-Max-Age","1728000")
//				return handler(ctx, req)
//			}
//			return handler(ctx, req)
//		}
//	}
//}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var filterHost = [...]string{"http://localhost.*"}
		// filterHost 做过滤器，防止不合法的域名访问
		var isAccess = false
		for _, v := range filterHost {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}
		if isAccess {
			// 核心处理方式
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Set-Cookie,Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Set-Cookie, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.Header("Access-Control-Max-Age", "1728000")
			c.JSON(h1.StatusOK, "Options Request!")
			return
		}
		c.Next()
	}
}
