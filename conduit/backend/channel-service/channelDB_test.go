package main

import (
	"backend/channel-service/database"
	"testing"
)

func TestChannelDBLogin(t *testing.T) {

	pool, err := database.GetChannelsDBConnPool()
	defer pool.Close()
	if err != nil {
		t.Errorf("There was error connecting " + err.Error())
		t.Failed()
	}
}
