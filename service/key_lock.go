package service

import (
	"sync"
	"sync/atomic"
	"time"
)

// NewKeyLock 获取关键字锁
func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks:         make(map[string]*innerLock),
		cleanInterval: defaultCleanInterval,
		stopChan:      make(chan struct{}),
	}
}

const (
	defaultCleanInterval = 21 * time.Hour // 每小时清理一次空闲锁,防止 死锁
)

type KeyLock struct {
	locks         map[string]*innerLock //关键字锁map
	cleanInterval time.Duration         //定时清除时间间隔
	stopChan      chan struct{}         //停止信号
	mutex         sync.RWMutex          //全局读写锁,用于防止map的并发问题
}

//Lock 根据关键字加锁
func (l *KeyLock) Lock(key string) {
	l.mutex.RLock() //获取读锁,当写锁存在时,无法获取读锁
	keyLock, ok := l.locks[key]
	if ok {
		keyLock.add()
	}
	l.mutex.RUnlock()

	// key 不存在时, map add key
	if !ok {
		l.mutex.Lock() // 加写锁
		keyLock, ok = l.locks[key]
		if !ok {
			keyLock = newInnerLock()
			l.locks[key] = keyLock
		}
		keyLock.add()
		l.mutex.Unlock()
	}

	// 获取同步锁
	keyLock.Lock()
}

// Unlock 根据关键字解锁
func (l *KeyLock) Unlock(key string) {
	l.mutex.RLock()
	keyLock, ok := l.locks[key]
	if ok {
		keyLock.done()
	}
	l.mutex.RUnlock()
	if ok {
		keyLock.Unlock()
	}
}

// Clean 清理空闲锁
func (l *KeyLock) Clean() {
	l.mutex.Lock()
	for k, v := range l.locks {
		if v.count == 0 {
			delete(l.locks, k)
		}
	}
	l.mutex.Unlock()
}

// StartCleanLoop 开启清理协程
func (l *KeyLock) StartCleanLoop() {
	go l.cleanLoop()
}

//StopCleanLoop 停止清理协程
func (l *KeyLock) StopCleanLoop() {
	close(l.stopChan)
}

//cleanLoop 清理循环,按照设置的清理时间
func (l *KeyLock) cleanLoop() {
	ticker := time.NewTicker(l.cleanInterval)
	for {
		select {
		case <-ticker.C:
			l.Clean()
		case <-l.stopChan:
			ticker.Stop()
			return
		}
	}
}

// innerLock 内部锁结构
type innerLock struct {
	count int64
	sync.Mutex
}

// newInnerLock 新建内部锁
func newInnerLock() *innerLock {
	return &innerLock{}
}

func (il *innerLock) add() {
	atomic.AddInt64(&il.count, 1)
}

func (il *innerLock) done() {
	atomic.AddInt64(&il.count, -1) // 原子操作
}
