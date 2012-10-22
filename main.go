package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

const seperator = ","
const newline = "\n"

type Lunch struct {
	Group [][]string
	Order string
	Date  string
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")

	files, _ := ioutil.ReadDir("data")
	lunches := make([]Lunch, len(files))
	for i, file := range files {
		bytes, _ := ioutil.ReadFile("data/" + file.Name())
		c := strings.Trim(string(bytes), "\r\n ")
		lines := strings.Split(c, newline)
		groups := make([][]string, len(lines))
		for j, line := range lines {
			groups[j] = strings.Split(line, seperator)
		}
		lunches[i] = Lunch{Group: groups, Date: file.Name()}
	}
	t.Execute(w, lunches)
}

func saveGroupHandler(w http.ResponseWriter, r *http.Request) {
	order := r.FormValue("order")
	ordered := strings.Split(order, seperator)
	now := time.Now()
	if now.Weekday() == time.Thursday {
		grouped := group(ordered, groupCount)
		lines := make([]string, len(grouped))
		for i := 0; i < len(grouped); i++ {
			lines[i] = strings.Join(grouped[i], seperator)
		}
		y, m, d := now.Date()
		name := fmt.Sprintf("data/%d-%d-%d", y, m, d)
		_, err := os.Stat(name)
		if err != nil {
			ioutil.WriteFile(name, []byte(strings.Join(lines, newline)), 0600)
			log.Print("write file ", name)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/s/error.html", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/s/error.html", http.StatusFound)
	}
}

func newLunchHanler(w http.ResponseWriter, r *http.Request) {
	ordered := shuffle(people)
	g := group(ordered, groupCount)
	t, _ := template.ParseFiles("templates/new.html")
	data := &Lunch{Group: g, Order: strings.Join(ordered, seperator)}
	t.Execute(w, data)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newLunchHanler)
	http.HandleFunc("/save", saveGroupHandler)
	log.Print("Listen at ", *addr)
	http.ListenAndServe(*addr, nil)
}
