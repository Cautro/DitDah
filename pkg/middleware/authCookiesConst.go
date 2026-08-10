package middleware

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"

	AccessTokenCookieMaxAge  = 3600 * 72
	RefreshTokenCookieMaxAge = 3600 * 24 * 30
)