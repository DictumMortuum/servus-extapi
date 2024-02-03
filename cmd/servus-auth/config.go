package main

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var SuperTokensConfig = supertokens.TypeInput{
	Supertokens: &supertokens.ConnectionInfo{
		ConnectionURI: "http://sol.dictummortuum.com:3567",
	},
	AppInfo: supertokens.AppInfo{
		AppName:   "Tables",
		APIDomain: "https://auth.dictummortuum.com",
		// WebsiteDomain: "https://tables.dictummortuum.com",
		GetOrigin: func(request *http.Request, userContext supertokens.UserContext) (string, error) {
			if request != nil {
				origin := request.Header.Get("origin")
				if origin == "" {
					// this means the client is in an iframe, it's a mobile app, or
					// there is a privacy setting on the frontend which doesn't send
					// the origin
				} else {
					if origin == "https://tables.dictummortuum.com" {
						// query from the test site
						return "https://tables.dictummortuum.com", nil
					} else if origin == "http://localhost:3000" {
						// query from local development
						return "http://localhost:3000", nil
					}
				}
			}
			// in case the origin is unknown or not set, we return a default
			// value which will be used for this request.
			return "https://tables.dictummortuum.com", nil
		},
	},
	RecipeList: []supertokens.Recipe{
		thirdpartyemailpassword.Init(&tpepmodels.TypeInput{}),
		session.Init(nil),
		dashboard.Init(nil),
		thirdparty.Init(nil),
	},
}
