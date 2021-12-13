package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// 因为连接池化不方便在浏览器测试channel，故用单测测试
func TestChanLimiter(t *testing.T) {
	for i := 0; i < 20; i++ {
		go func() {
			resp, _ := http.Get("http://localhost:8003/chan")
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			fmt.Println(string(body))
		}()
	}

	time.Sleep(5 * time.Second) // When your main() function ends, your program ends as well. It does not wait for other goroutines to finish.
}
