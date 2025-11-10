package internal

import (
	"testing"
	"time"
)

func TestToInt(t *testing.T) {
	if toInt("1") != 1 {
		t.Fatalf("\"1\" != 1")
	}
	if toInt(uint(1)) != 1 {
		t.Fatalf("uint(1) != 1")
	}
	if toInt(int8(1)) != 1 {
		t.Fatalf("int8(1) != 1")
	}
	if toInt(int16(1)) != 1 {
		t.Fatalf("int16(1) != 1")
	}
	if toInt(int32(1)) != 1 {
		t.Fatalf("int32(1) != 1")
	}
	if toInt(int64(1)) != 1 {
		t.Fatalf("int64(1) != 1")
	}
	if toInt(uint8(1)) != 1 {
		t.Fatalf("uint8(1) != 1")
	}
	if toInt(uint16(1)) != 1 {
		t.Fatalf("uint16(1) != 1")
	}
	if toInt(uint32(1)) != 1 {
		t.Fatalf("uint32(1) != 1")
	}
	if toInt(uint64(1)) != 1 {
		t.Fatalf("uint64(1) != 1")
	}
	if toInt(float32(1)) != 1 {
		t.Fatalf("float32(1) != 1")
	}
	if toInt(float64(1)) != 1 {
		t.Fatalf("float64(1) != 1")
	}
	if toInt(true) != 1 {
		t.Fatalf("true != 1")
	}
	if toInt(false) != 0 {
		t.Fatalf("false != 0")
	}
	if toInt(nil) != 0 {
		t.Fatalf("nil != 0")
	}
	if toInt("") != 0 {
		t.Fatalf("'' != 0")
	}
	if toInt(time.Weekday(1)) != 1 {
		t.Fatalf("time.Weekday(1) != 1")
	}
	if toInt(time.Month(1)) != 1 {
		t.Fatalf("time.Month(1) != 1")
	}
}

func TestToFloat64(t *testing.T) {
	if toFloat64("1") != 1.0 {
		t.Fatalf("\"1\" != 1.0")
	}
	if toFloat64(uint(1)) != 1.0 {
		t.Fatalf("uint(1) != 1.0")
	}
	if toFloat64(int8(1)) != 1.0 {
		t.Fatalf("int8(1) != 1.0")
	}
	if toFloat64(int16(1)) != 1.0 {
		t.Fatalf("int16(1) != 1.0")
	}
	if toFloat64(int32(1)) != 1.0 {
		t.Fatalf("int32(1) != 1.0")
	}
	if toFloat64(int64(1)) != 1.0 {
		t.Fatalf("int64(1) != 1.0")
	}
	if toFloat64(uint8(1)) != 1.0 {
		t.Fatalf("uint8(1) != 1.0")
	}
	if toFloat64(uint16(1)) != 1.0 {
		t.Fatalf("uint16(1) != 1.0")
	}
	if toFloat64(uint32(1)) != 1.0 {
		t.Fatalf("uint32(1) != 1.0")
	}
	if toFloat64(uint64(1)) != 1.0 {
		t.Fatalf("uint64(1) != 1.0")
	}
	if toFloat64(float32(1)) != 1.0 {
		t.Fatalf("float32(1) != 1.0")
	}
	if toFloat64(float64(1)) != 1.0 {
		t.Fatalf("float64(1) != 1.0")
	}
	if toFloat64(true) != 1.0 {
		t.Fatalf("true != 1.0")
	}
	if toFloat64(false) != 0.0 {
		t.Fatalf("false != 0.0")
	}
	if toFloat64(nil) != 0.0 {
		t.Fatalf("nil != 0.0")
	}
	if toFloat64("") != 0.0 {
		t.Fatalf("'' != 0.0")
	}
	if toFloat64(time.Weekday(1)) != 1.0 {
		t.Fatalf("time.Weekday(1) != 1.0")
	}
	if toFloat64(time.Month(1)) != 1.0 {
		t.Fatalf("time.Month(1) != 1.0")
	}
}
