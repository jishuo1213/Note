package main

import (
	"bufio"
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
	"sync"
	"time"
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

func main() {
	// findDuplicateLines()
	// animateGifs()
	// fetchUrls()
	// fetchUrlsConcurrently()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	})
	http.HandleFunc("/count", counter)
	http.HandleFunc("/favicon.ico", handlerICon)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
