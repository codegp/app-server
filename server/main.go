package main

import (
	"net/http"
	"os"

	"github.com/codegp/cloud-persister"
	"github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/codegp/env"
	"github.com/codegp/kube-client"
	"golang.org/x/oauth2"
)

var (
	cp           *cloudpersister.CloudPersister
	kc           *kubeclient.KubeClient
	isLocal      bool
	oAuthConfig  *oauth2.Config
	sessionStore sessions.Store
	logger       *logrus.Logger
)

func init() {
	var err error
	logger = logrus.New()
	logger.Out = os.Stderr
	logger.Infof("ENV %s %v %s", env.GCloudProjectID(), env.IsLocal(), os.Getenv("DATASTORE_EMULATOR_HOST"))
	cp, err = cloudpersister.NewCloudPersister()
	if err != nil {
		logger.Fatalf("Failed to start cloud persister: %v", err)
	}

	kc, err = kubeclient.NewClient()
	if err != nil {
		logger.Fatalf("Failed to start kube client: %v", err)
	}

	if !env.IsLocal() {
		oAuthConfig = configureOAuthClient()
		// TODO get better secret
		cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
		cookieStore.Options = &sessions.Options{
			HttpOnly: true,
		}
		sessionStore = cookieStore
	}
}

func main() {
	baseRouter := mux.NewRouter()
	apiRouter := baseRouter.PathPrefix("/console").Subrouter()

	apiRouter.HandleFunc("/user", sessionMiddleware(GetUser)).Methods("GET")

	apiRouter.HandleFunc("/project/{projectID}", sessionMiddleware(GetProject)).Methods("GET")
	apiRouter.HandleFunc("/project", sessionMiddleware(PostProject)).Methods("POST")
	apiRouter.HandleFunc("/projects/", sessionMiddleware(GetProjects)).Methods("GET")

	apiRouter.HandleFunc("/gametype/{gameTypeID}", sessionMiddleware(GetGameType)).Methods("GET")
	apiRouter.HandleFunc("/gametype", sessionMiddleware(PostGameType)).Methods("POST")
	apiRouter.HandleFunc("/gametype/{gameTypeID}/code", sessionMiddleware(PostGameTypeCode)).Methods("POST")

	apiRouter.HandleFunc("/gametypes/", sessionMiddleware(GetGameTypes)).Methods("GET")

	apiRouter.HandleFunc("/gametype/{gameTypeID}/map/{mapName}", sessionMiddleware(PostMap)).Methods("POST")

	apiRouter.HandleFunc("/game/{gameID}", sessionMiddleware(GetGame)).Methods("GET")
	apiRouter.HandleFunc("/project/{projectID}/map/{mapID}/game", sessionMiddleware(PostGame)).Methods("POST")
	apiRouter.HandleFunc("/game", sessionMiddleware(UpdateGame)).Methods("PUT")
	apiRouter.HandleFunc("/project/{projectID}/games/", sessionMiddleware(GetGames)).Methods("GET")

	apiRouter.HandleFunc("/project/{projectID}/file/{name}", sessionMiddleware(GetProjectFile)).Methods("GET")
	apiRouter.HandleFunc("/project/{projectID}/file/{name}", sessionMiddleware(PostProjectFile)).Methods("POST")
	apiRouter.HandleFunc("/project/{projectID}/file/{name}", sessionMiddleware(PutProjectFile)).Methods("PUT")

	apiRouter.HandleFunc("/types/", sessionMiddleware(GetTypes)).Methods("GET")

	// Get and add types
	apiRouter.HandleFunc("/botType/{botTypeID}", sessionMiddleware(GetBotType)).Methods("GET")
	apiRouter.HandleFunc("/botType", sessionMiddleware(PostBotType)).Methods("POST")

	apiRouter.HandleFunc("/attackType/{attackTypeID}", sessionMiddleware(GetAttackType)).Methods("GET")
	apiRouter.HandleFunc("/attackType", sessionMiddleware(PostAttackType)).Methods("POST")

	apiRouter.HandleFunc("/moveType/{moveTypeID}", sessionMiddleware(GetMoveType)).Methods("GET")
	apiRouter.HandleFunc("/moveType", sessionMiddleware(PostMoveType)).Methods("POST")

	apiRouter.HandleFunc("/terrainType/{terrainTypeID}", sessionMiddleware(GetTerrainType)).Methods("GET")
	apiRouter.HandleFunc("/terrainType", sessionMiddleware(PostTerrainType)).Methods("POST")

	apiRouter.HandleFunc("/itemType/{itemTypeID}", sessionMiddleware(GetItemType)).Methods("GET")
	apiRouter.HandleFunc("/itemType", sessionMiddleware(PostItemType)).Methods("POST")

	apiRouter.HandleFunc("/type/{typeID}/icon", sessionMiddleware(PostIcon)).Methods("POST")

	apiRouter.HandleFunc("/game/{gameID}/history", sessionMiddleware(GetHistory)).Methods("GET")

	// The following handlers are defined in auth.go and used in the
	// "Authenticating Users" part of the Getting Started guide.
	baseRouter.HandleFunc("/login", errorMiddleware(loginHandler)).Methods("GET")
	baseRouter.HandleFunc("/logout", errorMiddleware(logoutHandler)).Methods("POST")
	baseRouter.HandleFunc("/oauth2callback", errorMiddleware(oauthCallbackHandler)).Methods("GET")
	// baseRouter.HandleFunc("/{rest:.*}", ServeDistFile)

	http.Handle("/", &AppServer{baseRouter})

	logger.Info("Serving base router on port 8080")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}

type AppServer struct {
	r *mux.Router
}

func (gc *AppServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger.Infof("Executing request %v", req.URL.Path)
	if origin := req.Header.Get("Origin"); origin != "" && env.IsLocal() {
		logger.Info("Running on localhost, allowing CORs")

		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	gc.r.ServeHTTP(rw, req)
}

// func ServeDistFile(w http.ResponseWriter, r *http.Request) {
// 	path := fmt.Sprintf("../client/dist%v", r.URL.Path)
// 	if info, err := os.Stat(path); err != nil || info.IsDir() {
// 		http.ServeFile(w, r, "../client/dist/index.html")
// 		return
// 	}
//
// 	http.ServeFile(w, r, path)
// }
