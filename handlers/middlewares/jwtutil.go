package middlewares

import (
	jwt "github.com/dgrijalva/jwt-go"
	googleUuid "github.com/google/uuid"
	"github.com/mascot27/myAwesomeWebApi/config"
	"github.com/mascot27/myAwesomeWebApi/models"
	"time"
)

func GetToken(name string) (string, error) {
	signingKey := []byte(config.JWT_SIGNING_KEY)

	tokenIssuer := "web"

	rawId, _ := googleUuid.NewRandom()
	tokenUuid := rawId.String()

	secondsOfValidity := int64(1000)
	tokenExp := float64(time.Now().Unix() + secondsOfValidity)

	userUuid, _ := models.GetUuidForName(name)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss:":        tokenIssuer,  // issuer
		"exp":         tokenExp,  // expiration time
		"isSingleUse": true,  // begin customs claims
		"userUuid":    userUuid,
		"tokenUuid":   tokenUuid,
	})

	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(config.JWT_SIGNING_KEY)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, err
}

/*
good read on json token:
(taken from SO)

Here is the a solution called as JWT old for new exchange schema.

Because we can’t invalidate the issued token before expire time, we always use short-time token, such as 30 minute. When the token expired, we use the old token exchange a new token. The critical point is one old token can exchange one new token only.

In center auth server, we maintain a table like this:

// ->> it's a whitelist of tokens
table auth_tokens(
	user_id,
	jwt_hash,
	expire
)

- user_id contained in JWT string.
- jwt_hash is a hash value of whole JWT string,Such as SHA256. Can be replaced by the token id or something else who identify the token
- expire field is optional.


The following is work flow:

	1) User request the login API with username and password, the auth server issue one token, and register the token ( add one row in the table. )
	2) When the token expired, user request the exchange API with the old token. Firstly the auth server validate the old token as normal except expire checking, then create the token hash value, then lookup above table by user id:
		If found record and user_id and jwt_hash is match, then issue new token and update the table.
		If found record, but user_id and jwt_hash is not match , it means someone has use the token exchanged new token before. The token be hacked, delete records by user_id and response with alert information.
		if not found record, user need login again or only input password.
	when use changed the password or login out, delete record by user id.

To use token continuously ,both legal user and hacker need exchange new token continuously, but only one can succeed, when one fails, both need to login again at next exchange time.

So if hacker got the token, it can be used for a short time, but can't exchange for a new one if a legal user exchanged new one next time, because the token validity period is short. It is more secure this way.

If there is no hacker, normal user also need exchange new token periodically ,such as every 30 minutes, this is just like login automatically. The extra load is not high and we can adjust expire time for our application.


design flaws:
	- cannot use multiple devices
	- no granular access can be granted


--------------------------------------------
Usage of refresh tokens
-------------------------------------------

this design approach separate authentication and identification

firstly, we approach the authentication

table active_refresh_tokens (
	jwt_token_id,
	token_user_uid
	token_genetic
)

table used_refresh token (
	jwt_token_id,
	token_user_uid
	token_genetic
)

table genetics_blacklist (
	user_id
	tokens_genetic
)


The following is workflow
	1) User request the auth API with username and password, the auth server issue one token, and register the token in the active_tokens database's table
		- save refresh token in the actives refresh_tokens database's table, genetic is like the device id
		- return an refresh token and an access token
	2) when user request a ressource with an expired token, he need to refresh the token by requesting a new one to the auth server
		a) if found record in active record, issue new access and refresh token, moves the refresh token to the used list
			check also if genetic is blacklisted for this user
		b) if found record in used record, revoke all token of this genetic, a hacker might posses the refresh token but we cannot identify the user and the hacker, so we take precautions
			- blacklist token and return error "you might been hacked"
		c) else token is invalid return bad request



potential risk:
	- refresh token can be used by hacker if leaked
		-> so a new refresh token should be issued with each access token, and save it in a blacklist with refer the token_id with the user if the refresh token is reused, this might be a sign of MIT attack so revoke all the tokens of this tree of token
	- we use a whitelist because if the database is no longer accessible no unauthorized access can be granted


note:
	- these risks are eliminated by using good rock-solid implementation
	- this enable multiple device
	- this enable the use of fine grained access to ressource


-----------------------------------------------------------------------------------------------------------------------
How the resources can validate the use of the token?





Remarks:
		 - the second use is great for our use case but need a slightly more complex architecture based on micro-services pattern.



source: http://www.jianshu.com/p/b11accc40ba7

 */

