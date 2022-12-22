package database

import (
	"database/sql"
	"fmt"
	"internal/posts"
	"strconv"
)

func OpenSqlConnection(user string, password string, database string, host string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, database))
	if err != nil {
		fmt.Printf("Error connecting to db: %v\n", err)
	}

	return db
}

func GetAllBlogposts(db *sql.DB) []*posts.Blogpost {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Printf("some error -> %v\n", err)
	}
	defer rows.Close()

	blogPosts := make([]*posts.Blogpost, 0)

	for rows.Next() {
		var id int
		var datetime string
		var title string
		var content string

		err := rows.Scan(&id, &datetime, &title, &content)
		if err != nil {
			fmt.Printf("Some error -> %v\n", err)
		}

		blogPosts = append(blogPosts, posts.NewBlogpost(id, datetime, title, content))
	}

	return blogPosts
}

func GetLatestPosts(db *sql.DB) []*posts.Blogpost {
	rows, err := db.Query("SELECT * FROM posts ORDER BY datetimestamp DESC LIMIT 5")
	if err != nil {
		fmt.Printf("some error -> %v\n", err)
	}
	defer rows.Close()
	blogPosts := populateRows(rows)

	return blogPosts

}

func GetPostsFromMonth(db *sql.DB, year string, month string) ([]*posts.Blogpost, error) {

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		fmt.Printf("Error converting year to integer -- %v\n", err)
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		fmt.Printf("Error converting month to integer -- %v\n", err)
	}

	var ulYear, ulMonth int
	if monthInt == 12 {
		ulYear = yearInt + 1
		ulMonth = 1
	} else {
		ulYear = yearInt
		ulMonth = monthInt + 1
	}

	fromDate := year + "-" + fmt.Sprintf("%02d", monthInt) + "-" + "01"
	toDate := strconv.Itoa(ulYear) + "-" + fmt.Sprintf("%02d", ulMonth) + "-" + "01"

	rows, err := db.Query("SELECT * FROM posts WHERE datetimestamp >= ? AND datetimestamp < ?", fromDate, toDate)
	if err != nil {
		fmt.Printf("Some error -> %v\n", err)
	}
	defer rows.Close()
	blogPosts := populateRows(rows)

	return blogPosts, nil
}

func populateRows(rows *sql.Rows) []*posts.Blogpost {
	blogPosts := make([]*posts.Blogpost, 0)

	for rows.Next() {
		var id int
		var datetime string
		var title string
		var content string

		err := rows.Scan(&id, &datetime, &title, &content)
		if err != nil {
			fmt.Printf("Some error -> %v\n", err)
		}

		blogPosts = append(blogPosts, posts.NewBlogpost(id, datetime, title, content))
	}

	return blogPosts
}
