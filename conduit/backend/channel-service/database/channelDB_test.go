package database

import (
	"testing"
)

func TestChannelDBLogin(t *testing.T) {

	pool, err := GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		t.Errorf("There was error connecting " + err.Error())
		t.Failed()
	}
}
