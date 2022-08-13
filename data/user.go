package data

import (
	"time"
)

type User struct {
	Id				int
	Uuid			string
	Name			string
	Email			string
	Password	string
	CreatedAt time.Time
}

type Session struct {
	Id				int
	Uuid			string
	Email			string
	UserId		int
	CreatedAt time.Time
}

func (user *User) CreateSession() (session Session, err error) {
  db := db()
  defer db.Close()
  statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
  stmt, err := db.Prepare(statement)
  if err != nil {
    return
  }
  defer stmt.Close()
  
	// use QueryRow to return a row and scan the returned id into the Session struct
  err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
  if err != nil {
    return
  }    
  return
}

func (user *User) Session (session Session, err error) {
	db := db()
	defer db.Close()
	session = Session{}
	err = db.QueryRow("SELECT id, uuid, email, created_at FROM sesssions WHERE user_id = $1", user.Id).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
	
}

func (session *Session) Check() (valid bool, err error) {
	db := db()
	defer db.Close()
	err = db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)

	if err != nil {
		valid = false
		return
	}

	if session.Id != 0 {
		valid = true
	}
	return
}

func (session *Session) DeleteByUUID() (err error) {
	db := db()
	defer db.Close()
	statement := "DELETE FROM sessions WHERE uuid = $1"
	stmt, err := db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	if err != nil { 
		return
	}
	return
}

func (session *Session) User() (user User, err error) {
	db := db()
	defer db.Close()
	user = User{}
	err = db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func SessionDeleteAll() {
	db := db()
	defer db.Close()
	statement := "DELETE FROM sessions"
	_, err := db.Exec(statement)
	if err != nil {
		return
	}
	
}

func (user *User) Create() (err error) {
	db := db()
	defer db.Close()

	statement := "INSERT INTO users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (user *User) Delete() (err error) {
	db := db()
	defer db.Close()
	statement := "DELETE FROM users WHERE id = $1"
	stmt, err := db.Prepare(statement)
	if err != nil { 
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		return
	}
	return
}

func (user *User) Update() (err error) {
	db := db()
	defer db.Close()

	statement := "UPDATE users SET name = $2, email = $3 WHERE id = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	if err != nil {
		return
	}
	return
}

func UserDeleteAll() (err error) {
	db := db()
	defer db.Close()
	statement := "DELETE FROM users"

	_, err = db.Exec(statement)
	if err != nil {
		return
	}
	return
}

func Users() (users []User, err error) {
	db := db()
	defer db.Close()
	rows, err := db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}

		users = append(users, user)
	}
	rows.Close()
	return
}

func UserByEmail(email string) (user User, err error) {
	db := db()
	defer db.Close()
	user = User{}
	err = db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return
	}
	return
}

func UserByUUID(uuid string) (user User, err error) {
	db := db()
	defer db.Close()
	user = User{}
	err = db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return
	}
	return
}

