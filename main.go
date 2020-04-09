package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"ohmytech.io/platform/config"
	"ohmytech.io/platform/controllers"
	"ohmytech.io/platform/middlewares"
)

var (
	domain    string
	env       string
	debugMode bool
)

// init is invoked before main()
func init() {
	conf := config.New()
	domain = conf.Domain.Host
	env = conf.Env
	debugMode = conf.DebugMode
}

func main() {
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	fmt.Printf("Domain >> %s\n", "api."+domain)

	// Main router
	r := mux.NewRouter()

	home := controllers.HomeController{}

	// For all
	r.HandleFunc("/", home.Home).Methods(http.MethodGet, http.MethodHead, http.MethodOptions)
	r.HandleFunc("/healthz", home.Healtz).Methods(http.MethodGet)
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/{name}.{extension}").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/.well-known/dnt-policy.txt", home.DNTpolicy).Methods(http.MethodGet)

	if debugMode {
		r.Use(middlewares.Logger)
	}

	// Authentication API
	authenticationRoutes := r.Host("api." + domain).Subrouter()

	authentication := controllers.AuthenticationController{}

	authenticationRoutes.HandleFunc("/signin", authentication.SignIn).Methods(http.MethodPost, http.MethodOptions)
	authenticationRoutes.HandleFunc("/signup", authentication.SignUp).Methods(http.MethodPost, http.MethodOptions)
	authenticationRoutes.HandleFunc("/forgot_password", authentication.ForgotPassword).Methods(http.MethodPost, http.MethodOptions)
	authenticationRoutes.HandleFunc("/reset_password", authentication.ResetPassword).Methods(http.MethodPost, http.MethodOptions)
	authenticationRoutes.HandleFunc("/login_check", authentication.LogoutCheck).Methods(http.MethodGet, http.MethodOptions)
	authenticationRoutes.HandleFunc("/logout", authentication.Logout).Methods(http.MethodPost, http.MethodOptions)

	authenticationRoutes.Use(middlewares.Cors)

	// Secured API
	apiRoutes := r.Host("api." + domain).PathPrefix("/api/v1").Subrouter()

	country := controllers.CountryController{}

	apiRoutes.HandleFunc("/countries", country.List).Methods(http.MethodGet, http.MethodOptions)

	user := controllers.UserController{}

	apiRoutes.HandleFunc("/users", user.List).Methods(http.MethodGet, http.MethodOptions)
	apiRoutes.HandleFunc("/user", user.Save).Methods(http.MethodPost, http.MethodOptions)
	apiRoutes.HandleFunc("/profil", user.Profil).Methods(http.MethodGet, http.MethodOptions)
	apiRoutes.HandleFunc("/user/{id:[0-9]+}", user.Edit).Methods(http.MethodGet, http.MethodOptions)
	apiRoutes.HandleFunc("/user/{id:[0-9]+}", user.Update).Methods(http.MethodPut, http.MethodOptions)

	company := controllers.CompanyController{}

	apiRoutes.HandleFunc("/companies", company.List).Methods(http.MethodGet, http.MethodOptions)
	apiRoutes.HandleFunc("/company", company.Save).Methods(http.MethodPost, http.MethodOptions)
	apiRoutes.HandleFunc("/company/{id:[0-9]+}", company.Edit).Methods(http.MethodGet, http.MethodOptions)
	apiRoutes.HandleFunc("/company/{id:[0-9]+}", company.Update).Methods(http.MethodPut, http.MethodOptions)

	apiRoutes.Use(middlewares.Cors)
	apiRoutes.Use(middlewares.JwtAuthentication)

	r.NotFoundHandler = http.HandlerFunc(home.NotFound)

	srv := &http.Server{
		Addr: "0.0.0.0:8000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20,
		Handler:        r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); nil != err {
			// lw.ApplicationLogger(lw.ApplicationEventLogger{
			// 	Package:  "main",
			// 	Function: "main",
			// }, "ListenAndServe error : "+err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	// lw.ApplicationLogger(lw.ApplicationEventLogger{
	// 	Package:  "main",
	// 	Function: "main",
	// }, "ListenAndServe shutting down")

	os.Exit(0)
}
