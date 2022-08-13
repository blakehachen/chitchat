package main

import (
	"net/http"
	"time"
)

func main() {
	p("ChitChat", version(), "started at", config.Address)
	
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	
	/* All routes START here */

	mux.HandleFunc("/", index)
	mux.HandleFunc("/err", err)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/like", likeThread)
	mux.HandleFunc("/thread/likepost", likePost)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	mux.HandleFunc("/user", getDashboard)
	mux.HandleFunc("/user/posts", getUserPosts)
	mux.HandleFunc("/user/likedPosts", getUserLikedPosts)
	mux.HandleFunc("/user/viewPost", viewPost)
	mux.HandleFunc("/user/threads", getUserThreads)
	mux.HandleFunc("/user/likedThreads", getUserLikedThreads)
	
	/* All routes END here */
	
	server := &http.Server{
		Addr: config.Address,
		Handler: mux,
		ReadTimeout: time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout: time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	
	server.ListenAndServe()
}

