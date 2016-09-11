package main

import (
	"bufio"
	// "flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	// "strings"
	// "bytes"
	"links"
	"sync"
	"time"
	// "unicode/utf8"
)

func findDuplicateLines() {
	counts := make(map[string]int)
	files := os.Args[1:]
	fmt.Printf("%xaa\n", &counts)
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

func animateGifs() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

func countLines(f *os.File, counts map[string]int) {
	fmt.Printf("%xaa\n", &counts)
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
}

func fetchUrls(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch:%v\n", err)
		os.Exit(1)
	}
	b, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch:reading %s : %v \n", url, err)
		os.Exit(1)
	}
	seconds := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", seconds, b, url)
}

func fetchUrlsConcurrently() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetchUrls(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

var mu sync.Mutex
var count int

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q \n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
	fmt.Printf("handler\n")
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.PATH = %q\n", r.URL.Path)

	lissajous(w)
}

func counter(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("counter\n")
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

func handlerICon(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

// var n = flag.Bool("n", false, "omit trailing newline")
// var sep = flag.String("s", " ", "separtor")

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func svg() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(j, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

type a struct {
	X int
}

func (pa a) ma() int {
	return pa.X
}

type b struct {
	X int
}

func (pa b) ma() int {
	return pa.X
}

type son struct {
	a
	b
}

func testWaitGroup() {
	res := make(chan int)
	var wg sync.WaitGroup
	for i := range []int{1, 2, 3, 4, 5} {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 2)
			res <- i
		}(i)
	}
	wg.Wait()
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func basename(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

func testNilSlice(s []int) []int {
	s = append(s, 1)
	return s
}

func testchan() int {
	ch := make(chan int, 3)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("first")
		ch <- 1
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("second")
		ch <- 1
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("third")
		ch <- 1
	}()

	return <-ch
}

func main() {
	// s1 := make([]int, 4)
	// s2 := []int{1, 2, 3, 4, 56, 7, 8, 9}
	// res := copy(s1, s2)
	// fmt.Println(s1)
	// fmt.Println(res)

	// var m map[string][]string
	// m := make(map[string][]string)
	// m := map[string][]string{"aaa": {"bbb"}}
	// m["aaa"] = make([]string, 5)
	// fmt.Println(cap(m["aaa"]))
	// m["aaa"] = append(m["aaa"], "aaa")
	// s := son{a{1}, b{2}}
	// fmt.Println(s.a.ma())
	// findDuplicateLines()
	// animateGifs()
	// fetchUrls()

	// fetchUrlsConcurrently()

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	lissajous(w)
	// })
	// http.HandleFunc("/count", counter)
	// http.HandleFunc("/favicon.ico", handlerICon)
	// log.Fatal(http.ListenAndServe("localhost:8000", nil))

	// flag.Parse()
	// fmt.Print(strings.Join(flag.Args(), *sep))
	// if !*n {
	// 	fmt.Println()
	// }

	// p := new(struct{})
	// q := new(struct{})

	// fmt.Println(p == q)

	// ascii := 'a'
	// unicode := '国'
	// newLine := '\n'

	// for x := 0; x < 15; x++ {
	// 	fmt.Printf("x = %d e A = %5.3f\n", x, math.Exp(float64(x)))
	// }

	// fmt.Printf("%d %[1]c %[1]q \n", ascii)
	// fmt.Printf("%d %[1]c %[1]q \n", unicode)
	// fmt.Printf("%d %[1]q \n", newLine)
	// svg()
	// s := `asdfasdfas
	// asdfasdf
	// asdfasd\\\\\asdfasd
	// asdfas`
	// fmt.Println(s)

	// s := []int{1, 2, 3, 4, 5}
	// // s := "abcdef"
	// ss := s[:]
	// fmt.Println(&s)
	// fmt.Println(&ss)
	// for i, c := range ss {
	// 	c++
	// 	// s[i]++
	// 	fmt.Println(c)
	// 	fmt.Println(ss[i])
	// }
	// fmt.Println(s)
	// fmt.Println(ss)

	// var buf bytes.Buffer
	// length, _ := buf.WriteRune('[')
	// fmt.Println(length)
	// length, _ = buf.WriteRune('樊')
	// fmt.Println(length)
	// s := buf.String()
	// fmt.Println(len(s))
	// fmt.Println(buf.String())
	// tests := []int{}
	// tests = append(tests, 1)
	//s := testNilSlice(nil)
	//fmt.Println(s)
	//fmt.Println(len(s))
	//fmt.Println(cap(s))
	// var x interface{}
	// x = true
	// switch x := x.(type) {
	// case int:
	// 	fmt.Println(x)
	// 	break
	// case bool:
	// 	if x {
	// 		fmt.Println(x)
	// 	}
	// 	break
	// }
	// strings.Contains()
	// bytes.Contains

	// fmt.Println(s)
	// fmt.Println(ss)

	worklist := make(chan []string)
	go func() { worklist <- os.Args[1:] }()

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			fmt.Println(link)
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					// worklist <- crawl(link)
					crawl(link)
				}(link)
			}
		}
	}
}
