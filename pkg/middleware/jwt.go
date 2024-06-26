package middleware

import (
	"log"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

var jwksCache *sessmodels.GetJWKSResult = nil
var mutex sync.RWMutex
var JWKCacheMaxAgeInMs int64 = 60000

func getJWKSFromCacheIfPresent() *sessmodels.GetJWKSResult {
	mutex.RLock()
	defer mutex.RUnlock()
	if jwksCache != nil {
		// This means that we have valid JWKs for the given core path
		// We check if we need to refresh before returning
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)

		// This means that the value in cache is not expired, in this case we return the cached value
		//
		// Note that this also means that the SDK will not try to query any other Core (if there are multiple)
		// if it has a valid cache entry from one of the core URLs. It will only attempt to fetch
		// from the cores again after the entry in the cache is expired
		if (currentTime - jwksCache.LastFetched) < JWKCacheMaxAgeInMs {
			return jwksCache
		}
	}

	return nil
}

func getJWKS(coreUrl string) (*keyfunc.JWKS, error) {
	resultFromCache := getJWKSFromCacheIfPresent()

	if resultFromCache != nil {
		return resultFromCache.JWKS, nil
	}

	mutex.Lock()
	defer mutex.Unlock()
	// RefreshUnknownKID - Fetch JWKS again if the kid in the header of the JWT does not match any in
	// the keyfunc library's cache
	jwks, jwksError := keyfunc.Get(coreUrl, keyfunc.Options{
		RefreshUnknownKID: true,
	})

	if jwksError == nil {
		jwksResult := sessmodels.GetJWKSResult{
			JWKS:        jwks,
			Error:       jwksError,
			LastFetched: time.Now().UnixNano() / int64(time.Millisecond),
		}

		// Dont add to cache if there is an error to keep the logic of checking cache simple
		//
		// This also has the added benefit where if initially the request failed because the core
		// was down and then it comes back up, the next time it will try to request that core again
		// after the cache has expired
		jwksCache = &jwksResult

		return jwksResult.JWKS, nil
	}

	// This means that fetching from the core failed
	return nil, jwksError
}

func verifyJwt(url, jwtString string) (bool, error) {
	jwks, err := getJWKS(url)
	if err != nil {
		return false, err
	}

	parsedToken, parseError := jwt.Parse(jwtString, jwks.Keyfunc)
	if parseError != nil {
		log.Println(parseError)
		return false, err
	}

	if parsedToken.Valid {
		// claims, err := parsedToken.Claims.(jwt.MapClaims)
		// if !err {
		// 	return false, err
		// }
		return true, nil
		// Convert the claims to a key-value pair
		// claimsMap := make(map[string]interface{})
		// for key, value := range claims {
		// 	claimsMap[key] = value
		// }

		// emailVerifiedClaim, ok := claimsMap["st-ev"]
		// if ok {
		// 	emailVerified, ok := emailVerifiedClaim.(map[string]interface{})["v"].(bool)
		// 	if !ok {
		// 		panic("Claim format is wrong. Should never come here")
		// 	}

		// 	if !emailVerified {
		// 		// TODO: send a 403 to the client
		// 	}
		// } else {
		// TODO: send 403 to client
		// }

		// TODO: access token is valid

	} else {
		return false, nil
		// TODO: send a 401 to the client
	}
}

func Jwt(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Request.Header.Get("SA")
		verified, err := verifyJwt(url, h)
		if err != nil {

			c.AbortWithStatusJSON(401, gin.H{"error": "not authorized"})
		}

		if !verified {
			log.Println(verified)
			c.AbortWithStatusJSON(401, gin.H{"error": "not authorized"})
		}
	}
}

// func Body(c *gin.Context) {
// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	m, err := model.ToMap(c, "req")
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	m.Set("body", body)
// }
