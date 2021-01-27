package kafka

import (
	"bufio"
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"
)

func TestFindCoordinatorResponseV0(t *testing.T) {
	item := findCoordinatorResponseV0{
		ErrorCode: 2,
		Coordinator: findCoordinatorResponseCoordinatorV0{
			NodeID: 3,
			Host:   "b",
			Port:   4,
		},
	}

	b := bytes.NewBuffer(nil)
	w := &writeBuffer{w: b}
	item.writeTo(w)

	var found findCoordinatorResponseV0
	remain, err := (&found).readFrom(bufio.NewReader(b), b.Len())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if remain != 0 {
		t.Errorf("expected 0 remain, got %v", remain)
		t.FailNow()
	}
	if !reflect.DeepEqual(item, found) {
		t.Error("expected item and found to be the same")
		t.FailNow()
	}
}

func TestClientFindCoordinator(t *testing.T) {
	client, shutdown := newLocalClient()
	defer shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.WaitForCoordinatorIndefinitely(ctx, &FindCoordinatorRequest{
		Addr:    client.Addr,
		Key:     "TransactionalID-1",
		KeyType: CoordinatorKeyTypeTransaction,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.Coordinator.Host != "localhost" {
		t.Fatal("Coordinator should be found @ localhost")
	}
}
