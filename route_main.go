package main

import (
	"chitchat/data"
	"net/http"
)

func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")

	}else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	threads, err := data.Threads(); if err != nil {
		
		error_message(writer, request, "Cannot get threads")

	}else {
		
		session, err := session(writer, request)
		if err != nil {
			data1 := struct {
				Threads []data.Thread
				CurrentUser    string
				Uuid    string
			}{threads, "", ""}
			generateHTML(writer, data1, "layout", "public.navbar", "index")
		} else {
			
			user, err := data.UserByEmail(session.Email)
			if err != nil {
				error_message(writer, request, "User not found.")
			}else{
				data2 := struct {
					Threads 			[]data.Thread
					CurrentUser   string
					UserUUID			string
				}{threads, user.Name, user.Uuid}
				
				generateHTML(writer, data2, "layout", "private.navbar", "index")
			}
			
		}
	}
}