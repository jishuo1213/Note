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
	"flag"
	"net"
	"path/filepath"
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

var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
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

func walkDir(dir string, n *sync.WaitGroup, fileSize chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSize)
		} else {
			fileSize <- entry.Size()
		}
	}
}

func walkDir2(dir string, fileSize chan<- int64) {
	// defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			// n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			walkDir2(subdir, fileSize)
		} else {
			fileSize <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() {
		<-sema
	}()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1 : %v \n", err)
		return nil
	}
	return entries
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

var verbose = flag.Bool("v", false, "show verbose progerss messaged")

func du() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fmt.Println(roots)
	fileSize := make(chan int64)
	var n sync.WaitGroup
	// go func() {
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSize)
	}
	// close(fileSize)

	// }()

	go func() {
		n.Wait()
		close(fileSize)
	}()

	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	var nFiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSize:
			if !ok {
				break loop
			}
			nFiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nFiles, nbytes)
		}
	}

	printDiskUsage(nFiles, nbytes)
}

func print(pi *int) { fmt.Println(*pi) }

func testError() {
	for i := 0; i < 10; i++ {
		// defer fmt.Println(i) // OK; prints 9 ... 0
		// defer func() { fmt.Println(i) }() // WRONG; prints "10" 10 times
		// defer func(i int) { fmt.Println(i) }(i) // OK
		defer print(&i) // WRONG; prints "10" 10 times
		// go fmt.Println(i)                       // OK; prints 0 ... 9 in unpredictable order
		// go func() { fmt.Println(i) }()          // WRONG; totally unpredictable.
	}
}

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are" + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)

	for input.Scan() {
		messages <- who + ":" + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func chatServer() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

var deposits = make(chan int)
var balances = make(chan int)

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance: //这个条件表示从balance读的时候会触发
			log.Print("aaaa")
		}
	}
}

type Cake struct{ state string }

func getCake() *Cake {
	c := Cake{"aaa"}
	log.Printf("%p\n", &c)
	return &c
}

func main() {
	// cake := new(Cake)
	// testError()
	// chatServer()
	// go teller()
	// for i := 0; i < 10; i++ {
	// 	log.Print("bbb")
	// 	deposits <- 200
	// }
	// log.Print(<-balances)
	// log.Print(<-balances)
}
