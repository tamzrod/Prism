// internal/engine/codec_test.go

package engine

import (
    "encoding/hex"
    "testing"
)

func TestDecode_Bytes(t *testing.T) {
    data := []byte{0x01, 0x02, 0x03}
    out, err := Decode("bytes", data)
    if err != nil {
        t.Fatal(err)
    }
    if got, ok := out.([]byte); !ok || len(got) != len(data) {
        t.Errorf("unexpected output: %v", out)
    }
}

func TestDecode_U8(t *testing.T) {
    data := []byte{0x01, 0xFF}
    out, err := Decode("u8", data)
    if err != nil {
        t.Fatal(err)
    }
    arr := out.([]uint64)
    if arr[0] != 1 || arr[1] != 255 {
        t.Errorf("u8 decode incorrect: %v", arr)
    }
}

func TestDecode_I16Be(t *testing.T) {
    data := []byte{0x01, 0x02}
    out, err := Decode("i16be", data)
    if err != nil {
        t.Fatal(err)
    }
    arr := out.([]int64)
    if arr[0] != 0x0102 {
        t.Errorf("i16be result wrong: %v", arr)
    }
}

func TestEncode_Hex(t *testing.T) {
    data := []byte{0xAA, 0xBB}
    out, err := EncodeResult(data, "hex")
    if err != nil {
        t.Fatal(err)
    }
    if out.(string) != hex.EncodeToString(data) {
        t.Errorf("hex encode wrong: %v", out)
    }
}

func TestEncode_Base64(t *testing.T) {
    data := []byte("hello")
    out, err := EncodeResult(data, "base64")
    if err != nil {
        t.Fatal(err)
    }
    if out.(string) == "" {
        t.Error("empty base64 output")
    }
}
