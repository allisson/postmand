package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

// NewSecureHeadersMiddleware returns a new secure headers middleware
func NewSecureHeadersMiddleware() func(next http.Handler) http.Handler {
	options := secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "frame-ancestors 'none'",
	}
	return secure.New(options).Handler
}
