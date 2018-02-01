package numbers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Nums struct
type Nums struct {
	Numbers []int `json:"numbers"`
}

var (
	errMissingParam = errors.New("missing parameter 'u' on request")
)

// OrderHandler return numbers in ascending order on http.Response
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()
	urls, ok := r.URL.Query()["u"]
	if !ok {
		http.Error(w, errMissingParam.Error(), http.StatusBadRequest)
		return
	}
	intChan := make(chan []int)
	go exec(ctx, urls, intChan)
	w.Header().Set("Content-Type", "application/json")
	nums := make([]int, 0)
	select {
	case <-ctx.Done():
	case nums = <-intChan:
	}
	if err := json.NewEncoder(w).Encode(Nums{nums}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// parseURL fetch data
func parseURL(ctx context.Context, wg *sync.WaitGroup, url string, intChan chan []int) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusOK {
		ints := Nums{}
		err = json.Unmarshal(body, &ints)
		if err != nil {
			return
		}
		intChan <- ints.Numbers
	}
}

// exec sort
func exec(ctx context.Context, urls []string, intChan chan []int) {
	internalChan := make(chan []int)
	nums := make([]int, 0)
	wg := &sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go parseURL(ctx, wg, url, internalChan)
	}
	go func(ctx context.Context) {
		for resp := range internalChan {
			nums = append(nums, resp...)
		}
		sort.Ints(nums)
		nums = removeDuplicates(nums)
		intChan <- nums
	}(ctx)
	wg.Wait()
	close(internalChan)
}

// removeDuplicates entries
func removeDuplicates(elements []int) []int {
	encountered := make(map[int]bool)
	result := make([]int, 0)
	for _, e := range elements {
		if !encountered[e] {
			encountered[e] = true
			result = append(result, e)
		}
	}
	return result
}
