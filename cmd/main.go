package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	c "github.com/strjkc/cacheproxy/internal/cache"
	s "github.com/strjkc/cacheproxy/internal/state"
)

var validParams []string = []string{"--port", "--origin"}
var state *s.State
var cache *c.Cache

func main() {
	args := os.Args
	if len(args) != 5 {
		fmt.Println("Invalid use")
		return
	}
	if args[1] != validParams[0] || args[3] != validParams[1] {
		fmt.Println("Invalid param")
		return
	}
	port := args[2]
	origin := args[4]
	fmt.Printf("port: %v\norigin:%v", port, origin)
	state = s.NewState(origin, 30)
	cache = c.NewCache()
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	go cache.Reap()
	for {
		err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	cache.ReportStats(2)
	path := r.URL.Path
	newUrl := state.BaseURL + path
	item, ok := cache.GetItem(newUrl)
	if ok {
		fmt.Println("Item found in cache!")
		expired := cache.IsItemExpired(newUrl)
		if !expired {
			fmt.Println("Item not expired! Returning it")
			w.Header().Set("X-Cache", "HIT")
			w.Write(item.Reply)
			fmt.Println("")
			return
		}
		fmt.Println("Item expired! Fetching it")
		resp, err := getAndCache(newUrl)
		if err != nil {
			w.Write([]byte("error"))
		}
		w.Header().Set("X-Cache", "MISS")
		io.Copy(w, resp.Body)
	} else {
		fmt.Println("Item not found in cache! Fetching it")
		resp, err := getAndCache(newUrl)
		if err != nil {
			w.Write([]byte("error"))
		}
		w.Header().Set("X-Cache", "MISS")
		io.Copy(w, resp.Body)
	}
	fmt.Println("")
}

func getAndCache(newUrl string) (*http.Response, error) {
	resp, err := http.Get(newUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return &http.Response{}, err
	}
	item := c.CacheEntry{Reply: b, Fetched: time.Now()}
	//TODO: maybe have one function that either caches new or updates
	cache.CacheItem(item, newUrl)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}

/*
 * u sate imam i cache expiery u sekundama?
 * imamo reap background funkciju koja treba da cisti sve unose za koje je proslo x sekundi od timestampea
 * ova funkcija mora da ima neki lock nad mapom, tako da ako ona reapuje neki record zahtev koji dolazi u isto vreme mora da saceka.
 * ovo onda treba nekako poboljsati, da nije lock na celoj mapi nego na recordu, jer onda zahtev ceka da se ocisti cela mapa mozda.
 * shardovanje radimo
 */
