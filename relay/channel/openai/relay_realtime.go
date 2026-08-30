package openai

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/QuantumNous/new-api/service"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func OpenaiRealtimeHandler(c *gin.Context, info *relaycommon.RelayInfo) (*types.NewAPIError, *dto.RealtimeUsage) {
	if info == nil || info.ClientWs == nil || info.TargetWs == nil {
		return types.NewError(fmt.Errorf("invalid websocket connection"), types.ErrorCodeBadResponse), nil
	}

	info.IsStream = true
	clientConn := info.ClientWs
	targetConn := info.TargetWs

	clientClosed := make(chan struct{})
	targetClosed := make(chan struct{})
	sendChan := make(chan []byte, 100)
	receiveChan := make(chan []byte, 100)
	errChan := make(chan error, 2)

	localUsage := &dto.RealtimeUsage{}
	sumUsage := &dto.RealtimeUsage{}
	usageTracker := newRealtimeUsageTracker()
	var usageMu sync.Mutex

	gopool.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic in client reader: %v", r)
			}
		}()
		for {
			select {
			case <-c.Done():
				return
			default:
				_, message, err := clientConn.ReadMessage()
				if err != nil {
					if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						errChan <- fmt.Errorf("error reading from client: %v", err)
					}
					close(clientClosed)
					return
				}

				realtimeEvent := &dto.RealtimeEvent{}
				err = common.Unmarshal(message, realtimeEvent)
				if err != nil {
					errChan <- fmt.Errorf("error unmarshalling message: %v", err)
					return
				}

				if realtimeEvent.Type == dto.RealtimeEventTypeSessionUpdate {
					if realtimeEvent.Session != nil {
						if realtimeEvent.Session.Tools != nil {
							info.RealtimeTools = realtimeEvent.Session.Tools
						}
					}
				}

				textToken, audioToken, err := service.CountTokenRealtime(info, *realtimeEvent, info.UpstreamModelName)
				if err != nil {
					errChan <- fmt.Errorf("error counting text token: %v", err)
					return
				}
				logger.LogInfo(c, fmt.Sprintf("type: %s, textToken: %d, audioToken: %d", realtimeEvent.Type, textToken, audioToken))
				usageMu.Lock()
				localUsage.TotalTokens += textToken + audioToken
				localUsage.InputTokens += textToken + audioToken
				localUsage.InputTokenDetails.TextTokens += textToken
				localUsage.InputTokenDetails.AudioTokens += audioToken
				usageMu.Unlock()

				err = helper.WssString(c, targetConn, string(message))
				if err != nil {
					errChan <- fmt.Errorf("error writing to target: %v", err)
					return
				}

				select {
				case sendChan <- message:
				default:
				}
			}
		}
	})

	gopool.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic in target reader: %v", r)
			}
		}()
		for {
			select {
			case <-c.Done():
				return
			default:
				_, message, err := targetConn.ReadMessage()
				if err != nil {
					if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						errChan <- fmt.Errorf("error reading from target: %v", err)
					}
					close(targetClosed)
					return
				}
				info.SetFirstResponseTime()
				realtimeEvent := &dto.RealtimeEvent{}
				err = common.Unmarshal(message, realtimeEvent)
				if err != nil {
					errChan <- fmt.Errorf("error unmarshalling message: %v", err)
					return
				}

				if realtimeEvent.Type == dto.RealtimeEventTypeResponseDone {
					textToken, audioToken, countErr := service.CountTokenRealtime(info, *realtimeEvent, info.UpstreamModelName)
					if countErr != nil {
						errChan <- fmt.Errorf("error counting text token: %v", countErr)
						return
					}
					usageMu.Lock()
					localUsage.TotalTokens += textToken + audioToken
					localUsage.InputTokens += textToken + audioToken
					localUsage.InputTokenDetails.TextTokens += textToken
					localUsage.InputTokenDetails.AudioTokens += audioToken
					localSnapshot := *localUsage
					accepted := usageTracker.Accept(realtimeEvent)
					if accepted {
						localUsage = &dto.RealtimeUsage{}
					}
					usageMu.Unlock()

					if !accepted {
						service.RecordUsageAnomaly(c, "duplicate_realtime_response_done")
						message = rewriteRealtimeUsage(message, nil)
					} else {
						var reported *dto.RealtimeUsage
						if realtimeEvent.Response != nil {
							reported = realtimeEvent.Response.Usage
						}
						safeUsage := service.SanitizeRealtimeUsage(c, reported, &localSnapshot)
						if safeUsage.TotalTokens > 0 {
							if err := preConsumeUsage(c, info, safeUsage, sumUsage); err != nil {
								errChan <- fmt.Errorf("error consume usage: %v", err)
								return
							}
						}
						message = rewriteRealtimeUsage(message, safeUsage)
						info.IsFirstRequest = false
					}
					logger.LogInfo(c, fmt.Sprintf("realtime streaming sumUsage: %v", sumUsage))

				} else if realtimeEvent.Type == dto.RealtimeEventTypeSessionUpdated || realtimeEvent.Type == dto.RealtimeEventTypeSessionCreated {
					realtimeSession := realtimeEvent.Session
					if realtimeSession != nil {
						// update audio format
						info.InputAudioFormat = common.GetStringIfEmpty(realtimeSession.InputAudioFormat, info.InputAudioFormat)
						info.OutputAudioFormat = common.GetStringIfEmpty(realtimeSession.OutputAudioFormat, info.OutputAudioFormat)
					}
				} else {
					textToken, audioToken, err := service.CountTokenRealtime(info, *realtimeEvent, info.UpstreamModelName)
					if err != nil {
						errChan <- fmt.Errorf("error counting text token: %v", err)
						return
					}
					logger.LogInfo(c, fmt.Sprintf("type: %s, textToken: %d, audioToken: %d", realtimeEvent.Type, textToken, audioToken))
					usageMu.Lock()
					localUsage.TotalTokens += textToken + audioToken
					localUsage.OutputTokens += textToken + audioToken
					localUsage.OutputTokenDetails.TextTokens += textToken
					localUsage.OutputTokenDetails.AudioTokens += audioToken
					usageMu.Unlock()
				}

				err = helper.WssString(c, clientConn, string(message))
				if err != nil {
					errChan <- fmt.Errorf("error writing to client: %v", err)
					return
				}

				select {
				case receiveChan <- message:
				default:
				}
			}
		}
	})

	select {
	case <-clientClosed:
	case <-targetClosed:
	case err := <-errChan:
		//return service.OpenAIErrorWrapper(err, "realtime_error", http.StatusInternalServerError), nil
		logger.LogError(c, "realtime error: "+err.Error())
	case <-c.Done():
	}

	usageMu.Lock()
	remainingUsage := *localUsage
	localUsage = &dto.RealtimeUsage{}
	usageMu.Unlock()
	if remainingUsage.TotalTokens != 0 {
		_ = preConsumeUsage(c, info, &remainingUsage, sumUsage)
	}

	// check usage total tokens, if 0, use local usage

	return nil, sumUsage
}

type realtimeUsageTracker struct {
	seen          map[string]struct{}
	anonymousDone bool
}

func newRealtimeUsageTracker() *realtimeUsageTracker {
	return &realtimeUsageTracker{seen: make(map[string]struct{})}
}

func (t *realtimeUsageTracker) Accept(event *dto.RealtimeEvent) bool {
	if t == nil || event == nil || event.Type != dto.RealtimeEventTypeResponseDone {
		return false
	}
	key := ""
	if event.Response != nil {
		key = event.Response.Id
	}
	if key == "" {
		key = event.EventId
	}
	if key == "" {
		if t.anonymousDone {
			return false
		}
		t.anonymousDone = true
		return true
	}
	if _, exists := t.seen[key]; exists {
		return false
	}
	t.seen[key] = struct{}{}
	return true
}

func rewriteRealtimeUsage(message []byte, usage *dto.RealtimeUsage) []byte {
	var payload map[string]json.RawMessage
	if err := common.Unmarshal(message, &payload); err != nil {
		return message
	}
	delete(payload, "usage")
	rawResponse, ok := payload["response"]
	if !ok {
		return message
	}
	var response map[string]json.RawMessage
	if err := common.Unmarshal(rawResponse, &response); err != nil {
		return message
	}
	if usage == nil {
		delete(response, "usage")
	} else if raw, err := common.Marshal(usage); err == nil {
		response["usage"] = raw
	}
	rewrittenResponse, err := common.Marshal(response)
	if err != nil {
		return message
	}
	payload["response"] = rewrittenResponse
	rewritten, err := common.Marshal(payload)
	if err != nil {
		return message
	}
	return rewritten
}

func preConsumeUsage(ctx *gin.Context, info *relaycommon.RelayInfo, usage *dto.RealtimeUsage, totalUsage *dto.RealtimeUsage) error {
	if usage == nil || totalUsage == nil {
		return fmt.Errorf("invalid usage pointer")
	}

	totalUsage.TotalTokens = addRealtimeTokenCount(totalUsage.TotalTokens, usage.TotalTokens)
	totalUsage.InputTokens = addRealtimeTokenCount(totalUsage.InputTokens, usage.InputTokens)
	totalUsage.OutputTokens = addRealtimeTokenCount(totalUsage.OutputTokens, usage.OutputTokens)
	totalUsage.InputTokenDetails.CachedTokens = addRealtimeTokenCount(totalUsage.InputTokenDetails.CachedTokens, usage.InputTokenDetails.CachedTokens)
	totalUsage.InputTokenDetails.TextTokens = addRealtimeTokenCount(totalUsage.InputTokenDetails.TextTokens, usage.InputTokenDetails.TextTokens)
	totalUsage.InputTokenDetails.AudioTokens = addRealtimeTokenCount(totalUsage.InputTokenDetails.AudioTokens, usage.InputTokenDetails.AudioTokens)
	totalUsage.OutputTokenDetails.TextTokens = addRealtimeTokenCount(totalUsage.OutputTokenDetails.TextTokens, usage.OutputTokenDetails.TextTokens)
	totalUsage.OutputTokenDetails.AudioTokens = addRealtimeTokenCount(totalUsage.OutputTokenDetails.AudioTokens, usage.OutputTokenDetails.AudioTokens)
	// clear usage
	err := service.PreWssConsumeQuota(ctx, info, usage)
	return err
}

func addRealtimeTokenCount(current, delta int) int {
	if current < 0 || delta < 0 || current > math.MaxInt-delta {
		return math.MaxInt
	}
	return current + delta
}
