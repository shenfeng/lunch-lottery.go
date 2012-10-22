package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

var people = []string{
	"江宏", "小敏", "小江", "赵雪",
	"吴哲", "名丽", "孙宁", "山川", "朱增",
	"王臣汉", "丽辉", "邹剑", "李蠡",
	"吴江程", "俊文", "朝中", "彦民",
	"杨彤", "王斌", "姜汉", "倪华杰", "晓丹",
	"边边", "唐宇", "钱国祥", "莫倩", "沈锋",
}

var groupCount = 6

var addr = flag.String("addr", ":9292", "The addr to listen (':9292')")

type Lunch struct {
	Year int
	Month time.Month
	Day   int
	Group [][]string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func listAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List")
}

func shuffle(input []string) []string {
	ret := make([]string, len(input))
	idxs := rand.Perm(len(input))
	for i := 0; i < len(input); i++ {
		ret[i] = input[idxs[i]]
	}
	return ret
}

func group(people []string, count int) [][]string {
	ret := make([][]string, count)
	perGroup := len(people) / count
	less := count - len(people)%count
	for i := 0; i < count; i++ {
		if i < less {
			ret[i] = make([]string, perGroup)
		} else {
			ret[i] = make([]string, perGroup+1)
		}
	}

	idx := 0
	for _, arr := range ret {
		for i := 0; i < len(arr); i++ {
			arr[i] = people[idx]
			idx += 1
		}
	}
	return ret
}

func newLunch(w http.ResponseWriter, r *http.Request) {
	g := group(shuffle(people), groupCount)
	t, _ := template.ParseFiles("show.html")
	y, m, d := time.Now().Date()
	// fmt.Println(time.Now())
	// fmt.Println(y, m, d)
	t.Execute(w, &Lunch{Year: y,Month: m, Day: d, Group: g})
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/list", listAll)
	http.HandleFunc("/new", newLunch)
	http.ListenAndServe(*addr, nil)
}
