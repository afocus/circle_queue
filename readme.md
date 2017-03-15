# circle_queue

一个环形任务队列的简单实现

## 使用场景

### 超时任务
```go
// 创建一个由10个格子组成的环，每2秒移动一格， 20秒算超时
c := NewCircle(10, time.Second*2)
// 当移动到某个格子且有数据时，执行此方法
c.OnPart = func(values []interface{}){
    for _, v := range values{
        log.Println("time out", v)
    }
}
// 模拟一个任务
task := newTask(1)
// 把此任务的序号放入环中
c.Put(1)
// 模拟超时 任务执行30秒
time.Sleep(time.Secound*30)
c.Pop(1)
```



### 延时任务
```go
// 创建一个由10个格子组成的环，每2秒移动一格
c := NewCircle(10, time.Second*2)
// 当移动到某个格子且有数据时，执行此方法
c.OnPart = func(values []interface{}){
    for _,v:=range values{
        log.Println("do", v)
    }
}
// 模拟一个任务
task := newTask()
// 把此任务的序号放入环中 20秒后执行
c.Put(task)
```