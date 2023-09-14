package simple_sd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHTTPServer_handlePing(t *testing.T) {
	type XX struct {
		name    string
		req     interface{}
		wantRes *HttpRes
	}
	tests := []XX{
		{
			name:    "T-invalid ping body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(PingReq)).Error())),
		},
		{
			name:    "T-ping body:false",
			req:     PingReq{Ping: false},
			wantRes: newRes(nil, 200, nil),
		},
		{
			name:    "T-ping body:true",
			req:     PingReq{Ping: true},
			wantRes: newRes(PingRspBody{Pong: true}, 200, nil),
		},
	}

	for _, tt := range tests {
		Init() // 每次重新初始化

		req, err := http.NewRequest("POST", "/ping", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handlePing)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got：%s, not want:%+v", b, tt.wantRes)
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func TestHTTPServer_handleRegister(t *testing.T) {
	type XX struct {
		name    string
		req     interface{}
		wantRes *HttpRes
	}
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(RegisterReq)).Error())),
		},
		{
			name: "T-空service",
			req: ServiceInstance{
				Name: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, ServiceInstance{}.Check().Error())),
		},
		{
			name: "T-空id",
			req: ServiceInstance{
				Name: "go-user",
				Id:   "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, ServiceInstance{Name: "go-user", Id: ""}.Check().Error())),
		},
		{
			name: "T-空Host",
			req: ServiceInstance{
				Id:   "any id",
				Name: "go-user",
				Host: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, ServiceInstance{Name: "go-user", Id: "any id", Host: ""}.Check().Error())),
		},
		{
			name: "T-空Port",
			req: ServiceInstance{
				Id:   "any id",
				Name: "go-user",
				Host: "localhost",
				Port: 0,
			},
			wantRes: newRes(nil, 400,
				errors.Wrap(ErrParams, ServiceInstance{Name: "go-user", Id: "any id", Host: "localhost", Port: 0}.Check().Error())),
		},
		{
			name: "T-有效实例",
			req: ServiceInstance{
				Id:   "any id",
				Name: "go-user",
				Host: "localhost",
				Port: 8080,
			},
			wantRes: newRes([]*ServiceInstance{
				{
					Id:   "any id",
					Name: "go-user",
					Host: "localhost",
					Port: 8080,
				},
			}, 200, nil),
		},
	}

	for _, tt := range tests {
		Init() // 每次重新初始化

		req, err := http.NewRequest("POST", "/service/register", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handleRegister)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got：%s, not want:%+v", b, tt.wantRes)
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func checkItem(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("---NotEqual a:%+v b:%+v", a, b)
	}
}

func TestHTTPServer_handleDeregister(t *testing.T) {
	type XX struct {
		name     string
		req      interface{}
		regFirst bool
		wantRes  *HttpRes
	}
	instId := "a_instance_id"
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(DeregisterReq)).Error())),
		},
		{
			name: "T-空service",
			req: DeregisterReq{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "provide a valid service, id")),
		},
		{
			name: "T-空id",
			req: DeregisterReq{
				Service: "go-user",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "provide a valid service, id")),
		},
		{
			name: "T-注销未注册的实例",
			req: DeregisterReq{
				Service: "go-user",
				Id:      "any id",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrInstanceNotRegistered,
				fmt.Sprintf("service: %s, id: %s", "go-user", "any id"))),
		},
		{
			name: "T-正常注销实例",
			req: DeregisterReq{
				Service: "go-user",
				Id:      instId,
			},
			regFirst: true,
			wantRes:  newRes(nil, 200, nil),
		},
	}

	for _, tt := range tests {
		Init() // 每次重新初始化

		if tt.regFirst {
			inst := tt.req.(DeregisterReq)
			err := Sd.Register(ServiceInstance{
				Id:       instId,
				Name:     inst.Service,
				IsUDP:    false,
				Host:     inst.Id,
				Metadata: nil,
			})
			if err != nil {
				t.Fatalf(tt.name+"--- Register failed %v", err)
			}
		}
		req, err := http.NewRequest("POST", "/service/deregister", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handleDeregister)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got:\n%s\n--- not want:\n%s", b, ToJson(tt.wantRes))
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func TestHTTPServer_handleDiscovery(t *testing.T) {
	type XX struct {
		name              string
		req               interface{}
		regAfter          time.Duration
		regInstances      []ServiceInstance
		regInstancesAfter []ServiceInstance
		shouldCostMs      int64
		wantRes           *HttpRes
	}

	__inst1 := []ServiceInstance{{
		Id:   "inst1",
		Name: "go-user",
		Host: "localhost",
		Port: 8080,
	}}
	__inst2 := []ServiceInstance{{
		Id:   "inst2",
		Name: "go-user",
		Host: "localhost",
		Port: 8081,
	}}

	__oneInstHash := CalInstanceHash(__inst1)
	__twoInstHash := CalInstanceHash(append(__inst1, __inst2...))
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(DiscoveryReq)).Error())),
		},
		{
			name: "T-空service",
			req: &DiscoveryReq{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "params without service")),
		},
		{
			name: "T-设置Service、无Hash，有1个实例（应立即返回）",
			req: &DiscoveryReq{
				Service: "go-user",
			},
			regInstances: __inst1,
			wantRes: newRes(DiscoveryRspBody{
				Instances: __inst1,
				Hash:      __oneInstHash,
			}, 200, nil),
		},
		{
			name: "T-设置Service、无Hash、WaitMaxMs=10, 无实例（应立即返回）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10, // 即使设置了也不会等待，因为无Hash
			},
			wantRes: newRes(DiscoveryRspBody{Instances: __inst1[:0], Hash: EmptyInstanceHash}, 200, nil),
		},
		{
			name: "T-设置Service、无Hash、WaitMaxMs=10, 有实例（应立即返回）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10, // 即使设置了也不会等待，因为无Hash
			},
			regInstances: __inst1,
			wantRes:      newRes(DiscoveryRspBody{Instances: __inst1, Hash: __oneInstHash}, 200, nil),
		},
		{
			name: "T-设置Service、无实例对应Hash、WaitMaxMs=10, 无实例（应等待10ms再返回空实例以及相同hash）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10,
				LastHash:  EmptyInstanceHash,
			},
			shouldCostMs: time.Millisecond.Milliseconds() * 10,
			wantRes:      newRes(DiscoveryRspBody{Instances: __inst1[:0], Hash: EmptyInstanceHash}, 200, nil),
		},
		{
			name: "T-设置Service、无实例对应Hash、WaitMaxMs=20, 10ms后注册1个实例（应等待10ms再返回1个实例）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 20, // 稍微大于 regAfter
				LastHash:  EmptyInstanceHash,
			},
			regAfter:          time.Millisecond * 10,
			regInstancesAfter: __inst1,
			shouldCostMs:      time.Millisecond.Milliseconds() * 10,
			wantRes:           newRes(DiscoveryRspBody{Instances: __inst1, Hash: __oneInstHash}, 200, nil),
		},
		{
			name: "T-设置Service、单实例的Hash、WaitMaxMs=10, 注册单实例（应等待10ms再返回1个相同实例）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10,
				LastHash:  __oneInstHash,
			},
			shouldCostMs: time.Millisecond.Milliseconds() * 10,
			regInstances: __inst1,
			wantRes: newRes(DiscoveryRspBody{
				Instances: __inst1,
				Hash:      __oneInstHash,
			}, 200, nil),
		},
		{
			name: "T-设置Service、单实例的Hash、WaitMaxMs=1010, 注册单实例，1000ms后再注册另一实例（应等待1000ms再返回2个实例）",
			req: &DiscoveryReq{
				Service:   "go-user",
				WaitMaxMs: 1050, // 需要稍微大于 regAfter
				LastHash:  __oneInstHash,
			},
			shouldCostMs:      time.Millisecond.Milliseconds() * 1000,
			regInstances:      __inst1,
			regAfter:          time.Millisecond * 1000,
			regInstancesAfter: __inst2,
			wantRes: newRes(DiscoveryRspBody{
				Instances: append(__inst1, __inst2...),
				Hash:      __twoInstHash,
			}, 200, nil),
		},
	}

	for _, tt := range tests {
		Init() // 每次重新初始化
		t.Logf("run test: %v", tt.name)
		go func(instances []ServiceInstance) {
			time.Sleep(tt.regAfter)
			for _, inst := range instances {
				err := Sd.Register(inst)
				if err != nil {
					t.Fatalf(tt.name+"--- Register failed %v", err)
				}
			}
		}(tt.regInstancesAfter)

		for _, inst := range tt.regInstances {
			err := Sd.Register(inst)
			if err != nil {
				t.Fatalf(tt.name+"--- Register failed %v", err)
			}
		}
		req, err := http.NewRequest("POST", "/service/discovery", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		st := time.Now()
		handler := http.HandlerFunc(handleDiscovery)
		handler.ServeHTTP(rr, req)
		cost := time.Since(st).Milliseconds()

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		//body, ok := tt.req.(*DiscoveryReq)
		//if ok && body.WaitMaxMs > 0 && !isMillsTimeCostEqual(cost, tt.regAfter.Milliseconds()) {
		//	t.Fatalf(tt.name+"--handler unexpected waitMs: got %v, want %v", cost, tt.regAfter.Milliseconds())
		//}
		if tt.shouldCostMs < cost-20 { // 减去一个误差时间（不知道哪冒出来的）
			t.Fatalf(tt.name+"--handler unexpected cost-ms: got %v, want %v", cost, tt.shouldCostMs)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got:\n%s\n--- not want:\n%s", b, ToJson(tt.wantRes))
		}

		for _, inst := range append(__inst1, __inst2...) {
			_ = Sd.Deregister(inst.Name, inst.Id)
		}

		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func doRegister(insts []ServiceInstance) []byte {
	var res []byte
	for _, inst := range insts {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleRegister)
		req, _ := http.NewRequest("POST", "/service/register", bytes.NewReader(ToJson(RegisterReq{inst})))
		handler.ServeHTTP(rr, req)
		res = rr.Body.Bytes()
	}
	return res
}

func doDiscovery(service string) []byte {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDiscovery)
	req, _ := http.NewRequest("POST", "/service/discovery", bytes.NewReader(ToJson(&DiscoveryReq{
		Service: service,
	})))
	handler.ServeHTTP(rr, req)
	return rr.Body.Bytes()
}

func mustHaveInstance(t *testing.T, insts []ServiceInstance, hash string) {
	var wantRes = &DiscoveryRspBody{
		Instances: insts,
		Hash:      hash,
	}
	resBytes := doDiscovery("go-user")

	if !bytes.Equal(resBytes, ToJson(newRes(wantRes, 200, nil))) {
		t.Fatalf("Discover result failed, got res:%s", resBytes)
	}
}

func mustNoInstance(t *testing.T) {
	var wantRes = &DiscoveryRspBody{}
	resBytes := doDiscovery("go-user")

	if !bytes.Equal(resBytes, ToJson(newRes(wantRes, 200, nil))) {
		t.Fatalf("Discover result failed, got res:%s", resBytes)
	}
}

func TestHTTPServer_handleDiscoveryWithHealthCheckFail(t *testing.T) {
	Init() // 每次重新初始化

	__inst1 := []ServiceInstance{{
		Name: "go-user",
		Host: "localhost",
		Port: 8080,
	}}
	//__inst2 := []*core.ServiceInstance{{
	//	Name: "go-user",
	//	Host:    "localhost",
	//	Port:    8081,
	//}}

	__oneInstHash := CalInstanceHash(__inst1)
	//__twoInstHash := core.CalInstanceHash(append(__inst1, __inst2...))

	resBytes := doRegister(__inst1)
	if !bytes.Equal(resBytes, ToJson(newRes(__inst1, 200, nil))) {
		t.Fatalf("doRegister failed: %s", resBytes)
	}

	// 在sd进行健康检测期间，实例应该一直存在
	serverAliveSec := HealthCheckInterval.Seconds() * HealthCheckMaxFails
	for i := 0; i < int(serverAliveSec); i++ {
		println("--- No.1 sleep serverAliveSec", i+1)
		mustHaveInstance(t, __inst1, __oneInstHash)
		time.Sleep(time.Second)
	}

	println("No.1 sleep end...")
	time.Sleep(time.Second) // 给一个合理的间隙时间
	mustNoInstance(t)
}

func TestHTTPServer_handleDiscoveryWithHealthCheckPass(t *testing.T) {
	Init() // 每次重新初始化

	__inst1 := []ServiceInstance{{
		Name: "go-user",
		Host: "localhost",
		Port: 8080,
	}}
	//__inst2 := []*core.ServiceInstance{{
	//	Name: "go-user",
	//	Host:    "localhost",
	//	Port:    8081,
	//}}

	__oneInstHash := CalInstanceHash(__inst1)
	//__twoInstHash := core.CalInstanceHash(append(__inst1, __inst2...))

	//
	resBytes := doRegister(__inst1)
	// -- 返回实例是按注册时间排序的，所以这里顺序固定
	if !bytes.Equal(resBytes, ToJson(newRes(__inst1, 200, nil))) {
		t.Fatalf("doRegister failed: %s", resBytes)
	}

	// -- 同时启动实例 以通过健康检测
	__server := http.Server{Addr: __inst1[0].Addr()}
	serverAliveSec := HealthCheckInterval.Seconds() * (HealthCheckMaxFails + 1)
	go func() {
		err := __server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
		println("No.2 server closed")
	}()

	time.AfterFunc(time.Duration(serverAliveSec)*time.Second, func() {
		__server.Close()
	})

	// 以更长的时间来持续检测实例是否注册
	for i := 0; i < int(serverAliveSec); i++ {
		println("--- No.2 sleep serverAliveSec", i+1)
		mustHaveInstance(t, __inst1, __oneInstHash)
		time.Sleep(time.Second)
	}

	println("No.2 sleep end...")
	time.Sleep(time.Second)
	mustNoInstance(t)
}
