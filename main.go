package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/welschma/gowebdev/controllers"
	"github.com/welschma/gowebdev/migrations"
	"github.com/welschma/gowebdev/models"
	"github.com/welschma/gowebdev/templates"
	"github.com/welschma/gowebdev/views"
)

func RequestLoggerMw(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming %v request for %v from IP adress: %v", r.Method, r.URL, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}

func main() {

	//Set up the database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	//Set up services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	//Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)

	//Set up controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	//Set up router and routes
	r := chi.NewRouter()
    
    r.Use(RequestLoggerMw)
    r.Use(csrfMw)
    r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)

	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	r.Post("/signout", usersC.ProcessSignOut)

    r.Route("/users/me", func(r chi.Router) {
        r.Use(umw.RequireUser)
        r.Get("/", usersC.CurrentUser)
    })


    r.NotFound(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Page not found", http.StatusNotFound)
    })

    //Start the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
