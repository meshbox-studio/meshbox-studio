package auth

import (
	"github.com/danielgtaylor/huma/v2"
)

var store = NewSessionStore()

func RequireSession(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token, err := huma.ReadCookie(ctx, "meshbox_session")
		if err != nil || !store.isValid(token.Value) {
			huma.WriteErr(api, ctx, 401, "valid session required")
			return
		}
		next(ctx)
	}
}

func (s *SessionStore) isValid(token string) bool {
	_, ok := s.Get(token)
	return ok
}