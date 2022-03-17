package env

import (
	"reflect"
	"testing"
	"time"

	"github.com/iver-wharf/wharf-core/v2/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBind[T BindConstraint](t *testing.T, ptr T, envKey string, envValue string, want any) {
	t.Run(envKey, func(t *testing.T) {
		testutil.SetEnv(t, envKey, envValue)
		require.NoError(t, Bind(ptr, envKey))
		// dereference the pointer inside the interface,
		// turning e.g. interface{*int} into interface{int}
		got := reflect.ValueOf(ptr).Elem().Interface()
		assert.Equal(t, want, got)
	})
}

func TestBind(t *testing.T) {
	var (
		myString   string
		myBool     bool
		myInt      int
		myInt32    int32
		myInt64    int64
		myUint     uint
		myUint32   uint32
		myUint64   uint64
		myFloat32  float32
		myFloat64  float64
		myDuration time.Duration
	)
	testBind(t, &myString, "MY_STR", "bar", "bar")
	testBind(t, &myBool, "MY_BOOL", "true", true)
	testBind(t, &myInt, "MY_INT", "-123", int(-123))
	testBind(t, &myInt32, "MY_INT32", "-123", int32(-123))
	testBind(t, &myInt64, "MY_INT64", "-123", int64(-123))
	testBind(t, &myUint, "MY_UINT", "123", uint(123))
	testBind(t, &myUint32, "MY_UINT32", "123", uint32(123))
	testBind(t, &myUint64, "MY_UINT64", "123", uint64(123))
	testBind(t, &myFloat32, "MY_FLOAT32", "123.0", float32(123.0))
	testBind(t, &myFloat64, "MY_FLOAT64", "123.0", float64(123.0))
	testBind(t, &myDuration, "MY_DURATION", "5s", 5*time.Second)
}

func TestBindMultiple_noErrorOnNilMap(t *testing.T) {
	assert.NoError(t, BindMultiple[*int](nil))
}
