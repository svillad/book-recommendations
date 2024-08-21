package api

import (
	"net/http"

	"github.com/book-recommendations/service/config"
	"github.com/book-recommendations/service/controllers"
	"github.com/book-recommendations/service/mediators"
	"github.com/book-recommendations/service/stores"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Routes prepares the mux router to be served
func Routes() http.Handler {
	// initialize controllers
	bookController, authorController, genreController, sizeController, eraController := generateControllers()

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// routes
	router.HandleFunc("/books", bookController.Get).Methods(http.MethodGet)
	router.HandleFunc("/authors", authorController.Get).Methods(http.MethodGet)
	router.HandleFunc("/genres", genreController.Get).Methods(http.MethodGet)
	router.HandleFunc("/sizes", sizeController.Get).Methods(http.MethodGet)
	router.HandleFunc("/eras", eraController.Get).Methods(http.MethodGet)

	return cors.AllowAll().Handler(router)
}

// generateControllers constructs the needed controller with dependency injected mediators
func generateControllers() (controllers.BookController, controllers.AuthorController, controllers.GenreController, controllers.SizeController, controllers.EraController) {
	configValues, err := config.LoadConfig()
	if err != nil {
		log.WithField("error", err).Error("error init config application")
	}

	storeAdapter, err := stores.NewStore(configValues.DatabaseURL)
	if err != nil {
		log.WithField("error", err).Error("error initializing database")
	}

	// ------------------------ book ------------------------
	bookMediatorFactory := func() mediators.BookMediator {
		storeLog := log.WithField("*store", "Book")
		bookStore := stores.NewBookStore(storeLog, storeAdapter.GetDB())
		mediatorLog := log.WithField("*mediator", "Book")
		return mediators.NewBookMediator(mediatorLog, bookStore)
	}
	bookController := controllers.BookController{
		Logger:              log.WithField("*controller", "Book"),
		BookMediatorFactory: bookMediatorFactory,
	}

	// ------------------------ author ------------------------
	authorMediatorFactory := func() mediators.AuthorMediator {
		storeLog := log.WithField("*store", "Author")
		authorStore := stores.NewAuthorStore(storeLog, storeAdapter.GetDB())
		mediatorLog := log.WithField("*mediator", "Author")
		return mediators.NewAuthorMediator(mediatorLog, authorStore)
	}
	authorController := controllers.AuthorController{
		Logger:                log.WithField("*controller", "Author"),
		AuthorMediatorFactory: authorMediatorFactory,
	}

	// ------------------------ genre ------------------------
	genrerMediatorFactory := func() mediators.GenreMediator {
		storeLog := log.WithField("*store", "Genre")
		genreStore := stores.NewGenreStore(storeLog, storeAdapter.GetDB())
		mediatorLog := log.WithField("*mediator", "Genre")
		return mediators.NewGenreMediator(mediatorLog, genreStore)
	}
	genrerController := controllers.GenreController{
		Logger:               log.WithField("*controller", "Genre"),
		GenreMediatorFactory: genrerMediatorFactory,
	}

	// ------------------------ size ------------------------
	sizeMediatorFactory := func() mediators.SizeMediator {
		storeLog := log.WithField("*store", "Size")
		sizeStore := stores.NewSizeStore(storeLog, storeAdapter.GetDB())
		mediatorLog := log.WithField("*mediator", "Size")
		return mediators.NewSizeMediator(mediatorLog, sizeStore)
	}
	sizeController := controllers.SizeController{
		Logger:              log.WithField("*controller", "Size"),
		SizeMediatorFactory: sizeMediatorFactory,
	}

	// ------------------------ era ------------------------
	eraMediatorFactory := func() mediators.EraMediator {
		storeLog := log.WithField("*store", "Era")
		eraStore := stores.NewEraStore(storeLog, storeAdapter.GetDB())
		mediatorLog := log.WithField("*mediator", "Era")
		return mediators.NewEraMediator(mediatorLog, eraStore)
	}
	eraController := controllers.EraController{
		Logger:             log.WithField("*controller", "Era"),
		EraMediatorFactory: eraMediatorFactory,
	}

	return bookController, authorController, genrerController, sizeController, eraController
}
