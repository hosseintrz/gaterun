package main

import (
	"os"

	"github.com/hosseintrz/gaterun/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Could not run command")
	}
}

// func main() {
// 	port := flag.Int("p", 7800, "Port of the service")
// 	configFile := flag.String("c", "./config.json", "Path to the configuration filename")
// 	flag.Parse()

// 	serviceConfig, err := config.Parse(*configFile)
// 	if err != nil {
// 		log.Fatal("ERROR:", err.Error())
// 	}

// 	if *port != 0 {
// 		serviceConfig.Port = *port
// 	}

// 	routerFactory := gorilla.DefaultFactory(proxy.NewDefaultFactory())
// 	routerFactory.New().Run(serviceConfig)

// 	// logger := slog.NewLogLogger(, slog.LevelInfo)

// 	//routerFactory := gin.DefaultFactory(proxy.DefaultFactory(logger), logger)
// 	//routerFactory := gorilla.DefaultFactory(proxy.DefaultFactory(logger), logger)

// 	//running api-gateway
// 	//go routerFactory.New().Run(serviceConfig)

// 	// setting up backend service
// 	// r := mux.NewRouter()
// 	// initRoutes(r)

// 	// srv := &http.Server{
// 	// 	Handler: r,
// 	// 	Addr:    "127.0.0.1:9000",
// 	// 	// Good practice: enforce timeouts for servers you create!
// 	// 	WriteTimeout: 15 * time.Second,
// 	// 	ReadTimeout:  15 * time.Second,
// 	// }

// 	// slog.Info("starting server on 127.0.0.1:9000")
// 	// log.Fatal(srv.ListenAndServe())
// }

// func initRoutes(r *mux.Router) {
// 	r.HandleFunc("/users/{id}/permissions", userPermissionsHandler).Methods("GET")
// 	r.HandleFunc("/users/{id}", usersHandler).Methods("GET")
// }

// func userPermissionsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "users permissions handler")
// }

// func usersHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "users handler")
// }
