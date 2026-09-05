package relay

import (
	"encoding/json"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/require"
)

func TestTaskModel2LogDtoHidesUpstreamLocations(t *testing.T) {
	task := &model.Task{
		TaskID:     "task_public",
		Action:     constant.TaskActionTextToVideo,
		Status:     model.TaskStatusSuccess,
		FailReason: "legacy result https://origin.example.com/video.mp4",
		Properties: model.Properties{Input: `{"prompt":"keep request data","image":"https://user.example.com/input.png"}`},
		PrivateData: model.TaskPrivateData{
			ResultURL: "https://origin.example.com/video.mp4",
		},
		Data: json.RawMessage(`{"status":"succeeded","video_url":"https://cdn.example.com/video.mp4","download_url":"cdn.internal/private.mp4","diagnostic":"fetch https://origin.example.com/private"}`),
	}

	result := TaskModel2LogDto(task)
	require.Equal(t, "/v1/videos/task_public/content", result.ResultURL)
	require.NotContains(t, result.FailReason, "origin.example.com")
	require.NotContains(t, string(result.Data), "cdn.example.com")
	require.NotContains(t, string(result.Data), "cdn.internal")
	require.NotContains(t, string(result.Data), "origin.example.com")
	require.Contains(t, string(result.Data), "[hidden upstream URL]")
	require.Contains(t, result.Properties.(model.Properties).Input, "user.example.com")
}

func TestTaskModel2LogDtoPreservesNonURLResponseFields(t *testing.T) {
	task := &model.Task{
		TaskID: "task_public",
		Action: constant.TaskActionTextToVideo,
		Status: model.TaskStatusSuccess,
		Data: json.RawMessage(`{
			"id":"upstream-private",
			"task_id":"upstream-task-private",
			"model":"video-model",
			"duration":10,
			"billing_amount":0.24,
			"billing_currency":"CNY",
			"billing_quota":120000,
			"billing_status":"settled",
			"nested":{"status":"ready","totalCost":0.2,"apiKey":"secret"}
		}`),
	}

	result := TaskModel2LogDto(task)
	var data map[string]any
	require.NoError(t, common.Unmarshal(result.Data, &data))
	require.Equal(t, "video-model", data["model"])
	require.Equal(t, float64(10), data["duration"])
	require.NotContains(t, data, "id")
	require.NotContains(t, data, "task_id")
	require.NotContains(t, data, "billing_amount")
	require.NotContains(t, data, "billing_currency")
	require.NotContains(t, data, "billing_quota")
	require.NotContains(t, data, "billing_status")
	nested := data["nested"].(map[string]any)
	require.Equal(t, "ready", nested["status"])
	require.NotContains(t, nested, "totalCost")
	require.NotContains(t, nested, "apiKey")
}

func TestTaskModel2LogDtoDropsUnstructuredUpstreamResponse(t *testing.T) {
	task := &model.Task{Data: json.RawMessage(`upstream amount=0.24 token=secret`)}
	require.Empty(t, TaskModel2LogDto(task).Data)
}
