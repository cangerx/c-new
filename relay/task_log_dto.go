package relay

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
)

var taskLogHTTPURLPattern = regexp.MustCompile(`(?i)https?://[^\s"'<>\\\])}]+`)

// TaskModel2LogDto removes upstream locations at the dashboard API boundary.
// The original task data remains stored for polling and server-side diagnosis.
func TaskModel2LogDto(task *model.Task) *dto.TaskDto {
	result := TaskModel2Dto(task)
	result.FailReason = sanitizeTaskLogString("", result.FailReason)
	result.Data = sanitizeTaskLogResponseData(result.Data)
	if isVideoTaskAction(task.Action) && task.Status == model.TaskStatusSuccess {
		result.ResultURL = fmt.Sprintf("/v1/videos/%s/content", url.PathEscape(task.TaskID))
	} else {
		result.ResultURL = sanitizeTaskLogString("result_url", result.ResultURL)
	}
	return result
}

func isVideoTaskAction(action string) bool {
	switch constant.NormalizeTaskAction(action) {
	case constant.TaskActionImageToVideo,
		constant.TaskActionTextToVideo,
		constant.TaskActionFirstTailToVideo,
		constant.TaskActionReferenceToVideo,
		constant.TaskActionRemix:
		return true
	default:
		return false
	}
}

func sanitizeTaskLogResponseData(data json.RawMessage) json.RawMessage {
	if len(data) == 0 {
		return data
	}
	var value any
	if err := common.Unmarshal(data, &value); err != nil {
		sanitized, marshalErr := common.Marshal(sanitizeTaskLogString("", string(data)))
		if marshalErr != nil {
			return nil
		}
		return json.RawMessage(sanitized)
	}
	sanitized, err := common.Marshal(sanitizeTaskLogValue("", value))
	if err != nil {
		return nil
	}
	return json.RawMessage(sanitized)
}

func sanitizeTaskLogValue(key string, value any) any {
	switch typed := value.(type) {
	case string:
		return sanitizeTaskLogString(key, typed)
	case []any:
		result := make([]any, len(typed))
		for i, item := range typed {
			result[i] = sanitizeTaskLogValue(key, item)
		}
		return result
	case map[string]any:
		result := make(map[string]any, len(typed))
		for childKey, item := range typed {
			result[childKey] = sanitizeTaskLogValue(childKey, item)
		}
		return result
	default:
		return value
	}
}

func sanitizeTaskLogString(key string, value string) string {
	switch key {
	case "url", "uri", "video_url", "videoUrl", "result_url", "content_url", "download_url":
		if value != "" {
			return "[hidden upstream URL]"
		}
	}
	return taskLogHTTPURLPattern.ReplaceAllString(value, "[hidden upstream URL]")
}
