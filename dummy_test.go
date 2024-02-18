package store

import (
	"context"
	"testing"
	"time"
)

func TestNewDummyStore(t *testing.T) {
	conf := NewDummyStore()

	defer conf.Close()

	conf.MustGet("")
	conf.Get("")
	conf.Set("", nil)
	conf.SetComment("", "", "")
	conf.SetTag("", nil)
	conf.Remove("")
	conf.RemoveEx("")
	conf.Merge("", nil)
	conf.Has("")
	conf.Locate("")
	conf.GetString("", "")
	conf.MustString("", "")
	conf.GetStringSlice("", "")
	conf.MustStringSlice("", "")
	conf.GetStringMap("", nil)
	conf.MustStringMap("", nil)

	conf.GetInt64("", 0)
	conf.MustInt64("", 0)
	conf.GetInt32("", 0)
	conf.MustInt32("", 0)
	conf.GetInt16("", 0)
	conf.MustInt16("", 0)
	conf.GetInt8("", 0)
	conf.MustInt8("", 0)
	conf.GetInt("", 0)
	conf.MustInt("", 0)

	conf.GetInt64Slice("", 0)
	conf.MustInt64Slice("", 0)
	conf.GetInt32Slice("", 0)
	conf.MustInt32Slice("", 0)
	conf.GetInt16Slice("", 0)
	conf.MustInt16Slice("", 0)
	conf.GetInt8Slice("", 0)
	conf.MustInt8Slice("", 0)
	conf.GetIntSlice("", 0)
	conf.MustIntSlice("", 0)

	conf.GetInt64Map("", nil)
	conf.MustInt64Map("", nil)
	conf.GetInt32Map("", nil)
	conf.MustInt32Map("", nil)
	conf.GetInt16Map("", nil)
	conf.MustInt16Map("", nil)
	conf.GetInt8Map("", nil)
	conf.MustInt8Map("", nil)
	conf.GetIntMap("", nil)
	conf.MustIntMap("", nil)

	conf.GetUint64("", 0)
	conf.MustUint64("", 0)
	conf.GetUint32("", 0)
	conf.MustUint32("", 0)
	conf.GetUint16("", 0)
	conf.MustUint16("", 0)
	conf.GetUint8("", 0)
	conf.MustUint8("", 0)
	conf.GetUint("", 0)
	conf.MustUint("", 0)

	conf.GetUint64Slice("", 0)
	conf.MustUint64Slice("", 0)
	conf.GetUint32Slice("", 0)
	conf.MustUint32Slice("", 0)
	conf.GetUint16Slice("", 0)
	conf.MustUint16Slice("", 0)
	conf.GetUint8Slice("", 0)
	conf.MustUint8Slice("", 0)
	conf.GetUintSlice("", 0)
	conf.MustUintSlice("", 0)

	conf.GetUint64Map("", nil)
	conf.MustUint64Map("", nil)
	conf.GetUint32Map("", nil)
	conf.MustUint32Map("", nil)
	conf.GetUint16Map("", nil)
	conf.MustUint16Map("", nil)
	conf.GetUint8Map("", nil)
	conf.MustUint8Map("", nil)
	conf.GetUintMap("", nil)
	conf.MustUintMap("", nil)

	conf.GetKibiBytes("", 0)
	conf.MustKibiBytes("", 0)
	conf.GetKiloBytes("", 0)
	conf.MustKiloBytes("", 0)

	conf.GetFloat64("", 0.0)
	conf.MustFloat64("", 0.0)
	conf.GetFloat64Slice("", 0.0)
	conf.MustFloat64Slice("", 0.0)
	conf.GetFloat64Map("", nil)
	conf.MustFloat64Map("", nil)

	conf.GetFloat32("", 0.0)
	conf.MustFloat32("", 0.0)
	conf.GetFloat32Slice("", 0.0)
	conf.MustFloat32Slice("", 0.0)
	conf.GetFloat32Map("", nil)
	conf.MustFloat32Map("", nil)

	conf.GetComplex64("", 0.0)
	conf.MustComplex64("", 0.0)
	conf.GetComplex64Slice("", 0.0)
	conf.MustComplex64Slice("", 0.0)
	conf.GetComplex64Map("", nil)
	conf.MustComplex64Map("", nil)

	conf.GetComplex128("", 0.0)
	conf.MustComplex128("", 0.0)
	conf.GetComplex128Slice("", 0.0)
	conf.MustComplex128Slice("", 0.0)
	conf.GetComplex128Map("", nil)
	conf.MustComplex128Map("", nil)

	conf.GetBool("", false)
	conf.MustBool("", false)
	conf.GetBoolSlice("", false)
	conf.MustBoolSlice("", false)
	conf.GetBoolMap("", nil)
	conf.MustBoolMap("", nil)

	conf.GetDuration("", 0)
	conf.MustDuration("", 0)
	conf.GetDurationSlice("", 0)
	conf.MustDurationSlice("", 0)
	conf.GetDurationMap("", nil)
	conf.MustDurationMap("", nil)

	now := time.Now()
	conf.GetTime("", now)
	conf.MustTime("", now)
	conf.GetTimeSlice("", now)
	conf.MustTimeSlice("", now)
	conf.GetTimeMap("", nil)
	conf.MustTimeMap("", nil)

	conf.GetR("", nil)
	conf.MustR("", nil)

	conf.GetM("", nil)
	conf.MustM("", nil)
	conf.GetSectionFrom("", nil, nil)

	conf.Dump()
	conf.Clone()
	conf.Dup()
	conf.Walk("", nil)
	conf.WithPrefix("")
	conf.WithPrefixReplaced("")
	conf.SetPrefix("")
	conf.Prefix()
	conf.Delimiter()
	conf.SetDelimiter('\t')
	conf.Load(context.TODO(), nil)
	conf.WithinLoading(func() {})
}
