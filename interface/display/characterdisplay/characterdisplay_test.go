package characterdisplay

import (
	"reflect"
	"testing"
	"time"
)

const (
	rows = 4
	cols = 20
)

type mockController struct {
	calls chan call
}

type call struct {
	name      string
	arguments []interface{}
}

func noArgCall(name string) call {
	return call{name, []interface{}{}}
}

func (mock *mockController) DisplayOff() error   { mock.calls <- noArgCall("DisplayOff"); return nil }
func (mock *mockController) DisplayOn() error    { mock.calls <- noArgCall("DisplayOn"); return nil }
func (mock *mockController) CursorOff() error    { mock.calls <- noArgCall("CursorOff"); return nil }
func (mock *mockController) CursorOn() error     { mock.calls <- noArgCall("CursorOn"); return nil }
func (mock *mockController) BlinkOff() error     { mock.calls <- noArgCall("BlinkOff"); return nil }
func (mock *mockController) BlinkOn() error      { mock.calls <- noArgCall("BlinkOn"); return nil }
func (mock *mockController) ShiftLeft() error    { mock.calls <- noArgCall("ShiftLeft"); return nil }
func (mock *mockController) ShiftRight() error   { mock.calls <- noArgCall("ShiftRight"); return nil }
func (mock *mockController) BacklightOff() error { mock.calls <- noArgCall("BacklightOff"); return nil }
func (mock *mockController) BacklightOn() error  { mock.calls <- noArgCall("BacklightOn"); return nil }
func (mock *mockController) Home() error         { mock.calls <- noArgCall("Home"); return nil }
func (mock *mockController) Clear() error        { mock.calls <- noArgCall("Clear"); return nil }
func (mock *mockController) Close() error        { mock.calls <- noArgCall("Close"); return nil }
func (mock *mockController) WriteChar(b byte) error {
	mock.calls <- call{"WriteChar", []interface{}{b}}
	return nil
}
func (mock *mockController) SetCursor(col, row int) error {
	mock.calls <- call{"SetCursor", []interface{}{col, row}}
	return nil
}

func (mock *mockController) testExpectedCalls(expectedCalls []call, t *testing.T) {
	for _, expectedCall := range expectedCalls {
		select {
		case actualCall := <-mock.calls:
			if !reflect.DeepEqual(expectedCall, actualCall) {
				t.Errorf("Expected call %+v, actual call %+v", expectedCall, actualCall)
			}
		case <-time.After(time.Millisecond * 1):
			t.Errorf("Timeout reading next call. Expected call %+v", expectedCall)
		}
	}

ExtraCallsCheck:
	for {
		select {
		case extraCall := <-mock.calls:
			t.Errorf("Unexpected call %+v", extraCall)
		case <-time.After(time.Millisecond * 1):
			break ExtraCallsCheck
		}
	}
}

func newMockController() *mockController {
	return &mockController{make(chan call, 256)}
}

func TestNewline(t *testing.T) {
	mock := newMockController()
	disp := New(mock, cols, rows)
	disp.Newline()

	expectedCalls := []call{
		call{"SetCursor", []interface{}{0, 1}},
	}

	mock.testExpectedCalls(expectedCalls, t)
}

func TestMessage(t *testing.T) {
	mock := newMockController()
	disp := New(mock, cols, rows)
	disp.Message("ab")

	expectedCalls := []call{
		call{"WriteChar", []interface{}{byte('a')}},
		call{"WriteChar", []interface{}{byte('b')}},
	}

	mock.testExpectedCalls(expectedCalls, t)
}

func TestMessage_newLine(t *testing.T) {
	mock := newMockController()
	disp := New(mock, cols, rows)
	disp.Message("a\nb")

	expectedCalls := []call{
		call{"WriteChar", []interface{}{byte('a')}},
		call{"SetCursor", []interface{}{0, 1}},
		call{"WriteChar", []interface{}{byte('b')}},
	}

	mock.testExpectedCalls(expectedCalls, t)
}

func TestMessage_wrap(t *testing.T) {
	mock := newMockController()
	disp := New(mock, cols, rows)
	disp.SetCursor(cols-1, 0)
	disp.Message("ab")

	expectedCalls := []call{
		call{"SetCursor", []interface{}{cols - 1, 0}},
		call{"WriteChar", []interface{}{byte('a')}},
		call{"SetCursor", []interface{}{0, 1}},
		call{"WriteChar", []interface{}{byte('b')}},
	}

	mock.testExpectedCalls(expectedCalls, t)
}
