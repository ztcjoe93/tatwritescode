package utilities

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMapByKeyReverse(t *testing.T) {

	var hm map[string]map[string]int = map[string]map[string]int{
		"2017": {},
		"2004": {},
		"2022": {},
		"2013": {},
		"2008": {},
		"2012": {},
		"2020": {},
	}

	var keys []string = SortMapByKeyReverse(hm)

	var expectedKeys []string = []string{
		"2022", "2020", "2017", "2013", "2012", "2008", "2004",
	}

	assert.Equal(t, expectedKeys, keys)
}

func TestSortMapByValueReverse(t *testing.T) {
	var hm map[string]int = map[string]int{
		"February": 2,
		"October":  10,
		"December": 12,
		"August":   8,
		"March":    3,
		"November": 11,
	}

	var keys []string = SortMapByValueReverse(hm)

	var expectedKeys []string = []string{
		"December", "November", "October", "August", "March", "February",
	}

	assert.Equal(t, expectedKeys, keys)
}

func TestConvertMonthToIntRepr(t *testing.T) {
	var months []string = []string{
		"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December",
	}

	for i := 0; i < len(months); i++ {
		var monthIntRepr int = ConvertMonthToIntRepr(months[i])
		assert.EqualValues(t, i+1, monthIntRepr)
	}
}

func TestRenderAsHTML(t *testing.T) {
	var htmlString string = "test<br><b>someotherline</b><br><br>test542"
	var htmlRendered template.HTML = RenderAsHTML(htmlString)

	assert.EqualValues(t, htmlRendered, htmlString)
}

func TestHashPasswordAndCheckHash(t *testing.T) {
	var rawPassword string = "passw0rd123"
	hash, _ := HashPassword(rawPassword)

	var matches bool = CheckPasswordHash(rawPassword, hash)
	var notMatches bool = CheckPasswordHash("passw0rd1234", hash)
	assert.True(t, matches)
	assert.False(t, notMatches)
}
