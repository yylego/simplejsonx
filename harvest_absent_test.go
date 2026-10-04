package simplejsonx_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/simplejsonx"
)

func TestHarvestAbsent(t *testing.T) {
	object, err := simplejsonx.Load([]byte(`{"null":null,"zero":0,"text":"12","bad":"oops"}`))
	require.NoError(t, err)
	opt := &simplejsonx.HarvestOption[int]{AbsentOk: true, Strconv: strconv.Atoi}
	value, err := simplejsonx.Harvest(object, "absent", opt)
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest(object, "null", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest[int](nil, "absent", opt)
	require.Error(t, err)
	opt.Nullable = true
	value, err = simplejsonx.Harvest(object, "null", opt)
	require.NoError(t, err)
	require.Zero(t, value)
	value, err = simplejsonx.Harvest(object, "text", opt)
	require.NoError(t, err)
	require.Equal(t, 12, value)
	_, err = simplejsonx.Harvest(object, "bad", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest(object, "", opt)
	require.Error(t, err)
	_, err = simplejsonx.Harvest[int](object, "absent", nil)
	require.Error(t, err)
	value, err = simplejsonx.Harvest[int](object, "zero", nil)
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest[int](nil, "absent", nil)
	require.Error(t, err)
	value, err = simplejsonx.Harvest[int](nil, "absent", &simplejsonx.HarvestOption[int]{Nullable: true})
	require.NoError(t, err)
	require.Zero(t, value)
	_, err = simplejsonx.Harvest[int](object, "absent", &simplejsonx.HarvestOption[int]{Nullable: true})
	require.Error(t, err)
}
