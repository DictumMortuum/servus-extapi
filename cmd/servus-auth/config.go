package main

import (
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	Origins = []string{
		"https://tables.dictummortuum.com",
		"https://prices.dictummortuum.com",
		"https://boardgames.dictummortuum.com",
		"https://tools.dictummortuum.com",
		"http://localhost:3000",
	}
)

var SuperTokensConfig = supertokens.TypeInput{
	Supertokens: &supertokens.ConnectionInfo{
		ConnectionURI: "http://sol.dictummortuum.com:3567",
	},
	AppInfo: supertokens.AppInfo{
		AppName:   "DictumMortuum",
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
					for _, item := range Origins {
						if origin == item {
							return item, nil
						}
					}
				}
			}
			// in case the origin is unknown or not set, we return a default
			// value which will be used for this request.
			return "https://tables.dictummortuum.com", nil
		},
	},
	RecipeList: []supertokens.Recipe{
		emailpassword.Init(&epmodels.TypeInput{
			EmailDelivery: &emaildelivery.TypeInput{
				Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
					ogSendEmail := *originalImplementation.SendEmail

					(*originalImplementation.SendEmail) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
						// You can change the path, domain of the reset password link,
						// or even deep link it to your mobile app
						// This is: `${websiteDomain}${websiteBasePath}/reset-password`
						input.PasswordReset.PasswordResetLink = strings.Replace(
							input.PasswordReset.PasswordResetLink,
							"/auth/reset-password",
							"/#/auth/reset-password", 1,
						)
						return ogSendEmail(input, userContext)
					}
					return originalImplementation
				},
			},
		}),
		session.Init(nil),
		dashboard.Init(nil),
		thirdparty.Init(nil),
	},
}
