package captcha

import (
	"github.com/clearcodecn/gin"
	"github.com/make-money-fast/captcha"
)

func Captcha(w, h int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		captcha.Server(w, h).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
