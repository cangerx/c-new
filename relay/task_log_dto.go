package relay

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"
)

var (
	taskLogHTTPURLPattern     = regexp.MustCompile(`(?i)https?://[^\s"'<>\\\])}]+`)
	taskLogKeyBoundaryPattern = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

var taskLogSensitiveResponseKeyParts = map[string]struct{}{
	"amount":        {},
	"authorization": {},
	"balance":       {},
	"billing":       {},
	"charge":        {},
	"cost":          {},
	"credential":    {},
	"currency":      {},
	"fee":           {},
	"password":      {},
	"payment":       {},
	"price":         {},
	"quota":         {},
	"secret":        {},
	"signature":     {},
	"token":         {},
}

var taskLogSensitiveResponseIDs = map[string]struct{}{
	"api_key":             {},
	"apikey":              {},
	"id":                  {},
	"private_key":         {},
	"request_id":          {},
	"task_id":             {},
	"trace_id":            {},
	"upstream_request_id": {},
	"upstream_task_id":    {},
}

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
		// Unstructured upstream payloads cannot be redacted reliably.
		return nil
	}
	sanitizedValue, _ := sanitizeTaskLogValue("", value)
	sanitized, err := common.Marshal(sanitizedValue)
	if err != nil {
		return nil
	}
	return json.RawMessage(sanitized)
}

func sanitizeTaskLogValue(key string, value any) (any, bool) {
	if isSensitiveTaskLogResponseKey(key) {
		return nil, false
	}
	switch typed := value.(type) {
	case string:
		return sanitizeTaskLogString(key, typed), true
	case []any:
		result := make([]any, 0, len(typed))
		for _, item := range typed {
			sanitized, keep := sanitizeTaskLogValue("", item)
			if keep {
				result = append(result, sanitized)
			}
		}
		return result, true
	case map[string]any:
		result := make(map[string]any, len(typed))
		for childKey, item := range typed {
			sanitized, keep := sanitizeTaskLogValue(childKey, item)
			if keep {
				result[childKey] = sanitized
			}
		}
		return result, true
	default:
		return value, true
	}
}

func isSensitiveTaskLogResponseKey(key string) bool {
	normalized := taskLogKeyBoundaryPattern.ReplaceAllString(strings.TrimSpace(key), `${1}_${2}`)
	normalized = strings.NewReplacer("-", "_", ".", "_").Replace(strings.ToLower(normalized))
	if normalized == "" {
		return false
	}
	if _, sensitive := taskLogSensitiveResponseIDs[normalized]; sensitive {
		return true
	}
	for _, part := range strings.Split(normalized, "_") {
		if _, sensitive := taskLogSensitiveResponseKeyParts[part]; sensitive {
			return true
		}
	}
	return false
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
