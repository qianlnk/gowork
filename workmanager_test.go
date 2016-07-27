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
	time.Sleep(time.Second * 1)
}

func routine(name interface{}, res interface{}) {
	fmt.Printf("test routine name is %s\n", name)
	time.Sleep(time.Second * 1)
}

func timeoutDeal(song interface{}, res interface{}) {
	tmpsong := song.(string)
	tmpres := res.(*string)
	*tmpres = fmt.Sprintf("%s timeout", tmpsong)
}

func addSingRequest(gw interface{}, res interface{}) {
	var timeoutRes string
	tmpGw := gw.(*WorkManager)
	//tmpGw.SetExecptionHandler("sing", timeoutDeal, &timeoutRes)
	for i := 0; i < 1000; i++ {
		fmt.Println("sing")
		err := tmpGw.AddRequest("sing", fmt.Sprintf("song%d", i))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(timeoutRes)

	}
	tmpGw.Done("sing")
}

func TestWork(t *testing.T) {
	gw := NewWorkManager()
	res := new(MyJob)

	gw.NewGoroutine("hello", 5, hello, res)
	gw.NewGoroutine("sing", 20, sing, nil)
	gw.SetTimeout("sing", 4)

	gw.NewGoroutine("routine", 25, routine, nil)
	for i := 0; i < 2; i++ {
		gw.AddRequest("hello", strconv.Itoa(i))
	}
	gw.NewGoroutine("addSingRequest", 1, addSingRequest, nil)
	//gw.SetTimeout("addSingRequest", 50)
	gw.AddRequest("addSingRequest", gw)
	for i := 0; i < 1000; i++ {
		fmt.Println("routine")
		err := gw.AddRequest("routine", fmt.Sprintf("routine%d", i))
		if err != nil {
			fmt.Println(err)
		}
	}
	gw.Done("hello")
	gw.Done("addSingRequest")
	gw.Done("routine")

	for _, r := range res.name {
		fmt.Println(r)
	}
	time.Sleep(4 * time.Second)
}
