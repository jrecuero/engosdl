package engosdl_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jrecuero/engosdl"
)

func TestEvent_PoolManager(t *testing.T) {
	h := engosdl.NewEventManager("test-event-manager")
	if h == nil {
		t.Error("error creating event manager")
	}
	poolID, err := h.CreatePool("test-pool-1")
	if err != nil {
		t.Errorf("error creating pool: %s", err.Error())
	}
	pool := h.GetPool(poolID)
	if pool == nil {
		t.Errorf("error getting pool")
	}
	pools := h.GetPools()
	if pools == nil {
		t.Errorf("error getting pools")
	}
	if len(pools) != 1 {
		t.Errorf("error getting number of pools\nexp: %d\ngot: %d\n", 1, len(pools))
	}
	_, err = h.CreatePool("test-pool-2")
	if err != nil {
		t.Errorf("error creating pool: %s", err.Error())
	}
	if len(pools) != 2 {
		t.Errorf("error getting number of pools\nexp: %d\ngot: %d\n", 2, len(pools))
	}
	id, err := h.GetIDForName("test-pool-1")
	if err != nil {
		t.Errorf("error getting id from name: %s", err.Error())
	}
	if id != poolID {
		t.Errorf("error getting id from name\nexp: %s\ngot: %s\n", poolID, id)
	}
	if err := h.DeletePool(poolID); err != nil {
		t.Errorf("error deleting pool: %s", err.Error())
	}
}

func TestEvent_EventPool(t *testing.T) {
	h := engosdl.NewEventManager("test-event-manager")
	poolID, err := h.CreatePool("test-pool-1")
	if err != nil {
		t.Errorf("error creating pool: %s", err.Error())
	}
	for i := 0; i < 10; i++ {
		eventName := fmt.Sprintf("event-%s", strconv.Itoa(i))
		dataName := fmt.Sprintf("data-%s", strconv.Itoa(i))
		err = h.GetPool(poolID).Add(engosdl.NewEvent(eventName, engosdl.NewObject(dataName)))
		if err != nil {
			t.Errorf("error adding event to event pool: %s", err.Error())
		}
	}
	pool := h.GetPool(poolID)
	if pool == nil {
		t.Errorf("error getting pool")
	}
	if len(pool.Pool()) != 10 {
		t.Errorf("error length of pool\nexp: %d\ngot: %d\n", 10, len(pool.Pool()))
	}
	for i := 0; i < 10; i++ {
		event, err := pool.Next()
		if err != nil {
			t.Errorf("error getting next entry in pool: %s", err.Error())
		}
		if event == nil {
			t.Errorf("error getting next entry in pool: nil")
		}
		got := event.GetData().GetName()
		exp := fmt.Sprintf("data-%s", strconv.Itoa(i))
		if got != exp {
			t.Errorf("error getting next entry in pool\nexp: %s\ngot: %s\n", exp, got)
		}
		event, err = pool.Pop()
		if err != nil {
			t.Errorf("error pop entry in pool: %s", err.Error())
		}
		if event == nil {
			t.Errorf("error pop entry in pool: nil")
		}
		got = event.GetData().GetName()
		exp = fmt.Sprintf("data-%s", strconv.Itoa(i))
		if got != exp {
			t.Errorf("error pop entry in pool\nexp: %s\ngot: %s\n", exp, got)
		}
	}
	if len(pool.Pool()) != 0 {
		t.Errorf("error length of pool\nexp: %d\ngot: %d\n", 0, len(pool.Pool()))
	}
	for i := 0; i < 10; i++ {
		eventName := fmt.Sprintf("event-%s", strconv.Itoa(i))
		dataName := fmt.Sprintf("data-%s", strconv.Itoa(i))
		err = h.GetPool(poolID).Add(engosdl.NewEvent(eventName, engosdl.NewObject(dataName)))
		if err != nil {
			t.Errorf("error adding event to event pool: %s", err.Error())
		}
	}
	if len(pool.Pool()) != 10 {
		t.Errorf("error length of pool\nexp: %d\ngot: %d\n", 10, len(pool.Pool()))
	}
	if err = pool.Flush(); err != nil {
		t.Errorf("error flush event pool: %s", err.Error())
	}
	if len(pool.Pool()) != 0 {
		t.Errorf("error length of pool after flush\nexp: %d\ngot: %d\n", 0, len(pool.Pool()))
	}
}
