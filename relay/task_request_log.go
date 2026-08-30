package relay

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
)

const (
	maxTaskRequestLogBytes    = 32 * 1024
	maxTaskRequestSourceBytes = 1024 * 1024
	maxTaskRequestLogDepth    = 6
	maxTaskRequestLogItems    = 50
	maxTaskRequestTextRunes   = 16 * 1024
	maxTaskRequestValueRunes  = 2 * 1024
)

// BuildTaskRequestLogInput returns a bounded, redacted JSON snapshot of the
// parsed task request. The snapshot is user-visible and must never include
// credentials or large embedded image payloads.
func BuildTaskRequestLogInput(c *gin.Context) string {
	if snapshot := buildRawTaskRequestSnapshot(c); snapshot != nil {
		return encodeTaskRequestLogSnapshot(snapshot)
	}

	request, exists := c.Get("task_request")
	if !exists || request == nil {
		return ""
	}

	var snapshot any
	if taskRequest, ok := request.(relaycommon.TaskSubmitReq); ok {
		snapshot = buildVideoTaskRequestSnapshot(taskRequest)
	} else {
		raw, err := common.Marshal(request)
		if err != nil {
			return ""
		}
		if err := common.Unmarshal(raw, &snapshot); err != nil {
			return ""
		}
		snapshot = sanitizeTaskRequestLogValue("", snapshot, 0)
	}
	return encodeTaskRequestLogSnapshot(snapshot)
}

func buildRawTaskRequestSnapshot(c *gin.Context) any {
	if c == nil || c.Request == nil {
		return nil
	}
	if !strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "json") {
		return nil
	}
	storage, err := common.GetBodyStorage(c)
	if err != nil || storage.Size() == 0 || storage.Size() > maxTaskRequestSourceBytes {
		return nil
	}
	raw, err := storage.Bytes()
	if err != nil {
		return nil
	}

	var snapshot any
	if err := common.Unmarshal(raw, &snapshot); err != nil || snapshot == nil {
		return nil
	}
	return sanitizeTaskRequestLogValue("", snapshot, 0)
}

func encodeTaskRequestLogSnapshot(snapshot any) string {
	encoded, err := common.Marshal(snapshot)
	if err != nil {
		return ""
	}
	if len(encoded) <= maxTaskRequestLogBytes {
		return string(encoded)
	}

	if requestSnapshot, ok := snapshot.(map[string]any); ok {
		compact := map[string]any{"_truncated": true}
		for _, key := range []string{
			"prompt", "model", "mode", "size", "duration", "seconds",
			"resolution", "aspect_ratio", "image", "images", "input_reference",
		} {
			if value, exists := requestSnapshot[key]; exists {
				compact[key] = value
			}
		}
		encoded, err = common.Marshal(compact)
		if err == nil && len(encoded) <= maxTaskRequestLogBytes {
			return string(encoded)
		}
	}

	fallback, _ := common.Marshal(map[string]any{
		"_truncated": true,
		"summary":    "request exceeded task log size limit",
	})
	return string(fallback)
}

func buildVideoTaskRequestSnapshot(request relaycommon.TaskSubmitReq) map[string]any {
	snapshot := make(map[string]any)
	addTaskRequestLogString(snapshot, "prompt", request.Prompt)
	addTaskRequestLogString(snapshot, "model", request.Model)
	addTaskRequestLogString(snapshot, "mode", request.Mode)
	addTaskRequestLogString(snapshot, "size", request.Size)
	addTaskRequestLogString(snapshot, "seconds", request.Seconds)
	if request.Duration != 0 {
		snapshot["duration"] = request.Duration
	}
	addTaskRequestLogString(snapshot, "image", request.Image)
	addTaskRequestLogString(snapshot, "input_reference", request.InputReference)
	if len(request.Images) > 0 {
		snapshot["images"] = sanitizeTaskRequestLogValue("images", request.Images, 0)
	}
	if len(request.Metadata) > 0 {
		snapshot["metadata"] = sanitizeTaskRequestLogValue("metadata", request.Metadata, 0)
	}
	return snapshot
}

func addTaskRequestLogString(snapshot map[string]any, key string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	snapshot[key] = sanitizeTaskRequestLogString(key, value)
}

func sanitizeTaskRequestLogValue(key string, value any, depth int) any {
	if depth >= maxTaskRequestLogDepth {
		return "[nested value omitted]"
	}
	if isSensitiveTaskRequestLogKey(key) {
		return "[redacted]"
	}

	switch typed := value.(type) {
	case string:
		return sanitizeTaskRequestLogString(key, typed)
	case []string:
		limit := min(len(typed), maxTaskRequestLogItems)
		result := make([]any, 0, limit+1)
		for _, item := range typed[:limit] {
			result = append(result, sanitizeTaskRequestLogString(key, item))
		}
		if len(typed) > limit {
			result = append(result, fmt.Sprintf("[%d more items omitted]", len(typed)-limit))
		}
		return result
	case []any:
		limit := min(len(typed), maxTaskRequestLogItems)
		result := make([]any, 0, limit+1)
		for _, item := range typed[:limit] {
			result = append(result, sanitizeTaskRequestLogValue(key, item, depth+1))
		}
		if len(typed) > limit {
			result = append(result, fmt.Sprintf("[%d more items omitted]", len(typed)-limit))
		}
		return result
	case map[string]interface{}:
		result := make(map[string]any)
		count := 0
		for childKey, childValue := range typed {
			if count >= maxTaskRequestLogItems {
				result["_truncated"] = true
				break
			}
			result[childKey] = sanitizeTaskRequestLogValue(childKey, childValue, depth+1)
			count++
		}
		return result
	default:
		return value
	}
}

func sanitizeTaskRequestLogString(key string, value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(trimmed), "data:") {
		return fmt.Sprintf("[embedded data omitted: %d bytes]", len(value))
	}

	limit := maxTaskRequestValueRunes
	if strings.EqualFold(key, "prompt") {
		limit = maxTaskRequestTextRunes
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit]) + fmt.Sprintf("\n[%d characters omitted]", len(runes)-limit)
}

func isSensitiveTaskRequestLogKey(key string) bool {
	normalized := strings.NewReplacer("-", "", "_", "", " ", "").Replace(
		strings.ToLower(strings.TrimSpace(key)),
	)
	switch normalized {
	case "apikey", "authorization", "auth", "credential", "credentials", "key", "password", "passwd", "secret", "token", "accesstoken", "refreshtoken":
		return true
	default:
		return strings.Contains(normalized, "token") ||
			strings.Contains(normalized, "secret") ||
			strings.Contains(normalized, "password") ||
			strings.Contains(normalized, "authorization")
	}
}
