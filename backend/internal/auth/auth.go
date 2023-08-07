package auth

import (
	"os"
	"serverless/internal/requests"
	"serverless/pkg/tweetlock/datastructs"
	"serverless/pkg/tweetlock/dtos"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	User dtos.User `json:"user"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	jwt.StandardClaims
}

func Guard(r events.APIGatewayProxyRequest) (string, bool) {
	accessCookie := requests.GetCookieValue(r.Headers["Cookie"], "x-access")

	if accessCookie == "" {
		return "", false
	}

	c := &AccessClaims{}
	_, err := jwt.ParseWithClaims(accessCookie, c, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return "", false
	}

	if c.NotBefore == 0 || c.ExpiresAt < time.Now().Unix() || c.NotBefore > time.Now().Unix() {
		return "", false
	}

	return c.Subject, true
}

func SignAccessToken(user datastructs.User) (string, time.Time, error) {
	expires := time.Now().Add(time.Minute * 10)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		User: dtos.User{
			ID:          user.ID,
			TwitterID:   user.TwitterID,
			TwitterName: user.Handle,
			Name:        user.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   user.ID,
		},
	})

	signed, err := tok.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", time.Unix(0, 0), err
	}

	return signed, expires, nil
}
