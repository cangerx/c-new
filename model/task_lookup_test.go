package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetByTaskIdForAdminIsNotOwnerScoped(t *testing.T) {
	truncateTables(t)
	insertTask(t, &Task{TaskID: "task_admin_lookup", UserId: 42})

	_, exists, err := GetByTaskId(7, "task_admin_lookup")
	require.NoError(t, err)
	require.False(t, exists)

	task, exists, err := GetByTaskIdForAdmin("task_admin_lookup")
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, 42, task.UserId)
}
