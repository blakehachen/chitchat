package main

import (
	"chitchat/data"
	"fmt"
	"net/http"
	"text/template"
)

func newThread(writer http.ResponseWriter, request *http.Request) {
	session, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		user, err := data.UserByEmail(session.Email)
			if err != nil {
				error_message(writer, request, "User not found.")
			}else {
				data := struct {
					CurrentUser	string
					UserUUID string 
				}{user.Name, user.Uuid}
				generateHTML(writer, data, "layout", "private.navbar", "new.thread")
			}
	}
}

func createThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}

		user, err := sess.User(); if err != nil {
			danger(err, "Cannot get user from session")
		}

		topic := request.PostFormValue("topic")
		escaped_topic := template.HTMLEscapeString(topic)
		
		if _, err := user.CreateThread(escaped_topic); err != nil {
			danger(err, "Cannot create thread")

		}
		
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func likeThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}

		user, err := sess.User(); if err != nil {
			danger(err, "Cannot get user from session")
		}

		uuid := request.PostFormValue("uuid")
		if _, err := user.LikeThread(uuid); err != nil {
			danger(err, "Cannot create thread")
		}
		
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func likePost(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}

		user, err := sess.User(); if err != nil {
			danger(err, "Cannot get user from session")
		}

		uuid := request.PostFormValue("uuid")
		if _, err := user.LikePost(uuid); err != nil {
			danger(err, "Cannot LIKE post")
		}

		thread, err := data.ThreadByPostUUID(uuid)
		if err != nil {
			danger(err, "Cannot find thread by post uuid")
		}

		url := fmt.Sprint("/thread/read?id=", thread.Uuid)
		http.Redirect(writer, request, url, http.StatusFound)
	}
}

func getDashboard(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	user, err := data.UserByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find user account")
	} else {
		_, err := session(writer, request)
		if err != nil {
			http.Redirect(writer, request, "/login", http.StatusFound)
		}
		data := struct {
			User				data.User
			UserUUID 		string
			CurrentUser string
		}{user, user.Uuid, user.Name}
		
		generateHTML(writer, &data, "layout", "private.navbar", "user", "user.info")

	}

}

func readThread(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		sess, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &thread, "layout", "public.navbar", "public.thread")
		} else {
			user, err := data.UserByEmail(sess.Email)
			if err != nil{
				error_message(writer, request, "Cannot find user")
			}else {
				posts, err := thread.Posts(); if err != nil {
					error_message(writer, request, "Cannot find posts for given thread")
				}

				data1 := struct {
					ThreadUser					data.User
					Thread							data.Thread
					Posts								[]data.Post
					ThreadCreatedAtDate string
					CurrentUser 				string
					UserUUID						string
				} {thread.User(), thread, posts, thread.CreatedAtDate(), user.Name, user.Uuid}
				
				generateHTML(writer, &data1, "layout", "private.navbar", "private.thread")
			}
		}
	}
}

func postThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User(); if err != nil {
			danger(err, "Cannot get user from session")
		}

		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")

		thread, err := data.ThreadByUUID(uuid); if err != nil {
			error_message(writer, request, "Cannot read thread")
		}

		if _, err = user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}

		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(writer, request, url, http.StatusFound)
	}
}

func viewPost(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByPostUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find thread corresponding to post")
	} else {
		url := fmt.Sprint("/thread/read?id=", thread.Uuid)
		http.Redirect(writer, request, url, http.StatusFound)
	}
}

func getUserPosts(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	user, err := data.UserByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find user account")
	} else {
		sess, err := session(writer, request)
		if err != nil && sess.Email != user.Email {
			http.Redirect(writer, request, "/login", http.StatusFound)
		}
		
		posts, err := user.Posts(); if err != nil {
			error_message(writer, request, "Cannot find posts for given user")
		}

		data := struct {
			User				data.User
			Posts				[]data.Post
			UserUUID 		string
			CurrentUser string
		}{user, posts, user.Uuid, user.Name}

		generateHTML(writer, &data, "layout", "private.navbar", "user", "user.posts")
	}
}

func getUserThreads(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	user, err := data.UserByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find user account")
	} else {
		sess, err := session(writer, request)
		if err != nil && sess.Email != user.Email {
			http.Redirect(writer, request, "/login", http.StatusFound)
		}
		
		threads, err := user.Threads(); if err != nil {
			error_message(writer, request, "Cannot find posts for given user")
		}

		data := struct {
			User				data.User
			Threads			[]data.Thread
			UserUUID 		string
			CurrentUser string
		}{user, threads, user.Uuid, user.Name}

		generateHTML(writer, data, "layout", "private.navbar", "user", "user.threads")
	}
}

func getUserLikedPosts(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	user, err := data.UserByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find user account")
	} else {
		sess, err := session(writer, request)
		if err != nil && sess.Email != user.Email {
			http.Redirect(writer, request, "/login", http.StatusFound)
		}
		likedPosts, err := user.LikedPosts(); if err != nil {
			error_message(writer, request, "Cannot find liked posts for user account")
		}
		data := struct {
			User				data.User
			Posts				[]data.Post
			UserUUID 		string
			CurrentUser string
		}{user, likedPosts, user.Uuid, user.Name}

		generateHTML(writer, &data, "layout", "private.navbar", "user", "user.posts")

		
	}
}

func getUserLikedThreads(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	user, err := data.UserByUUID(uuid); if err != nil {
		error_message(writer, request, "Cannot find user account")
	} else {
		sess, err := session(writer, request)
		if err != nil && sess.Email != user.Email {
			http.Redirect(writer, request, "/login", http.StatusFound)
		}
		likedThreads, err := user.LikedThreads(); if err != nil {
			error_message(writer, request, "Cannot find liked posts for user account")
		}
		data := struct {
			User					data.User
			Threads				[]data.Thread
			UserUUID 			string
			CurrentUser 	string
		}{user, likedThreads, user.Uuid, user.Name}

		generateHTML(writer, &data, "layout", "private.navbar", "user", "user.threads")

		
	}
}

