package engosdl_test

import (
	"fmt"
	"testing"

	"github.com/jrecuero/engosdl"
)

var TEST_RESULTS []string = []string{}

func Benchmark_Delegate(b *testing.B) {
	for i := 0; i < 10000; i++ {
		TEST_RESULTS = []string{}
		h := engosdl.NewDelegateHandler("test-handler")
		obj := engosdl.NewObject("test-object")
		delegate := h.CreateDelegate(obj, "active")
		h.RegisterToDelegate(obj, delegate, test_create_register)
		h.TriggerDelegate(delegate, true)
		if len(TEST_RESULTS) != 1 {
			b.Errorf("Trigger Delegate error method not called")
		}
		if len(TEST_RESULTS) == 1 && TEST_RESULTS[0] != "signature was called" {
			b.Errorf("Trigger Delegate error method not called")
		}
	}
}

func TestDelegate_CreateDelegate(t *testing.T) {
	obj := engosdl.NewObject("test-object")
	delegate := engosdl.NewDelegate("test-delegate", obj, "active")
	if delegate.GetObject() != obj {
		t.Errorf("new delegate object error\nexp %#v\ngot: %#v", obj, delegate.GetObject())
	}
	if delegate.GetEventName() != "active" {
		t.Errorf("new delegate event error\nexp: %#v\ngot: %#v", "active", delegate.GetEventName())
	}
}

func test_create_register(...interface{}) bool {
	TEST_RESULTS = append(TEST_RESULTS, "signature was called")
	return true
}

func TestDelegate_CreateRegister(t *testing.T) {
	obj := engosdl.NewObject("test-object")
	delegate := engosdl.NewDelegate("test-delegate", obj, "active")
	register := engosdl.NewRegister("register-test", obj, nil, nil, delegate, test_create_register)
	if register.GetDelegate() != delegate {
		t.Errorf("new register delegate error\nexp: %#v\ngot: %#v", delegate, register.GetDelegate())
	}
	if fmt.Sprintf("%p", register.GetSignature()) != fmt.Sprintf("%p", test_create_register) {
		t.Errorf("mew register signature error\nexp: %p\ngot %p", test_create_register, register.GetSignature())
	}
}

func TestDelegate_DelegateHandler(t *testing.T) {
	TEST_RESULTS = []string{}
	h := engosdl.NewDelegateHandler("test-handler")
	obj := engosdl.NewObject("test-object")
	delegate := h.CreateDelegate(obj, "active")
	h.RegisterToDelegate(obj, delegate, test_create_register)
	h.TriggerDelegate(delegate, true)
	if len(TEST_RESULTS) != 1 {
		t.Errorf("Trigger Delegate error method not called")
	}
	if len(TEST_RESULTS) == 1 && TEST_RESULTS[0] != "signature was called" {
		t.Errorf("Trigger Delegate error method not called")
	}
}
