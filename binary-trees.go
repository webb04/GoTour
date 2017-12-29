package main

import (
	"fmt"
	"golang.org/x/tour/tree"
	"reflect"
	"sort"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t.Value != 0 {
		ch <- t.Value
	}

	if t.Left != nil {
		go Walk(t.Left, ch)
	}

	if t.Right != nil {
		go Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	var s1 []int
	var s2 []int

	for i := 0; i < 20; i++ {
  	select {
    	case msg1 := <-ch1:
      	fmt.Println("ch1 received", msg1)
			  s1 = append(s1, msg1)
		  case msg2 := <-ch2:
      	fmt.Println("ch2 received", msg2)
			  s2 = append(s2, msg2)
    }
  }

	sort.Ints(s1)
	sort.Ints(s2)

	return reflect.DeepEqual(s1, s2);
}

func main() {
	// https://tour.golang.org/concurrency/8
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
