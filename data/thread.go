package data

import (
	"time"
)

type Thread struct {
	Id					int
	Uuid				string
	Topic				string
	UserId			int
	Author			string
	CreatedAt		time.Time
	Likes				int
}

type Post struct {
	Id					int
	Uuid				string
	Body				string
	UserId			int
	ThreadId		int
	CreatedAt		time.Time
	Likes				int
}

func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")

}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (thread *Thread) NumReplies() (count int) {
	db := db()
	defer db.Close()
	rows, err := db.Query("SELECT count(*) FROM posts WHERE thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
	return
}

func (thread *Thread) Posts() (posts []Post, err error) {
	db := db()
	defer db.Close()
	rows, err := db.Query("SELECT id, uuid, body, user_id, thread_id, created_at, likes FROM posts WHERE thread_id = $1 ORDER BY created_at DESC", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt, &post.Likes); err != nil {
			return
		}

		posts = append(posts, post)
	}
	rows.Close()
	return
}

func (user *User) Posts() (posts []Post, err error) {
	db := db()
	defer db.Close()
	rows, err := db.Query("SELECT id, uuid, body, user_id, thread_id, created_at, likes FROM posts WHERE user_id = $1 ORDER BY created_at DESC", user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}

		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt, &post.Likes); err != nil {
			return
		}

		posts = append(posts, post)
	}
	rows.Close()
	return
}

func (user *User) CreateThread(topic string) (conv Thread, err error) {
	db := db()
	defer db.Close()
	statement1 := "INSERT INTO threads (uuid, topic, user_id, author, created_at, likes) values ($1, $2, $3, $4, $5, 0) returning id, uuid, topic, user_id, author, created_at, likes"
	stmt, err := db.Prepare(statement1)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), topic, user.Id, user.Email, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.Author, &conv.CreatedAt, &conv.Likes)
	
	if err != nil {
		return
	}
	
	return
}

func (user *User) LikeThread(thread string) (count int, err error) {
	db := db()
	defer db.Close()
	
	//Search for user and thread in likedThreads
	rows, err := db.Query("SELECT count(*) FROM likedThreads WHERE thread_uuid = $1 AND user_uuid = $2", thread, user.Uuid)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()

	//Has the user liked a specific thread?
	if count > 0 {

		/*Yes? Delete users uuid and thread uuid from liked threads and update thread
		update threads at uuid value likes - 1*/
		removeLike := "UPDATE threads SET likes = likes - 1 WHERE uuid = $1"
		_, err = db.Exec(removeLike, thread)
		if err != nil {
			return
		}
		removeFromLikedThreads := "DELETE FROM likedThreads WHERE user_uuid = $1 AND thread_uuid = $2"
		_, err = db.Exec(removeFromLikedThreads, user.Uuid, thread)
		if err != nil {
			return
		}

		return
	}else{

		/*No? Insert users uuid and thread uuid into liked threads and update thread
		update threads at uuid value likes + 1*/
		statement1 := "INSERT INTO likedThreads (user_uuid, thread_uuid) values ($1, $2)"
		_, err = db.Exec(statement1, user.Uuid, thread)
	

		statement2 := "UPDATE threads SET likes = likes + 1 WHERE uuid = $1"
		_, err = db.Exec(statement2, thread)
	}

	
	return
}

func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	db := db()
	defer db.Close()
	statement := "INSERT INTO posts (uuid, body, user_id, thread_id, created_at, likes) values ($1, $2, $3, $4, $5, 0) returning id, uuid, body, user_id, thread_id, created_at, likes"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt, &post.Likes)
	if err != nil {
		return
	}
	return
}

func (user *User) LikePost(post string) (count int, err error) {
	db := db()
	defer db.Close()
	
	//Search for thread, post and user in likedPosts
	err = db.QueryRow("SELECT count(*) FROM likedPosts WHERE (user_uuid = $1 AND post_uuid = $2)", user.Uuid, post).Scan(&count)
	if err != nil {
		return
	}
	
	
	if count > 0 {
		/*Yes? Delete users uuid and thread uuid from liked threads and update thread
		update threads at uuid value likes - 1*/

		removeFromLikedPosts := "DELETE FROM likedPosts WHERE (user_uuid = $1 AND post_uuid = $2)"
		_, err = db.Exec(removeFromLikedPosts, user.Uuid,  post)
		if err != nil {
			return
		}

		removeLike := "UPDATE posts SET likes = likes - 1 WHERE uuid = $1"
		_, err = db.Exec(removeLike, post)
		if err != nil {
			return
		}
		
	}else{
		
		/*No? Insert users uuid and thread uuid into liked threads and update thread
		update threads at uuid value likes + 1*/
		statement1 := "INSERT INTO likedPosts (user_uuid, post_uuid) values ($1, $2)"
		_, err = db.Exec(statement1, user.Uuid, post)
	

		statement2 := "UPDATE posts SET likes = likes + 1 WHERE uuid = $1"
		_, err = db.Exec(statement2, post)
	}

	return
}

func Threads() (threads []Thread, err error) {
	db := db()
	defer db.Close()
	
	rows, err := db.Query("SELECT id, uuid, topic, user_id, author, created_at, likes FROM threads ORDER BY created_at DESC")
	if err != nil {
		
		return
	}
	
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.Author, &conv.CreatedAt, &conv.Likes); err != nil {
			
			return
		}
		
		threads = append(threads, conv)
	}
	rows.Close()
	
	return
}

func (user *User) Threads() (threads []Thread, err error) {
	db := db()
	defer db.Close()

	rows, err := db.Query("SELECT id, uuid, topic, user_id, author, created_at, likes FROM threads WHERE user_id = $1 ORDER BY created_at DESC", user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.Author, &conv.CreatedAt, &conv.Likes); err != nil {
			
			return
		}
		
		threads = append(threads, conv)
	}
	rows.Close()

	return

} 

func ThreadByUUID(uuid string) (conv Thread, err error) {
	db := db()
	defer db.Close()
	conv = Thread{}
	err = db.QueryRow("SELECT id, uuid, topic, user_id, author, created_at, likes FROM threads WHERE uuid = $1", uuid).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.Author, &conv.CreatedAt, &conv.Likes)
	return
}

func ThreadByPostUUID(uuid string) (conv Thread, err error) {
	db := db()
	defer db.Close()
	post := Post{}
	err = db.QueryRow("SELECT id, uuid, body, user_id, thread_id, created_at, likes FROM posts WHERE uuid = $1", uuid).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt, &post.Likes)
	conv = Thread{}
	err = db.QueryRow("SELECT id, uuid, topic, user_id, author, created_at, likes FROM threads WHERE id = $1", post.ThreadId).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.Author, &conv.CreatedAt, &conv.Likes)
	return
}



func (thread *Thread) User() (user User) {
	db := db()
	defer db.Close()
	user = User{}
	db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (post *Post) User() (user User) {
	db := db()
	defer db.Close()
	user = User{}
	db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", post.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (user *User) LikedPosts() (posts []Post, err error){
	db := db()
	defer db.Close()
	var post_uuids []string

	rows, err := db.Query("SELECT post_uuid FROM likedPosts WHERE user_uuid = $1", user.Uuid)
	if err != nil {
		return
	}
	for rows.Next() {
		var uuid string
		if err = rows.Scan(&uuid); err != nil {
			return
		}
		
		post_uuids = append(post_uuids, uuid)
	}
	rows.Close()
	
	for _, item := range post_uuids {
		post := Post{}
		err = db.QueryRow("SELECT id, uuid, body, user_id, thread_id, created_at, likes FROM posts WHERE uuid = $1", item).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt, &post.Likes)
		if err != nil{
			return
		}

		posts = append(posts, post)
	}
	
	return
}
func (user *User) LikedThreads() (threads []Thread, err error){
	db := db()
	defer db.Close()
	var thread_uuids []string

	rows, err := db.Query("SELECT thread_uuid FROM likedThreads WHERE user_uuid = $1", user.Uuid)
	if err != nil {
		return
	}
	for rows.Next() {
		var uuid string
		if err = rows.Scan(&uuid); err != nil {
			return
		}
		
		thread_uuids = append(thread_uuids, uuid)
	}
	rows.Close()
	
	for _, item := range thread_uuids {
		thread := Thread{}
		err = db.QueryRow("SELECT id, uuid, topic, user_id, author, created_at, likes FROM threads WHERE uuid = $1", item).Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.Author, &thread.CreatedAt, &thread.Likes)
		if err != nil{
			return
		}

		threads = append(threads, thread)
	}
	
	return
}

