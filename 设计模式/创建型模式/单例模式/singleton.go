package singleton

import "sync"

// 私有内部类 结构体
type singleton struct {
	count int
}

var instance *singleton

// 懒汉模式-非并发
func GetInstance() *singleton {
	if instance == nil {
		//instance = new(singleton)
		instance = &singleton{}
	}

	return instance
}

func (s *singleton) AddOne() int {
	s.count++
	return s.count
}

func (s *singleton) GetCount() int {
	return s.count
}

/**************************** 知识扩展分解线 ********************************/
// 懒汉模式-并发-耗资源版
var mutx sync.Mutex

func GetInstance2() *singleton {
	mutx.Lock()
	defer mutx.Unlock()
	if instance == nil {
		instance = &singleton{}
	}

	return instance
}

// 懒汉模式-并发-优化资源版（双检查）
func GetInstancConcurrent() *singleton {
	if instance == nil {
		// 这里可能多个线程处于此处
		mutx.Lock()
		if instance == nil {
			instance = &singleton{} 
		}
		mutx.Unlock()
	}
	return instance // 当一线程初始化后返回并赋值instance，则 instance 为 非nil
}

var once sync.Once

// go 特有包单例模式高并发终极版推荐
// 支持高并发单例模式（优雅）
func GetInstanceOnce() *singleton {
	once.Do(func() {
		if instance == nil {
			instance = new(singleton)
		}
	})
	return instance
}
