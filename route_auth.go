package main

import (
	"chitchat/data"
	"net/http"
)

func login(writer http.ResponseWriter, request *http.Request) {
  generateHTML(writer, nil, "login.layout", "public.navbar", "login")
}

func signup(writer http.ResponseWriter, request *http.Request) {
  generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

func signupAccount(writer http.ResponseWriter, request *http.Request){
  err := request.ParseForm()
  if err != nil {
    danger(err, "Cannot parse form")
  }
  user := data.User{
    Name: request.PostFormValue("name"),
    Email: request.PostFormValue("email"),
    Password: request.PostFormValue("password"),
  }
  if err := user.Create(); err != nil {
    danger(err, "Cannot create user")
  }
  http.Redirect(writer, request, "/login", http.StatusFound)
}

func authenticate(writer http.ResponseWriter, request *http.Request) {  
  
  
  user, err := data.UserByEmail(request.PostFormValue("email"))
  if err != nil {
    danger(err, "Cannot find user")
  }
  if user.Password == data.Encrypt(request.PostFormValue("password")) {
    session, err := user.CreateSession()
    if err != nil {
      danger(err, "Cannot create session")
    }
    cookie := http.Cookie{
      Name:      "_cookie", 
      Value:     session.Uuid,
      HttpOnly:  true,
    }
    http.SetCookie(writer, &cookie)
    http.Redirect(writer, request, "/", http.StatusFound)
  } else {
    http.Redirect(writer, request, "/login", http.StatusFound)
  }
  
}

func logout(writer http.ResponseWriter, request *http.Request) {
  cookie, err := request.Cookie("_cookie")
  if err != http.ErrNoCookie {
    warning(err, "Failed to get cookie")
    session := data.Session{Uuid: cookie.Value}
    session.DeleteByUUID()
  }
  http.Redirect(writer, request, "/", http.StatusFound)
}