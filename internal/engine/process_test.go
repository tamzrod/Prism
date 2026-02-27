// internal/engine/process_test.go

package engine

import (
    "encoding/json"
    "testing"

    "github.com/tamzrod/prism/internal/errors"
    "github.com/tamzrod/prism/internal/protocol"
)

func TestProcess_Success(t *testing.T) {
    payload := []byte{0x01, 0x02}
    req := &protocol.Request{
        PayloadFormat: "u16be",
        TargetFormat:  "hex",
        PayloadBytes:  &payload,
    }
    respBytes, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(respBytes, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if _, ok := res["value_hex"]; !ok {
        t.Errorf("expected value_hex in response: %v", res)
    }
}

func TestProcess_FromNumericPayload(t *testing.T) {
    // convert an array of u8 values to base64
    arr := []uint64{0x41, 0x42, 0x43}
    req := &protocol.Request{
        PayloadFormat: "u8",
        TargetFormat:  "base64",
        PayloadU8:     &arr,
    }
    respBytes, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(respBytes, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if _, ok := res["value_base64"]; !ok {
        t.Errorf("expected value_base64 in response: %v", res)
    }
}

func TestProcess_HexPayloadString(t *testing.T) {
    // bytes type can also be provided as hex string
    hexstr := "deadbeef"
    req := &protocol.Request{
        PayloadFormat: "bytes",
        TargetFormat:  "u16be",
        PayloadHex:    &hexstr,
    }
    respBytes, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(respBytes, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if _, ok := res["value_u16be"]; !ok {
        t.Errorf("expected value_u16be in response: %v", res)
    }
}

func TestProcess_Unix32Conversion(t *testing.T) {
    // a single 32-bit unix timestamp -> u32le
    payload := []byte{0x5F, 0x5E, 0x10, 0x00} // little endian
    req := &protocol.Request{
        PayloadFormat: "unix32le",
        TargetFormat:  "u32le",
        PayloadBytes:  &payload,
    }
    resp, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(resp, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if vals, ok := res["value_u32le"]; !ok {
        t.Errorf("expected value_u32le, got %v", res)
    } else {
        arr := vals.([]interface{})
        if len(arr) != 1 || uint64(arr[0].(float64)) != 0x00105E5F {
            t.Errorf("unexpected conversion result: %v", arr)
        }
    }
}

func TestProcess_InvalidFormat(t *testing.T) {
    payload := []byte{0x00}
    req := &protocol.Request{
        PayloadFormat: "bogus",
        TargetFormat:  "u8",
        PayloadBytes:  &payload,
    }
    respBytes, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    json.Unmarshal(respBytes, &res)
    if code, ok := res["error_code"]; !ok || code.(float64) != float64(errors.CodeUnknownFormat) {
        t.Errorf("expected unknown format error, got %v", res)
    }
}

func TestProcess_PayloadFormatHexAlias(t *testing.T) {
    // same data used earlier but with payload_format = "hex"
    hexstr := "42f6e666"
    req := &protocol.Request{
        PayloadFormat: "hex",
        TargetFormat:  "u16be",
        PayloadHex:    &hexstr,
    }
    resp, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(resp, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if vals, ok := res["value_u16be"]; !ok {
        t.Errorf("expected value_u16be, got %v", res)
    } else {
        arr := vals.([]interface{})
        if len(arr) != 2 || uint64(arr[0].(float64)) != 0x42f6 || uint64(arr[1].(float64)) != 0xe666 {
            t.Errorf("unexpected conversion: %v", arr)
        }
    }
}

func TestProcess_Unix32ToRFC3339(t *testing.T) {
    // choose a known timestamp, e.g. 0x5E2D0BDC = 1581231230 -> 2020-02-09T...
    payload := []uint64{1581231230}
    req := &protocol.Request{
        PayloadFormat: "unix32be",
        TargetFormat:  "rfc3339",
        PayloadU32:    &payload,
    }
    resp, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(resp, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if str, ok := res["value_rfc3339"]; !ok {
        t.Errorf("expected rfc3339 value, got %v", res)
    } else {
        // just ensure it's a nonempty string
        if str.(string) == "" {
            t.Error("empty rfc3339 output")
        }
    }
}

func TestProcess_Unix64ToRFC3339(t *testing.T) {
    payload := []uint64{1581231230}
    req := &protocol.Request{
        PayloadFormat: "unix64le",
        TargetFormat:  "rfc3339",
        PayloadU64:    &payload,
    }
    resp, err := Process(req)
    if err != nil {
        t.Fatalf("process error: %v", err)
    }
    var res map[string]interface{}
    if err := json.Unmarshal(resp, &res); err != nil {
        t.Fatalf("unmarshal response: %v", err)
    }
    if str, ok := res["value_rfc3339"]; !ok {
        t.Errorf("expected rfc3339 value, got %v", res)
    } else {
        if str.(string) == "" {
            t.Error("empty rfc3339 output")
        }
    }
}
