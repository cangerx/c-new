package sora

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTaskResultAcceptsNumericAndStringProgress(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantStatus   string
		wantProgress string
	}{
		{
			name:         "numeric progress",
			body:         `{"status":"in_progress","progress":42}`,
			wantStatus:   model.TaskStatusInProgress,
			wantProgress: "42%",
		},
		{
			name:         "numeric string progress",
			body:         `{"status":"in_progress","progress":"42"}`,
			wantStatus:   model.TaskStatusInProgress,
			wantProgress: "42%",
		},
		{
			name:         "percent string and succeeded status",
			body:         `{"status":"succeeded","progress":"100%"}`,
			wantStatus:   model.TaskStatusSuccess,
			wantProgress: "",
		},
	}

	adaptor := &TaskAdaptor{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adaptor.ParseTaskResult([]byte(tt.body))

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, result.Status)
			assert.Equal(t, tt.wantProgress, result.Progress)
		})
	}
}

func TestParseTaskResultRejectsInvalidProgress(t *testing.T) {
	_, err := (&TaskAdaptor{}).ParseTaskResult([]byte(`{"status":"in_progress","progress":"half"}`))

	require.ErrorContains(t, err, "invalid task progress")
}

func TestSoraBuildRequestBodyReturnsReplayablePassThroughBody(t *testing.T) {
	payload := []byte("opaque-sora-request-body")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/octet-stream")
	defer common.CleanupBodyStorage(c)

	info := &relaycommon.RelayInfo{}
	body, err := (&TaskAdaptor{}).BuildRequestBody(c, info)
	require.NoError(t, err)
	replayable, ok := body.(common.ReplayableBody)
	require.True(t, ok)

	sent, err := io.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, payload, sent)
	assert.EqualValues(t, len(payload), replayable.Size())

	replayBody, err := replayable.NewReader()
	require.NoError(t, err)
	replay, err := io.ReadAll(replayBody)
	require.NoError(t, err)
	require.NoError(t, replayBody.Close())
	assert.Equal(t, payload, replay)
}
