package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewBlogpost(t *testing.T) {
	var id int = 1
	var datetime string = "2022-10-10 05:00:25"
	var title string = "test_title"
	var content string = "content<br>somestuffhere<br>"
	var post *Blogpost = NewBlogpost(id, datetime, title, content)

	assert.Equal(t, post.Id, id)
	assert.Equal(t, post.Datetime, datetime)
	assert.Equal(t, post.Title, title)
	assert.Equal(t, post.Content, content)
}

func TestGetMmYyyy(t *testing.T) {
	blogposts := createBlogposts()

	sideNavEntries := GetMmYyyy(blogposts)

	expectedNavEntries := map[string]map[string]int{
		"2018": {"January": 1},
		"2021": {"December": 12},
		"2022": {"February": 2, "October": 10},
	}

	assert.True(t, cmp.Equal(sideNavEntries, expectedNavEntries))
}

func TestGetNavigationLinks(t *testing.T) {
	blogposts := createBlogposts()

	navigationLinks := GetNavigationLinks(blogposts)

	expectedNaviLink := NavigationLinks{
		YearLinks: []*YearLink{
			{
				Year:       "2022",
				MonthOrder: []string{"October", "February"},
			},
			{
				Year:       "2021",
				MonthOrder: []string{"December"},
			},
			{
				Year:       "2018",
				MonthOrder: []string{"January"},
			},
		},
		YearOrder: []string{"2022", "2021", "2018"},
	}

	assert.True(t, reflect.DeepEqual(navigationLinks.YearLinks, expectedNaviLink.YearLinks))
	assert.True(t, reflect.DeepEqual(navigationLinks.YearOrder, expectedNaviLink.YearOrder))
}

func createBlogposts() []*Blogpost {
	datetimes := []string{
		"2022-10-10 05:00:25", "2021-12-05 12:00:25", "2018-01-01 03:41:41", "2022-02-10 05:00:25",
	}

	blogposts := make([]*Blogpost, len(datetimes))

	for i := 1; i < 5; i++ {
		blogposts[i-1] = NewBlogpost(i, datetimes[i-1], "test_title", "test_content")
	}

	return blogposts
}
