package gls

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// https://chai2010.cn/advanced-go-programming-book/ch3-asm/ch3-08-goroutine-id.html
// goid的应用: 局部存储
// 有了goid之后，构造Goroutine局部存储就非常容易了。我们可以定义一个gls包提供goid的特性：

// gls包变量简单包装了map，同时通过sync.Mutex互斥量支持并发访问。
var gls struct {
	m map[int32]map[interface{}]interface{}
	sync.Mutex
}

func init() {
	gls.m = make(map[int32]map[interface{}]interface{})
}

// 定义一个getMap内部函数，用于获取每个Goroutine字节的map：
func getMap() map[interface{}]interface{} {
	gls.Lock()
	defer gls.Unlock()

	goid := GetGoid()
	if m, _ := gls.m[goid]; m != nil {
		return m
	}

	m := make(map[interface{}]interface{})
	gls.m[goid] = m
	return m
}

// 获取到Goroutine私有的map之后，就是正常的增、删、改操作接口了：
func Get(key interface{}) interface{} {
	return getMap()[key]
}
func Put(key interface{}, v interface{}) {
	getMap()[key] = v
}
func Delete(key interface{}) {
	delete(getMap(), key)
}

// 再提供一个Clean函数，用于释放Goroutine对应的map资
func Clean() {
	gls.Lock()
	defer gls.Unlock()

	delete(gls.m, GetGoid())
}

func GetGoid() int32 {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return int32(id)
}
