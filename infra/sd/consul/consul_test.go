package consul

import (
	"context"
	"testing"
	"time"
)

func TestNewConsulSD(t *testing.T) {
	// 前提：在本机启动consul进程
	sd, err := New()
	if err != nil {
		t.Fatalf("New %v", err)
	}
	svc := "user"
	addr := "127.0.0.1"
	port := 8500 // 使用consul的端口
	err = sd.Register(svc, addr, port, nil)
	if err != nil {
		t.Fatalf("Register %v", err)
	}
	//return
	// 首次查询 不阻塞 所以不会超时
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
	list, err := sd.Discover(ctx, svc, false)
	if err != nil {
		t.Fatalf("discover %v", err)
	}
	//time.Sleep(time.Second) // consul delay
	if len(list) == 0 {
		t.Fatalf("len(list) == 0")
	}
	if list[0].Host != addr && list[0].Port != port {
		t.Fatalf("unexpected ret:%+v", list)
	}
	//return
	err = sd.Deregister(svc)
	if err != nil {
		t.Fatalf("Deregister %v", err)
	}
	list, err = sd.Discover(ctx, svc, false)
	if err != nil {
		t.Fatalf("discover %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("len(list) != 0")
	}
	//now := time.Now()
	// 观察多次查询的阻塞时间是否符合预期
	// 第一次3s但返回超时err，第二次和第三次 1min (Discover方法内部设置的)

	//firstDur := 3
	//subsequentDur := 60
	//for i := 0; i < 3; i++ {
	//	if i == 0 {
	//		ctx, _ = context.WithTimeout(context.TODO(), time.Second*3)
	//	} else {
	//		ctx = context.TODO()
	//	}
	//	list, err = abstract.Discover(ctx, svc)
	//	dur := int(time.Now().Sub(now).Seconds())
	//	if i == 0 {
	//		if err != context.DeadlineExceeded || dur != firstDur || len(list) != 0 {
	//			t.Errorf("no.%v Discover in for loop， unexpected result, err:%v dur:%ds", i, err, dur)
	//			now = time.Now()
	//			continue
	//		}
	//	} else {
	//		if err != nil || dur != subsequentDur || len(list) != 1 {
	//			t.Errorf("no.%v Discover in for loop， unexpected result, err:%v dur:%ds, list:%+v", i, err, dur, list)
	//			now = time.Now()
	//			continue
	//		}
	//	}
	//	now = time.Now()
	//}

}
