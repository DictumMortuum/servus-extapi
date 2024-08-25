package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.9",
	}
	c.AbortWithStatusJSON(200, rs)
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = supertokens.Init(SuperTokensConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: Origins,
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: append([]string{"content-type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	// Adding the SuperTokens middleware
	r.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	})

	r.GET("/auth/version", Version)
	r.GET("/auth/sessioninfo", verifySession(nil), sessionInfo)
	r.GET("/auth/userinfo", verifySession(nil), userInfo)
	r.Run(":10004")
}

func verifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.Request = c.Request.WithContext(r.Context())
			c.Next()
		})(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	}
}

func sessionInfo(c *gin.Context) {
	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
	w := c.Writer
	r := c.Request

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}

	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}

	w.WriteHeader(200)
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}

func userInfo(c *gin.Context) {
	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
	w := c.Writer

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}

	userID := sessionContainer.GetUserID()
	userInfo, err := thirdparty.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("could not get user by ID"))
		return
	}
	if userInfo == nil {
		emailPasswordUserInfo, err := emailpassword.GetUserByID(userID)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("could not get user by ID"))
			return
		}

		w.WriteHeader(200)
		w.Header().Add("content-type", "application/json")
		bytes, err := json.Marshal(emailPasswordUserInfo)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error in converting to json"))
		} else {
			w.Write(bytes)
		}
	} else {
		fmt.Println(userInfo)
	}
}
