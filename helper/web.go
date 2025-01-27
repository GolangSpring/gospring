package helper

import (
	"bytes"
	"context"
	"fmt"
	securityService "github.com/GolangSpring/gospring/pkg/security/service"
	"github.com/go-fuego/fuego"
	"io"
	"net/http"
	"time"
)

const CookieKey = "token"

func ReadRequestBody(request *http.Request) ([]byte, error) {
	body := request.Body

	// Read the body
	var buf bytes.Buffer
	if body != nil {
		_, err := io.Copy(&buf, body)
		if err != nil {
			return []byte{}, err
		}
	}
	// Reset the body so it can be read again
	request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	return buf.Bytes(), nil
}

func GetUserFromContext(c context.Context) (*securityService.UserClaims, error) {
	optionalUser := c.Value("user")
	user, ok := optionalUser.(*securityService.UserClaims)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func WriteTokenCookie[T any](c fuego.ContextWithBody[T], token string, expiration time.Duration) {
	cookie := http.Cookie{
		Name:     CookieKey,
		Value:    token,
		Expires:  time.Now().Add(expiration),
		Path:     "/",
		HttpOnly: true,
	}

	c.SetCookie(cookie)
}
