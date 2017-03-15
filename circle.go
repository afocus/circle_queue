package circle_queue

import "time"

type CirCle struct {
	// 这个环由多少个块构成
	part int
	// 拼成一个环的块集合
	list []*Set
	// 当前移动到的块索引
	index int
	// 值对应的块索引关系
	bindmap map[interface{}]int

	taskIn  chan indata
	taskOut chan interface{}
	closed  chan struct{}

	// 当移动到位置后且有数据会执行此方法，
	call callbackFunc
}

type indata struct {
	value interface{}
	must  bool
}

type callbackFunc func([]interface{})

// NewCirCle 创建
func NewCirCle(part int, intevel time.Duration, callback callbackFunc) *CirCle {
	part++
	setlist := make([]*Set, part)
	for i := 0; i < part; i++ {
		setlist[i] = NewSet()
	}
	c := &CirCle{
		list:    setlist,
		part:    part,
		bindmap: make(map[interface{}]int),
		taskOut: make(chan interface{}, 1024),
		taskIn:  make(chan indata, 1024),
		closed:  make(chan struct{}),
		call:    callback,
	}
	go c.task(intevel)
	return c
}

func (c *CirCle) task(intevel time.Duration) {
	//
	timer := time.NewTimer(intevel)

	for {
		select {
		case in := <-c.taskIn:
			if index, has := c.bindmap[in.value]; has {
				if !in.must {
					return
				}
				// 已存在同样的则删除以前的
				// 主要用阶梯超时等场景
				delete(c.bindmap, in.value)
				c.list[index].Del(in.value)
			}
			// 放入前一格
			index := (c.part + c.index - 2) % c.part
			c.list[index].Add(in.value)
			c.bindmap[in.value] = index

		case v := <-c.taskOut:
			// 如果任务完成则从里面剔除
			if index, has := c.bindmap[v]; has {
				delete(c.bindmap, v)
				c.list[index].Del(v)
			}

		case <-timer.C:
			// 检查当前格内是否有超时对象
			item := c.list[c.index]
			if item.Len() != 0 {
				ele := item.Elements()
				for k := range ele {
					delete(c.bindmap, ele[k])
				}
				item.Clear()
				// 异步执行防止堵塞
				go c.call(ele)
			}
			// 向后移动一格
			c.index = (c.index + 1) % c.part
			// 重置定时器
			timer.Reset(intevel)

		case <-c.closed:
			return
		}
	}
}

// Put 把v放入环中，flag==true 不管之前是否存在依然放入/替换 false不在环中才放入
func (c *CirCle) Put(v interface{}, flag bool) {
	c.taskIn <- indata{
		value: v,
		must:  flag,
	}
}

// Pop 从环中移除
func (c *CirCle) Pop(v interface{}) {
	c.taskOut <- v
}

// Close 不再接收put,pop 返回当前环中所有数据
func (c *CirCle) Close() (ret [][]interface{}) {
	close(c.closed)
	for _, v := range c.list {
		ret = append(ret, v.Elements())
	}
	return
}
