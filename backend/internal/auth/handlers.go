package auth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Summary:     "Log in",
		Description: "Any email/password combination logs in as the demo user.",
		Tags:        []string{"Auth"},
	}, loginHandler)

	huma.Register(api, huma.Operation{
		OperationID: "auth-logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Summary:     "Log out",
		Description: "Clears the session cookie.",
		Tags:        []string{"Auth"},
	}, logoutHandler)

	huma.Register(api, huma.Operation{
		OperationID: "auth-session",
		Method:      http.MethodGet,
		Path:        "/api/auth/session",
		Summary:     "Get current session",
		Description: "Returns the authenticated user or 401.",
		Tags:        []string{"Auth"},
		Middlewares: huma.Middlewares{RequireSession(api)},
	}, sessionHandler)
}

type LoginInput struct {
	Body struct {
		Email string `json:"email" required:"true"`
	}
}

type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      struct {
		User User `json:"user"`
	}
}

type SessionOutput struct {
	Body struct {
		User User `json:"user"`
	}
}

func loginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user := DemoUser
	user.Email = input.Body.Email

	token := store.Create(user)

	resp := &LoginOutput{}
	resp.SetCookie = http.Cookie{
		Name:     "meshbox_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	resp.Body.User = user

	return resp, nil
}

func logoutHandler(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	resp := &LogoutOutput{}
	resp.SetCookie = http.Cookie{
		Name:     "meshbox_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	return resp, nil
}

type LogoutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func sessionHandler(ctx context.Context, input *SessionInput) (*SessionOutput, error) {
	user, ok := store.Get(input.Session.Value)
	if !ok {
		return nil, huma.Error401Unauthorized("invalid session")
	}

	return &SessionOutput{
		Body: struct {
			User User `json:"user"`
		}{User: user},
	}, nil
}

type SessionInput struct {
	Session http.Cookie `cookie:"meshbox_session"`
}