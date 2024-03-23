package store

import (
	"context"
	"testing"
	"time"
)

func TestNewDummyStore(t *testing.T) { //nolint:revive
	conf := NewDummyStore()

	defer conf.Close()

	conf.MustGet("")
	conf.Get("")
	conf.Set("", nil)
	conf.SetComment("", "", "")
	conf.SetTag("", nil)
	conf.Remove("")
	conf.RemoveEx("")
	_ = conf.Merge("", nil)
	conf.Has("")
	conf.Locate("")
	_, _ = conf.GetString("", "")
	conf.MustString("", "")
	_, _ = conf.GetStringSlice("", "")
	conf.MustStringSlice("", "")
	_, _ = conf.GetStringMap("", nil)
	conf.MustStringMap("", nil)

	_, _ = conf.GetInt64("", 0)
	conf.MustInt64("", 0)
	_, _ = conf.GetInt32("", 0)
	_ = conf.MustInt32("", 0)
	_, _ = conf.GetInt16("", 0)
	_ = conf.MustInt16("", 0)
	_, _ = conf.GetInt8("", 0)
	_ = conf.MustInt8("", 0)
	_, _ = conf.GetInt("", 0)
	_ = conf.MustInt("", 0)

	_, _ = conf.GetInt64Slice("", 0)
	conf.MustInt64Slice("", 0)
	_, _ = conf.GetInt32Slice("", 0)
	conf.MustInt32Slice("", 0)
	_, _ = conf.GetInt16Slice("", 0)
	conf.MustInt16Slice("", 0)
	_, _ = conf.GetInt8Slice("", 0)
	conf.MustInt8Slice("", 0)
	_, _ = conf.GetIntSlice("", 0)
	conf.MustIntSlice("", 0)

	_, _ = conf.GetInt64Map("", nil)
	conf.MustInt64Map("", nil)
	_, _ = conf.GetInt32Map("", nil)
	conf.MustInt32Map("", nil)
	_, _ = conf.GetInt16Map("", nil)
	conf.MustInt16Map("", nil)
	_, _ = conf.GetInt8Map("", nil)
	conf.MustInt8Map("", nil)
	_, _ = conf.GetIntMap("", nil)
	conf.MustIntMap("", nil)

	_, _ = conf.GetUint64("", 0)
	conf.MustUint64("", 0)
	_, _ = conf.GetUint32("", 0)
	conf.MustUint32("", 0)
	_, _ = conf.GetUint16("", 0)
	conf.MustUint16("", 0)
	_, _ = conf.GetUint8("", 0)
	conf.MustUint8("", 0)
	_, _ = conf.GetUint("", 0)
	conf.MustUint("", 0)

	_, _ = conf.GetUint64Slice("", 0)
	conf.MustUint64Slice("", 0)
	_, _ = conf.GetUint32Slice("", 0)
	conf.MustUint32Slice("", 0)
	_, _ = conf.GetUint16Slice("", 0)
	conf.MustUint16Slice("", 0)
	_, _ = conf.GetUint8Slice("", 0)
	conf.MustUint8Slice("", 0)
	_, _ = conf.GetUintSlice("", 0)
	conf.MustUintSlice("", 0)

	_, _ = conf.GetUint64Map("", nil)
	conf.MustUint64Map("", nil)
	_, _ = conf.GetUint32Map("", nil)
	conf.MustUint32Map("", nil)
	_, _ = conf.GetUint16Map("", nil)
	conf.MustUint16Map("", nil)
	_, _ = conf.GetUint8Map("", nil)
	conf.MustUint8Map("", nil)
	_, _ = conf.GetUintMap("", nil)
	conf.MustUintMap("", nil)

	_, _ = conf.GetKibiBytes("", 0)
	conf.MustKibiBytes("", 0)
	_, _ = conf.GetKiloBytes("", 0)
	conf.MustKiloBytes("", 0)

	_, _ = conf.GetFloat64("", 0.0)
	conf.MustFloat64("", 0.0)
	_, _ = conf.GetFloat64Slice("", 0.0)
	conf.MustFloat64Slice("", 0.0)
	_, _ = conf.GetFloat64Map("", nil)
	conf.MustFloat64Map("", nil)

	_, _ = conf.GetFloat32("", 0.0)
	conf.MustFloat32("", 0.0)
	_, _ = conf.GetFloat32Slice("", 0.0)
	conf.MustFloat32Slice("", 0.0)
	_, _ = conf.GetFloat32Map("", nil)
	conf.MustFloat32Map("", nil)

	_, _ = conf.GetComplex64("", 0.0)
	conf.MustComplex64("", 0.0)
	_, _ = conf.GetComplex64Slice("", 0.0)
	conf.MustComplex64Slice("", 0.0)
	_, _ = conf.GetComplex64Map("", nil)
	conf.MustComplex64Map("", nil)

	_, _ = conf.GetComplex128("", 0.0)
	conf.MustComplex128("", 0.0)
	_, _ = conf.GetComplex128Slice("", 0.0)
	conf.MustComplex128Slice("", 0.0)
	_, _ = conf.GetComplex128Map("", nil)
	conf.MustComplex128Map("", nil)

	_, _ = conf.GetBool("", false)
	conf.MustBool("", false)
	_, _ = conf.GetBoolSlice("", false)
	conf.MustBoolSlice("", false)
	_, _ = conf.GetBoolMap("", nil)
	conf.MustBoolMap("", nil)

	_, _ = conf.GetDuration("", 0)
	conf.MustDuration("", 0)
	_, _ = conf.GetDurationSlice("", 0)
	conf.MustDurationSlice("", 0)
	_, _ = conf.GetDurationMap("", nil)
	conf.MustDurationMap("", nil)

	now := time.Now()
	_, _ = conf.GetTime("", now)
	conf.MustTime("", now)
	_, _ = conf.GetTimeSlice("", now)
	conf.MustTimeSlice("", now)
	_, _ = conf.GetTimeMap("", nil)
	conf.MustTimeMap("", nil)

	_, _ = conf.GetR("", nil)
	conf.MustR("", nil)

	_, _ = conf.GetM("", nil)
	conf.MustM("", nil)
	_ = conf.GetSectionFrom("", nil, nil)

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
	_, _ = conf.Load(context.TODO(), nil)
	conf.WithinLoading(func() {})
}
