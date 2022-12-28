package main

import (
	"html/template"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

// sorts a map by its key in lexicographically reverse order
func SortMapByKeyReverse(hm map[string]map[string]int) []string {

	keys := make([]string, len(hm))

	count := 0
	for k := range hm {
		keys[count] = k
		count++
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	return keys
}

// sorts a map by its value in month-reverse order
func SortMapByValueReverse(hm map[string]int) []string {

	keys := make([]string, len(hm))

	count := 0
	for k := range hm {
		keys[count] = k
		count++
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return hm[keys[i]] > hm[keys[j]]
	})

	return keys
}

func ConvertMonthToIntRepr(month string) int {

	hm := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}

	return hm[month]
}

func RenderAsHTML(htmlString string) template.HTML {
	return template.HTML(htmlString)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
