package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

type AuthenService struct {
	jwtVerifier *jwtverifier.JwtVerifier
}

func NewAuthenService() *AuthenService {
	tv := map[string]string{}
	//tv["aud"] = app.Config.Okta.OAuth2.Audience
	//tv["cid"] = app.Config.Okta.OAuth2.ClientId
	//tv["scp"] = app.Config.Okta.OAuth2.Scopes

	jvConfig := jwtverifier.JwtVerifier{
		//Issuer:           app.Config.Okta.OAuth2.Issuer,
		ClaimsToValidate: tv,
	}

	return &AuthenService{
		jwtVerifier: jvConfig.New(),
	}
}

// redirect and return un-authorized exception
func (authSvc *AuthenService) RedirectURL(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token"})
	c.Abort()
}

func (authSvc *AuthenService) Verify(bearerToken string) error {
	tokenString := strings.Split(bearerToken, "Bearer ")

	if len(tokenString) > 1 {
		jwt, err := authSvc.jwtVerifier.VerifyAccessToken(tokenString[1])
		fmt.Println(jwt.Claims)
		return err
	}

	return errors.New("bearer token requires")
}
