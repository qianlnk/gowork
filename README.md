# gowork
	gowork is a package for control goroutine number.You may compile func worker by your self, then func type is:</br>
```golang
type WorkFunction func(request interface{}, response interface{})
```
	you can transfer params by `request`, of course,if you have more params you may packaging them as a struct.
# use
if the worker has no param and result.
```golang
package main

import(
	"fmt"

	"github.com/qianlnk/gowork"
)
func hello(req interface{}, res interface{}) {
	fmt.Println("hello, I'm qianlnk")
}

func main() {
	wm := gowork.NewWorkManager()
	wm.NewGoroutine("hello", 5, hello, nil)

	for i := 0; i < 10; i++ {
		wm.AddRequest("hello", nil)
	}

	wm.Done("hello")
}
```
the worker with params and result.
```golang
package main

import (
	"fmt"
	"sync"

	"github.com/qianlnk/gowork"
)

type People struct {
	Id   int
	Name string
}

type Hello struct {
	Hellos []string
	Mutex  sync.Mutex
}

func SayHello(req, res interface{}) {
	tmpreq := req.(People)
	tmpres := res.(*Hello)
	tmpres.Mutex.Lock()
	defer tmpres.Mutex.Unlock()
	tmpres.Hellos = append(tmpres.Hellos, fmt.Sprintf("Hello, I'm %s, My ID is %d.", tmpreq.Name, tmpreq.Id))
}

func main() {
	wm := gowork.NewWorkManager()
	result := new(Hello)
	wm.NewGoroutine("sayhello", 5, SayHello, result)
	for i := 1; i < 50; i++ {
		pl := People{
			Id:   i,
			Name: fmt.Sprintf("test%d", i),
		}
		wm.AddRequest("sayhello", pl)
	}
	wm.Done("sayhello")

	for _, res := range result.Hellos {
		fmt.Println(res)
	}
}
```
and, you can create more worker also.
```golang
package main

import (
	"fmt"

	"github.com/qianlnk/gowork"
)

type People struct {
	Id   int
	Name string
}

func SayHello(request, response interface{}) {
	req := request.(People)
	fmt.Printf("Hello, I'm %s, My ID is %d.\n", req.Name, req.Id)
}

func SingSong(request, response interface{}) {
	req := request.(People)
	fmt.Printf("My ID is %d, what I will sing is test_%s\n", req.Id, req.Name)
}

func main() {
	wm := gowork.NewWorkManager()
	wm.NewGoroutine("sayhello", 5, SayHello, nil)
	wm.NewGoroutine("singsong", 5, SingSong, nil)
	go func() {
		for i := 1; i < 20; i++ {
			pl := People{
				Id:   i,
				Name: fmt.Sprintf("hello%d", i),
			}
			wm.AddRequest("sayhello", pl)
		}
	}()
	for i := 0; i < 20; i++ {
		pl := People{
			Id:   i,
			Name: fmt.Sprintf("song%d", i),
		}
		wm.AddRequest("singsong", pl)
	}
	wm.Done("sayhello")
	wm.Done("singsong")
}

```
result is 
```golang
My ID is 0, what I will sing is test_song0
My ID is 1, what I will sing is test_song1
My ID is 2, what I will sing is test_song2
My ID is 3, what I will sing is test_song3
Hello, I'm hello1, My ID is 1.
Hello, I'm hello6, My ID is 6.
Hello, I'm hello7, My ID is 7.
Hello, I'm hello8, My ID is 8.
Hello, I'm hello9, My ID is 9.
Hello, I'm hello10, My ID is 10.
Hello, I'm hello11, My ID is 11.
Hello, I'm hello12, My ID is 12.
Hello, I'm hello13, My ID is 13.
Hello, I'm hello14, My ID is 14.
Hello, I'm hello15, My ID is 15.
My ID is 4, what I will sing is test_song4
Hello, I'm hello16, My ID is 16.
Hello, I'm hello17, My ID is 17.
Hello, I'm hello18, My ID is 18.
Hello, I'm hello19, My ID is 19.
My ID is 9, what I will sing is test_song9
My ID is 10, what I will sing is test_song10
Hello, I'm hello2, My ID is 2.
My ID is 11, what I will sing is test_song11
My ID is 12, what I will sing is test_song12
My ID is 13, what I will sing is test_song13
Hello, I'm hello3, My ID is 3.
My ID is 14, what I will sing is test_song14
My ID is 15, what I will sing is test_song15
My ID is 16, what I will sing is test_song16
Hello, I'm hello4, My ID is 4.
My ID is 17, what I will sing is test_song17
My ID is 18, what I will sing is test_song18
Hello, I'm hello5, My ID is 5.
My ID is 5, what I will sing is test_song5
My ID is 6, what I will sing is test_song6
My ID is 19, what I will sing is test_song19
My ID is 7, what I will sing is test_song7
My ID is 8, what I will sing is test_song8
```