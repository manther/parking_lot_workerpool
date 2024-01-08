package parkinglot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParkinglotGroupComp(t *testing.T) {
	pl := NewParparkinglot(50, 50)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageGroupParkComp(20)
	for _, err := range err {
		require.NoError(t, err)
	}
	require.Equal(t, 30, pl.CompactAvailable())
}

func TestParkinglotGroupCompTooMuch(t *testing.T) {
	pl := NewParparkinglot(10, 10)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageGroupParkComp(20)
	for _, err := range err {
		require.Error(t, err)
	}
	require.Equal(t, 10, pl.CompactAvailable())
}

func TestParkinglotGroupReg(t *testing.T) {
	pl := NewParparkinglot(50, 50)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageGroupParkReg(20)
	for _, err := range err {
		require.NoError(t, err)
	}
	require.Equal(t, 30, pl.RegularAvailable())
}

func TestParkinglotGroupRegTooMuch(t *testing.T) {
	pl := NewParparkinglot(1, 1)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageGroupParkReg(2)
	for _, err := range err {
		require.Error(t, err)
	}
	require.Equal(t, 1, pl.RegularTot())
	require.Equal(t, 1, pl.RegularAvailable())
}

func TestParkinglotSingleReg(t *testing.T) {
	pl := NewParparkinglot(50, 50)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageSingleParkReg()
	require.NoError(t, err)
	require.Equal(t, 49, pl.RegularAvailable())
}

func TestParkinglotSingleRegTooMuch(t *testing.T) {
	pl := NewParparkinglot(1, 1)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageSingleParkReg()
	require.NoError(t, err)
	require.Equal(t, 0, len(pl.regular))
	err = sp.ManageSingleParkReg()
	require.Error(t, err)
	require.Equal(t, 0, pl.RegularAvailable())
}

func TestParkinglotSingleComp(t *testing.T) {
	pl := NewParparkinglot(50, 50)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageSingleParkComp()
	require.NoError(t, err)
	require.Equal(t, 49, pl.CompactAvailable())
}

func TestParkinglotSingleCompTooMuch(t *testing.T) {
	pl := NewParparkinglot(1, 1)
	sp := NewLotSupervisor(&pl, 8)
	err := sp.ManageSingleParkComp()
	require.NoError(t, err)
	require.Equal(t, 0, pl.CompactAvailable())
	err = sp.ManageSingleParkComp()
	require.Error(t, err)
	require.Equal(t, 0, pl.CompactAvailable())
	require.Equal(t, 1, pl.Total())
	require.Equal(t, false, pl.Full())
}
