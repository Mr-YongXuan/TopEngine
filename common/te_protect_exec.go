package common

import "sync"

var lock sync.Mutex

func Protect(function func()) {
	lock.Lock()
	function()
	lock.Unlock()
}