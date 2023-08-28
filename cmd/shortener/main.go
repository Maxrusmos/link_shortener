package main

import (
	"flag"
	config "link_shortener/internal/configs"
	"link_shortener/internal/startprep"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	cookie := &http.Cookie{
	// 		Name:  "auth_cookie",
	// 		Value: "user",
	// 	}
	// 	http.SetCookie(w, cookie)
	// })
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	conf := config.GetConfig()
	r := chi.NewRouter()
	flag.StringVar(&conf.Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "Base address for shortened URL")
	flag.StringVar(&conf.FileStore, "f", "", "File storage")
	flag.StringVar(&conf.DBConnect, "d", "", "db Connection String")
	flag.Parse()
	storageURL := startprep.GetStorageURL(conf)
	startprep.RegisterRoutes(r, storageURL, conf.BaseURL, logger)
	startprep.StartServer(conf.Address, r, logger)
}
