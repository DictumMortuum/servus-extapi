package main

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var SuperTokensConfig = supertokens.TypeInput{
	Supertokens: &supertokens.ConnectionInfo{
		ConnectionURI: "http://localhost:3567",
	},
	AppInfo: supertokens.AppInfo{
		AppName:       "Tables",
		APIDomain:     "https://auth.dictummortuum.com",
		WebsiteDomain: "https://tables.dictummortuum.com",
	},
	RecipeList: []supertokens.Recipe{
		thirdpartyemailpassword.Init(&tpepmodels.TypeInput{}),
		session.Init(nil),
		dashboard.Init(nil),
		thirdparty.Init(nil),
	},
}
