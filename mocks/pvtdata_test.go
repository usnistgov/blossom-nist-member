package mocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPvtData(t *testing.T) {
	pvtData := NewPvtData()
	pvtData.CreateNewCollection("coll1", []string{"r1", "r2"}, []string{"w1"})

	err := pvtData.PutPrivateData("r1", "coll1", "key", []byte("value"))
	require.Error(t, err)

	err = pvtData.PutPrivateData("w1", "coll1", "key", []byte("value"))
	require.NoError(t, err)

	bytes, err := pvtData.GetPrivateData("w1", "coll1", "key")
	require.NoError(t, err)
	require.Equal(t, []byte("value"), bytes)

	bytes, err = pvtData.GetPrivateData("r1", "coll1", "key")
	require.NoError(t, err)
	require.Equal(t, []byte("value"), bytes)

	pvtData.CreateNewCollection("coll2", []string{"r1", "r2"}, []string{"w1"})

	err = pvtData.PutPrivateData("r1", "coll2", "key", []byte("value"))
	require.Error(t, err)

	err = pvtData.PutPrivateData("w1", "coll2", "key", []byte("value"))
	require.NoError(t, err)

	bytes, err = pvtData.GetPrivateData("w1", "coll2", "key")
	require.NoError(t, err)
	require.Equal(t, []byte("value"), bytes)

	bytes, err = pvtData.GetPrivateData("r1", "coll2", "key")
	require.NoError(t, err)
	require.Equal(t, []byte("value"), bytes)
}
