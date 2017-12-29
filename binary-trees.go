package main

import (
	"fmt"
	"golang.org/x/tour/tree"
	"reflect"
	"sort"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int, quit chan int) {
	if t.Value != 0 {
		ch <- t.Value
	}

	if (t.Left == nil && t.Right == nil) {
		// quit <- 1
	}

	if t.Left != nil {
		go Walk(t.Left, ch, quit)
	}

	if t.Right != nil {
		go Walk(t.Right, ch, quit)
	}

	// fmt.Println(reflect.TypeOf(t.Left))
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	quit := make(chan int, 1)

	go Walk(t1, ch1, quit)
	go Walk(t2, ch2, quit)

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
        case msg3 := <-quit:
            fmt.Println("quit received", msg3)
			      // return false;
        }
    }

	sort.Ints(s1)
	sort.Ints(s2)

	return reflect.DeepEqual(s1, s2);
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
