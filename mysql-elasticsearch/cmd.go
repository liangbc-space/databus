package mysql_elasticsearch

import (
	"fmt"
)

func Run(){
	done := make(chan bool)

	go func() {
		//	获取消息
		poolMessages()
		done <- true
	}()

	select {
	case <-done:
		fmt.Printf("退出")
	}
}

