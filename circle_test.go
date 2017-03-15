package circle_queue

import "time"
import "testing"
import "fmt"

func TestCricle(t *testing.T) {
	work := []int{1, 5, 6, 8, 2, 4, 5, 3, 9, 0, 3, 4, 1, 0, 0, 0, 2, 3}

	c := NewCirCle(6, time.Second, func(ele []interface{}) {
		fmt.Println("time out", ele)
	})
	for _, v := range work {
		time.Sleep(time.Millisecond * 400)
		c.Put(v, true)
		fmt.Println("--->", v)
	}
	fmt.Println("===============")
	time.Sleep(time.Second * 20)

}
