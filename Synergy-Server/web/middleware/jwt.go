package middleware

import (
	"ARTeamViewService/global"
	"ARTeamViewService/models"
	"ARTeamViewService/utils"
	"github.com/dgrijalva/jwt-go"
	jwtMiddleWare "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

var tokenSecret = utils.RandomCharString(16)

func JwtHandler() *jwtMiddleWare.Middleware {
	return jwtMiddleWare.New(jwtMiddleWare.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func JwtAuthToken(ctx iris.Context) {
	token := ctx.Values().Get("jwt").(*jwt.Token)
	tokenStr := token.Raw
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		global.GLogger.Info("token, ok: ", ok, token.Valid, claims)
		global.GLogger.Info("expired: ", claims["exp"])
		global.GLogger.Info("userid: ", claims["userid"])

		rdsKey := utils.RedisUserKey(claims["userid"].(string))
		tokenRds := global.GRedisClient.Get(global.Context, rdsKey).Val()

		if utils.Equals(tokenStr, tokenRds) {
			ctx.Values().Set("userid", claims["userid"])
			ctx.Next()
		} else {
			global.GLogger.Info("token different")
			ctx.Values().Set("userid", "")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.StopExecution()
			return
		}
	} else {
		global.GLogger.Info("invalid token: ", ok, token.Valid, claims)
		ctx.Values().Set("userid", "")
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.Next()
	}
}

func GenerateToken(userId string) (string, error) {
	timeNow := time.Now()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	var deltaTime time.Duration = time.Duration(global.GConfig.TokenExpire) * time.Second

	claims["exp"] = int(timeNow.Add(deltaTime).Unix())
	claims["iat"] = int(timeNow.Unix())
	claims["userid"] = userId
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	oauthToken := new(models.OauthToken)
	oauthToken.Token = tokenString
	oauthToken.UserId = userId
	oauthToken.Secret = ""
	oauthToken.Revoked = false
	oauthToken.ExpireAt = timeNow.Add(deltaTime).Unix()
	oauthToken.CreateAt = timeNow.Unix()
	global.GLogger.Info(oauthToken)
	return tokenString, nil
}

// JWTUserId 获取jwt保存的userid
func JWTUserId(ctx iris.Context) string {
	userid := ctx.Values().GetString("userid")
	return userid
}
