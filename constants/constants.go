package constants

const (
	JwtSecretEnvKey = "JWT_SECRET"
)

const (
	RequestIDKey = "X-Request-Id"
	EmailKey     = "Email"
)

type ContextKey string

var RequestIDCtxKey = ContextKey("X-Request-Id")
var EmailCtxKey = ContextKey("Email")
