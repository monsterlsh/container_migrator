package util

import "testing"

func TestGetDir(t *testing.T) {
	GetCheckpointDir("/opt/migrator/client/myredis")
}
