package tlog

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestTLog(t *testing.T) {
	writer := NewWriteToConsole()
	//writer := NewWriteToFileSeparate(FileStoreMode(AppendOneFile))
	//writer := NewWriteToFileMixed(FileStoreMode(AppendOneFile))
	tl := New(OmitEmpty(true), TimeFormat(HumanReadableTimeMs), SetWriter(writer), Format(FormatText))
	s1 := `i'm sorry, "cuisw" is right! ohh.\n`
	tl.Debug().Fmt("fmt", "n=%d type=%s v=%v %s", 10, reflect.TypeOf(*tl).String(), *tl, s1).Msg("")
	tl.Debug().Fmt("> ", "n=%d type=%s v=%v %s", 10, reflect.TypeOf(*tl).String(), *tl, s1).Msg("")

	tl.Debug().Str(s1, "val").Msg("")
	tl.Debug().FastStr("str", "val").Msg("")
	tl.Debug().FastStr("emptystr", "").Str("after emptystr", "1").Go()
	tl.Debug().Strs("strs", []string{"val1", "val\tad", s1}).Msg("")

	tl.Info().Bool("bool", false).Msg("info")
	tl.Info().Bools("bools", []bool{true, false}).Msg("info")

	tl.Debug().Int("int", 32).Msg("")
	tl.Debug().Ints("ints", []int{32, 21}).Msg("")
	tl.Debug().Int8("int8", 127).Msg("")
	tl.Debug().Ints8("ints8", []int8{39, 21}).Msg("")
	tl.Debug().Int16("int16", 32767).Msg("")
	tl.Debug().Ints16("ints16", []int16{39, 21}).Msg("")
	tl.Debug().Int32("int32", 65536).Msg("")
	tl.Debug().Ints32("ints32", []int32{39, 21}).Msg("")
	tl.Debug().Int64("int64", 32).Msg("")
	tl.Debug().Ints64("ints64", []int64{32, 21}).Msg("")

	tl.Debug().Uint("uint", 32).Msg("")
	tl.Debug().Uints("uints", []uint{32, 21}).Msg("")
	tl.Debug().Uint8("uint8", 255).Msg("")
	tl.Debug().Uints8("uints8", []uint8{39, 21}).Msg("")
	tl.Debug().Uint16("uint16", 32768).Msg("")
	tl.Debug().Uints16("uints16", []uint16{39, 21}).Msg("")
	tl.Debug().Uint32("uint32", 65536).Msg("")
	tl.Debug().Uints32("uints32", []uint32{39, 21}).Uints32("empty uints32", []uint32{}).Msg("")
	tl.Debug().Uint64("uint64", 32).Msg("")
	tl.Debug().Uints64("uints64", []uint64{32, 21}).Uints64("nil uints64", nil).Msg("")

	tl.Debug().Float32("float32", 1212.32001).Msg("")
	tl.Debug().Floats32("floats32", []float32{32.00001, 21.000012, 1.000012}).Msg("")
	tl.Debug().Float64("float64", 31212.00000000012).Msg("")
	tl.Debug().Floats64("floats64", []float64{32.12121, 21.999999}).Msg("")

	tl.Debug().OmitEmpty(false).Ints("nilints", nil).Ints("emptyints", []int{}).Msg("")

	t1 := 1
	tl.Debug().Type("type", t1).Go()

	tl.Debug().Time("time", time.Now(), time.RFC1123Z).Msg("")

	type js struct {
		Name    string `json:"name"`
		Empty   string `json:"empty,omitempty"`
		Age     int    `json:"age"`
		Address []byte `json:"addr,omitempty"`
	}
	jsdata, _ := json.Marshal(&js{Name: "cuisw", Empty: "", Age: 18})
	tl.Debug().RawJSON("json", jsdata).Msgf("it's %s end #%d", "oooh", 12)
}
