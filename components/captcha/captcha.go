package captcha

import (
	"github.com/make-money-fast/captcha"
	"github.com/make-money-fast/gin"
)

func Captcha(w, h int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		captcha.Server(w, h).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
