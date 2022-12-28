package posts

import (
	"fmt"
	"strconv"
	"time"

	"twc-app/utilities"
)

type NavigationLinks struct {
	YearLinks []*YearLink
	YearOrder []string
}

type YearLink struct {
	Year       string
	MonthOrder []string
}

type Blogpost struct {
	Id       int
	Datetime string
	Title    string
	Content  string
}

func NewBlogpost(id int, datetime string, title string, content string) *Blogpost {
	post := Blogpost{
		Id:       id,
		Datetime: datetime,
		Title:    title,
		Content:  content,
	}

	return &post
}

func GetNavigationLinks(posts []*Blogpost) *NavigationLinks {

	postMap := GetMmYyyy(posts)
	yearLinks := make([]*YearLink, len(postMap))

	yearOrder := utilities.SortMapByKeyReverse(postMap)

	count := 0
	for _, year := range yearOrder {
		monthOrder := utilities.SortMapByValueReverse(postMap[year])
		yearLink := &YearLink{
			Year:       year,
			MonthOrder: monthOrder,
		}

		yearLinks[count] = yearLink
		count++
	}

	return &NavigationLinks{
		YearLinks: yearLinks,
		YearOrder: yearOrder,
	}
}

func GetMmYyyy(posts []*Blogpost) map[string]map[string]int {

	/*
		To store entries side navigation links
		{
			"YYYY": {
				"MM": <MM_int>,
				...
			}
			...
		}
	*/
	hm := make(map[string]map[string]int)
	links := make([]string, len(posts))

	for idx, post := range posts {
		ts, err := time.Parse("2006-01-02 15:04:05", post.Datetime)
		if err != nil {
			fmt.Printf("Unable to parse time: %v\n", err)
		}

		yearStr := strconv.Itoa(ts.Year())
		monthStr := ts.Month().String()

		if _, ok := hm[yearStr]; !ok {
			hm[yearStr] = make(map[string]int)
		}
		hm[yearStr][monthStr] = int(ts.Month())

		links[idx] = ts.Format(time.RFC822Z)
	}

	return hm
}
