package gowork

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type MyJob struct {
	name []string
	mu   sync.Mutex
}

func hello(name interface{}, res interface{}) {
	tmp := res.(*MyJob)
	tmp.mu.Lock()
	defer tmp.mu.Unlock()
	tmp.name = append(tmp.name, fmt.Sprintf("hello, %s.", name.(string)))
}

func sing(song interface{}, res interface{}) {
	fmt.Printf("song name is %s\n", song)
	time.Sleep(time.Second * 5)
}

func routine(name interface{}, res interface{}) {
	fmt.Printf("test routine name is %s\n", name)
	time.Sleep(time.Nanosecond * 100)
}

func TestWork(t *testing.T) {
	gw := NewWorkManager()
	res := new(MyJob)

	gw.NewGoroutine("hello", 5, hello, res)
	gw.NewGoroutine("sing", 8, sing, nil)
	gw.SetTimeout("sing", 4)
	gw.NewGoroutine("routine", 7, routine, nil)
	for i := 0; i < 2; i++ {
		gw.AddRequest("hello", strconv.Itoa(i))
	}
	go func() {
		for i := 0; i < 4000; i++ {
			fmt.Println("sing")
			gw.AddRequest("sing", fmt.Sprintf("song%d", i))
		}
		gw.Done("sing")
	}()
	for i := 0; i < 4000; i++ {
		fmt.Println("routine")
		err := gw.AddRequest("routine", fmt.Sprintf("routine%d", i))
		if err != nil {
			fmt.Println(err)
		}
	}
	gw.Done("hello")
	gw.Done("routine")

	for _, r := range res.name {
		fmt.Println(r)
	}
	time.Sleep(40 * time.Second)
}
