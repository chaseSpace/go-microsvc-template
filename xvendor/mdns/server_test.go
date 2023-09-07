// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MIT

package mdns

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestServer_StartStop(t *testing.T) {
	s := makeService(t)
	serv, err := NewServer(&Config{Zone: s})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(time.Second * 30)

	// waiting, you can execute `dig @224.0.0.251 -p 5353 _http._tcp.local PTR` to see the result
	if err := serv.Shutdown(); err != nil {
		t.Fatalf("err: %v", err)
	}
}

// In early Windows 10 versions, this test fails (always timeout), you must
// install bonjour service to pass the test.
// - bonjour download: https://support.apple.com/kb/DL999?locale=en_US
// Checkout reason: https://web.archive.org/web/20220807151616/https://qiita.com/maccadoo/items/48ace84f8aca030a12f1
func TestServer_Lookup(t *testing.T) {
	instance := "inst-007"
	service := "_http._tcp"
	domain := "local."
	txtInfo := []string{"Local web server"}
	serv, err := NewServer(&Config{Zone: makeServiceWithServiceName(t, instance, service, domain, txtInfo)})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer func() {
		if err := serv.Shutdown(); err != nil {
			t.Fatalf("err: %v", err)
		}
	}()

	entries := make(chan *ServiceEntry, 1)

	errCh := make(chan error, 1)
	timeout := time.Millisecond * 50
	go func() {
		select {
		case e := <-entries:
			if e == nil {
				errCh <- fmt.Errorf("Entry nil")
				return
			}
			if e.Name != strings.Join([]string{instance, service, domain}, ".") {
				errCh <- fmt.Errorf("Entry has the wrong name: %+v\n", e)
				return
			}
			if e.Port != 80 {
				errCh <- fmt.Errorf("Entry has the wrong port: %+v\n", e)
				return
			}
			if e.Info != txtInfo[0] {
				errCh <- fmt.Errorf("Entry as the wrong Info: %+v\n", e)
				return
			}
			errCh <- nil
		case <-time.After(timeout):
			errCh <- fmt.Errorf("Timed out waiting for response")
		}
	}()

	params := &QueryParam{
		Service:     service,
		Domain:      domain,
		Timeout:     timeout,
		Entries:     entries,
		DisableIPv6: true,
		LogLevel:    LogLevelInfo,
	}
	err = Query(params)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = <-errCh
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}
