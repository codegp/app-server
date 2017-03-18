package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"golang.org/x/net/context"

	"github.com/satori/go.uuid"
)

const githubProfileSessionKey = "github_user"

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&github.User{})
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) *requestError {
	sessionID := uuid.NewV4().String()

	oauthFlowSession, err := sessionStore.New(r, sessionID)
	log.Printf("CRETATE SESS %v", oauthFlowSession)
	if err != nil {
		return requestErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	// redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	// if err != nil {
	// 	return requestErrorf(err, "invalid redirect URL: %v", err)
	// }
	oauthFlowSession.Values[oauthFlowRedirectKey] = "http://localhost:3000/"

	if err := oauthFlowSession.Save(r, w); err != nil {
		return requestErrorf(err, "could not save session: %v", err)
	}

	url := githubOAuthConfig.AuthCodeURL(sessionID, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) *requestError {
	oauthFlowSession, err := sessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return requestErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	log.Println("redirectURL", redirectURL)
	// Validate this callback request came from the app.
	if !ok {
		return requestErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return requestErrorf(err, "could not get auth token: %v", err)
	}

	session, err := sessionStore.New(r, defaultSessionID)
	if err != nil {
		return requestErrorf(err, "could not get default session: %v", err)
	}

	user, err := fetchGithubUser(context.Background(), token)
	if err != nil {
		return requestErrorf(err, "fetchGithubUser failed: %v", err)
	}

	session.Values[oauthTokenSessionKey] = token
	session.Values[githubProfileSessionKey] = user

	logger.Infof("token %v \n\nsession %v", token, session)
	if err := session.Save(r, w); err != nil {
		return requestErrorf(err, "could not save session: %v", err)
	}

	logger.Infof("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}

func fetchGithubUser(ctx context.Context, tok *oauth2.Token) (*github.User, error) {
	oauthClient := oauth2.NewClient(ctx, githubOAuthConfig.TokenSource(ctx, tok))
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func configureGithubOAuthClient() *oauth2.Config {
	redirectURL := os.Getenv("GITHUB_OAUTH2_CALLBACK")
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	if redirectURL == "" || clientID == "" || clientSecret == "" {
		logger.Fatal("OAuth2 environment variables not found!")
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"user:email"},
		Endpoint: githuboauth.Endpoint,
	}
}

// profileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func githubUserFromSession(r *http.Request) *github.User {
	session, err := sessionStore.Get(r, defaultSessionID)
	log.Printf("sess %v %v", session, err)
	if err != nil {
		return nil
	}
	tok, ok := session.Values[oauthTokenSessionKey].(*oauth2.Token)
	log.Printf("tok %v %v", tok, ok)
	if !ok || !tok.Valid() {
		return nil
	}

	profile, ok := session.Values[githubProfileSessionKey].(*github.User)
	log.Printf("profile %v %v", profile, ok)
	if !ok {
		return nil
	}
	return profile
}
