package queue

import (
	"testing"
	"time"
)

func TestPush(t *testing.T) {

	type Foo struct {
		ID   int
		Name string
	}

	//id := 0
	//go func() {
	//	for {
	//		id++
	//		foo := Foo{ID: id, Name: "test1"}
	//		err := Push("test_queue", foo)
	//		if err != nil {
	//			panic(err)
	//		}
	//		//t.Logf("push succ:%+v", foo)
	//		time.Sleep(1 * time.Second)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		id++
	//		foo := Foo{ID: id, Name: "test2"}
	//		err := Push("test_queue", foo)
	//		if err != nil {
	//			panic(err)
	//		}
	//		//t.Logf("push succ:%+v", foo)
	//		time.Sleep(200 * time.Millisecond)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		id++
	//		foo := Foo{ID: id, Name: "test3"}
	//		err := Push("test_queue", foo)
	//		if err != nil {
	//			panic(err)
	//		}
	//		//t.Logf("push succ:%+v", foo)
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//}()

	go func() {
		for {
			foo := &Foo{}
			_, _, err := Pop("test_queue", foo, true)
			if err != nil {
				panic(err)
			}
			t.Logf("pop1 succ:%+v", foo)
		}
	}()

	go func() {
		for {
			foo := &Foo{}
			_, _, err := Pop("test_queue", foo, true)
			if err != nil {
				panic(err)
			}
			t.Logf("pop2 succ:%+v", foo)
		}
	}()

	time.Sleep(1 * time.Hour)
}
