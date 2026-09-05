package plugins_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/jsplugin"
	builtinplugins "github.com/QuantumNous/new-api/plugins"
	jspluginadaptor "github.com/QuantumNous/new-api/relay/channel/task/jsplugin"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compileSoraPlugin(t *testing.T) *jsplugin.LoadedPlugin {
	t.Helper()
	source, err := builtinplugins.Source("sora")
	require.NoError(t, err)
	plugin, err := jsplugin.NewRegistry().RegisterFactory(source, jsplugin.Options{Key: "sora"})
	require.NoError(t, err)
	return plugin
}

func decodeSoraPluginMap(t *testing.T, value any) map[string]any {
	t.Helper()
	encoded, err := common.Marshal(value)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, common.Unmarshal(encoded, &decoded))
	return decoded
}

func decodeSoraPartValues(t *testing.T, value any) map[string][]any {
	t.Helper()
	submit := decodeSoraPluginMap(t, value)
	require.Equal(t, "multipart", submit["bodyType"])
	values := map[string][]any{}
	for _, raw := range submit["parts"].([]any) {
		part := raw.(map[string]any)
		name := part["name"].(string)
		if value, ok := part["value"]; ok {
			values[name] = append(values[name], value)
		}
	}
	return values
}

func decodeSoraBody(t *testing.T, value any) map[string]any {
	t.Helper()
	submit := decodeSoraPluginMap(t, value)
	assert.NotContains(t, submit, "bodyType")
	body, ok := submit["body"].(map[string]any)
	require.True(t, ok)
	return body
}

func TestSoraOpenAIVideoPreservesVendorFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	references := []any{"https://assets.example/one.png", "https://assets.example/two.png"}
	videos := []any{"https://assets.example/guide.mp4"}
	audios := []any{"https://assets.example/voice.mp3"}
	request := map[string]any{
		"model":           "public-video-alias",
		"prompt":          "animate the references",
		"seconds":         15,
		"ratio":           "9:16",
		"resolution":      "720p",
		"referenceImages": references,
		"referenceVideos": videos,
		"referenceAudios": audios,
	}

	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "alibaba/wan-3.0",
		"body":          map[string]any{"kind": "json", "value": request},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	assert.Equal(t, "image_to_video", resolved["action"])
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, float64(15), body["seconds"])
	assert.NotContains(t, body, "duration")
	assert.Equal(t, "9:16", body["ratio"])
	assert.Equal(t, "720p", body["resolution"])
	assert.Equal(t, references, body["referenceImages"])
	assert.Equal(t, audios, body["referenceAudios"])

	submitValue, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "alibaba/wan-3.0",
		"requestBody":   body,
	})
	require.NoError(t, err)
	outbound := decodeSoraBody(t, submitValue)
	assert.Equal(t, references, outbound["referenceImages"])
	assert.Equal(t, videos, outbound["referenceVideos"])
	assert.Equal(t, audios, outbound["referenceAudios"])
	assert.NotContains(t, outbound, "reference_images")
	assert.NotContains(t, outbound, "reference_videos")
	assert.NotContains(t, outbound, "reference_audios")
}

func TestSoraReferenceMediaSurviveHTTPAdaptorSerialization(t *testing.T) {
	var captured map[string]any
	var captureErr error
	upstream := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		captureErr = common.DecodeJson(request.Body, &captured)
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"id":"mock-upstream-task"}`))
	}))
	defer upstream.Close()

	plugin := compileSoraPlugin(t)
	adaptor := jspluginadaptor.New(plugin)
	info := &relaycommon.RelayInfo{
		OriginModelName: "grok-imagine-video-1.5-preview",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelBaseUrl:    upstream.URL,
			ApiKey:            "test-key",
			UpstreamModelName: "grok-imagine-video-1.5-preview",
		},
		TaskRelayInfo: &relaycommon.TaskRelayInfo{},
	}
	adaptor.Init(info)
	requestBody := map[string]any{
		"model":            "grok-imagine-video-1.5-preview",
		"prompt":           "simulate reference transmission",
		"duration":         15,
		"ratio":            "9:16",
		"resolution":       "720p",
		"reference_images": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
		"reference_videos": []any{"https://assets.example/guide.mp4"},
		"reference_audios": []any{"https://assets.example/voice.mp3"},
	}
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", bytes.NewReader(nil))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Set("task_request", requestBody)

	body, err := adaptor.BuildRequestBody(context, info)
	require.NoError(t, err)
	requestURL, err := adaptor.BuildRequestURL(info)
	require.NoError(t, err)
	outbound, err := http.NewRequest(http.MethodPost, requestURL, body)
	require.NoError(t, err)
	require.NoError(t, adaptor.BuildRequestHeader(context, outbound, info))
	response, err := upstream.Client().Do(outbound)
	require.NoError(t, err)
	response.Body.Close()
	require.NoError(t, captureErr)

	assert.Equal(t, []any{"https://assets.example/one.png", "https://assets.example/two.png"}, captured["reference_images"])
	assert.Equal(t, []any{"https://assets.example/guide.mp4"}, captured["reference_videos"])
	assert.Equal(t, []any{"https://assets.example/voice.mp3"}, captured["reference_audios"])
	assert.Equal(t, "9:16", captured["ratio"])
	assert.Equal(t, "720p", captured["resolution"])
	assert.Equal(t, float64(15), captured["duration"])
	assert.NotContains(t, captured, "images")
	assert.NotContains(t, captured, "aspect_ratio")
	assert.NotContains(t, captured, "seconds")
	assert.NotContains(t, captured, "referenceImages")
	assert.NotContains(t, captured, "referenceVideos")
	assert.NotContains(t, captured, "referenceAudios")
}

func TestSoraGrokUsesAICostCompatibleFieldNames(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://www.aicost.me",
		"apiKey":        "secret",
		"upstreamModel": "grok-imagine-video-1.5-preview",
		"requestBody": map[string]any{
			"model":            "grok-imagine-video-1.5-preview",
			"prompt":           "animate the references",
			"duration":         15,
			"ratio":            "9:16",
			"resolution":       "720p",
			"reference_images": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
			"images_base64":    []any{"data:image/png;base64,aW1hZ2U="},
			"reference_videos": []any{"https://assets.example/guide.mp4"},
			"audio_reference":  []any{"https://assets.example/voice.mp3"},
		},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, value)
	body := submit["body"].(map[string]any)
	assert.Equal(t, []any{"data:image/png;base64,aW1hZ2U=", "https://assets.example/one.png", "https://assets.example/two.png"}, body["images"])
	assert.Equal(t, "15", body["seconds"])
	assert.Equal(t, "9:16", body["aspect_ratio"])
	assert.Equal(t, "720p", body["resolution"])
	assert.Equal(t, []any{"https://assets.example/guide.mp4"}, body["videos"])
	assert.Equal(t, []any{"https://assets.example/voice.mp3"}, body["audios"])
	assert.NotContains(t, body, "reference_images")
	assert.NotContains(t, body, "reference_videos")
	assert.NotContains(t, body, "audio_reference")
	assert.NotContains(t, body, "images_base64")
	assert.NotContains(t, body, "duration")
	assert.NotContains(t, body, "ratio")
	assert.NotContains(t, submit, "bodyType")
}

func TestSoraMegabyUsesLegacyVideoFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://newapi.megabyai.cc",
		"apiKey":        "secret",
		"upstreamModel": "alibaba/wan-3.0-720p",
		"requestBody": map[string]any{
			"model":         "alibaba/wan-3.0-720p",
			"prompt":        "animate",
			"seconds":       "6",
			"aspect_ratio":  "9:16",
			"resolution":    "720p",
			"images":        []any{"https://assets.example/one.png", "https://assets.example/two.png"},
			"videos":        []any{"https://assets.example/guide.mp4"},
			"audios":        []any{"https://assets.example/voice.mp3"},
			"custom_option": "preserved",
		},
	})
	require.NoError(t, err)
	body := decodeSoraBody(t, value)
	assert.Equal(t, float64(6), body["duration"])
	assert.Equal(t, "9:16", body["ratio"])
	assert.Equal(t, []any{"https://assets.example/one.png", "https://assets.example/two.png"}, body["referenceImages"])
	assert.Equal(t, []any{"https://assets.example/guide.mp4"}, body["referenceVideos"])
	assert.Equal(t, []any{"https://assets.example/voice.mp3"}, body["referenceAudios"])
	assert.Equal(t, "preserved", body["custom_option"])
	assert.NotContains(t, body, "seconds")
	assert.NotContains(t, body, "aspect_ratio")
	assert.NotContains(t, body, "images")
	assert.NotContains(t, body, "videos")
	assert.NotContains(t, body, "audios")
}

func TestSoraMiniMaxH3BuildsMultimodalMetadata(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://newapi.megabyai.cc",
		"apiKey":        "secret",
		"upstreamModel": "minimax-h3-f",
		"requestBody": map[string]any{
			"model":            "minimax-h3-f",
			"prompt":           "use all references",
			"duration":         5,
			"ratio":            "16:9",
			"resolution":       "768P",
			"reference_images": []any{"https://assets.example/frame.png"},
			"reference_videos": []any{"https://assets.example/guide.mp4"},
			"reference_audios": []any{"https://assets.example/voice.mp3"},
			"metadata":         map[string]any{"trace": "preserved"},
		},
	})
	require.NoError(t, err)
	body := decodeSoraBody(t, value)
	assert.Equal(t, "5", body["seconds"])
	assert.NotContains(t, body, "duration")
	assert.NotContains(t, body, "reference_images")
	assert.NotContains(t, body, "reference_videos")
	assert.NotContains(t, body, "reference_audios")
	metadata := body["metadata"].(map[string]any)
	assert.Equal(t, "16:9", metadata["ratio"])
	assert.Equal(t, "768P", metadata["resolution"])
	assert.Equal(t, "preserved", metadata["trace"])
	assert.Equal(t, []any{
		map[string]any{"type": "image_url", "role": "reference_image", "image_url": map[string]any{"url": "https://assets.example/frame.png"}},
		map[string]any{"type": "video_url", "role": "reference_video", "video_url": map[string]any{"url": "https://assets.example/guide.mp4"}},
		map[string]any{"type": "audio_url", "role": "reference_audio", "audio_url": map[string]any{"url": "https://assets.example/voice.mp3"}},
	}, metadata["content"])
}

func TestSoraMiniMaxH3RejectsFrameAndReferenceMix(t *testing.T) {
	plugin := compileSoraPlugin(t)
	_, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://www.aicost.me",
		"apiKey":        "secret",
		"upstreamModel": "minimax-h3-f",
		"requestBody": map[string]any{
			"prompt":           "animate",
			"start_frame":      "https://assets.example/start.png",
			"reference_images": []any{"https://assets.example/reference.png"},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be mixed")
}

func TestSoraMiniMaxH3MetadataImageSetsImageAction(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model": "minimax-h3-f",
		"body": map[string]any{"kind": "json", "value": map[string]any{
			"model":  "minimax-h3-f",
			"prompt": "animate",
			"metadata": map[string]any{"content": []any{
				map[string]any{"type": "image_url", "role": "reference_image", "image_url": map[string]any{"url": "https://assets.example/frame.png"}},
			}},
		}},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	assert.Equal(t, "image_to_video", resolved["action"])
}

func TestSoraSubmitPreservesAllImageFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":           "animate",
			"images":           []any{"https://assets.example/one.png"},
			"referenceImages":  []any{"https://assets.example/two.png"},
			"reference_images": []any{map[string]any{"url": "https://assets.example/three.png"}},
			"image_urls": []any{
				"https://assets.example/one.png",
				map[string]any{"image_url": map[string]any{"url": "https://assets.example/four.png"}},
			},
		},
	})
	require.NoError(t, err)
	body := decodeSoraBody(t, value)
	assert.Equal(t, []any{"https://assets.example/one.png"}, body["images"])
	assert.Equal(t, []any{"https://assets.example/two.png"}, body["referenceImages"])
	assert.Equal(t, []any{map[string]any{"url": "https://assets.example/three.png"}}, body["reference_images"])
	assert.Equal(t, []any{
		"https://assets.example/one.png",
		map[string]any{"image_url": map[string]any{"url": "https://assets.example/four.png"}},
	}, body["image_urls"])
	assert.NotContains(t, body, "reference_image_urls")
}

func TestSoraSubmitPreservesVideoAndAudioFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":           "animate",
			"referenceVideos":  []any{"https://assets.example/one.mp4"},
			"reference_videos": []any{map[string]any{"video_url": "https://assets.example/two.mp4"}},
			"audio_urls":       []any{"https://assets.example/one.mp3"},
			"referenceAudios":  []any{map[string]any{"url": "https://assets.example/two.mp3"}},
		},
	})
	require.NoError(t, err)
	body := decodeSoraBody(t, value)
	assert.Equal(t, []any{"https://assets.example/one.mp4"}, body["referenceVideos"])
	assert.Equal(t, []any{map[string]any{"video_url": "https://assets.example/two.mp4"}}, body["reference_videos"])
	assert.Equal(t, []any{"https://assets.example/one.mp3"}, body["audio_urls"])
	assert.Equal(t, []any{map[string]any{"url": "https://assets.example/two.mp3"}}, body["referenceAudios"])
	assert.NotContains(t, body, "video_references")
	assert.NotContains(t, body, "audio_references")
}

func TestSoraMultipartPreservesOriginalFieldNames(t *testing.T) {
	plugin := compileSoraPlugin(t)
	file := map[string]any{
		"ref":      "request_file:input_reference",
		"field":    "input_reference",
		"filename": "frame.png",
		"mimeType": "image/png",
		"size":     12,
	}
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":          "animate",
			"referenceImages": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
			"referenceVideos": []any{"https://assets.example/guide.mp4"},
			"referenceAudios": []any{"https://assets.example/voice.mp3"},
		},
		"files": []any{file},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, value)
	assert.Equal(t, "multipart", submit["bodyType"])
	parts := submit["parts"].([]any)
	assert.Contains(t, parts, map[string]any{"name": "referenceImages", "value": "https://assets.example/one.png"})
	assert.Contains(t, parts, map[string]any{"name": "referenceImages", "value": "https://assets.example/two.png"})
	assert.Contains(t, parts, map[string]any{"name": "referenceVideos", "value": "https://assets.example/guide.mp4"})
	assert.Contains(t, parts, map[string]any{"name": "referenceAudios", "value": "https://assets.example/voice.mp3"})
	assert.Contains(t, parts, map[string]any{
		"name":     "input_reference",
		"fileRef":  "request_file:input_reference",
		"filename": "frame.png",
	})
}

func TestSoraPreservesKnownUpstreamMediaFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	request := map[string]any{
		"prompt":               "animate",
		"image":                "https://assets.example/image.png",
		"images":               []any{"https://assets.example/one.png", "https://assets.example/two.png"},
		"reference_images":     []any{"https://assets.example/reference.png"},
		"reference_image_urls": []any{"https://assets.example/reference-url.png"},
		"start_frame":          "https://assets.example/start.png",
		"end_frame":            "https://assets.example/end.png",
		"image_reference":      "https://assets.example/image-reference.png",
		"input_reference":      "https://assets.example/input-reference.png",
		"video":                "https://assets.example/video.mp4",
		"videos":               []any{"https://assets.example/one.mp4", "https://assets.example/two.mp4"},
		"reference_videos":     []any{"https://assets.example/reference.mp4"},
		"video_references":     []any{"https://assets.example/video-reference.mp4"},
		"input_video":          "https://assets.example/input.mp4",
		"audio":                "https://assets.example/audio.mp3",
		"audios":               []any{"https://assets.example/one.mp3", "https://assets.example/two.mp3"},
		"reference_audios":     []any{"https://assets.example/reference.mp3"},
		"audio_references":     []any{"https://assets.example/audio-reference.mp3"},
	}
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody":   request,
	})
	require.NoError(t, err)
	body := decodeSoraBody(t, value)
	for key, expected := range request {
		assert.Equal(t, expected, body[key], key)
	}
	assert.Equal(t, "vendor-video", body["model"])
}

func TestSoraMultipartAllowsRepeatedMediaAndArbitraryFileFields(t *testing.T) {
	plugin := compileSoraPlugin(t)
	files := []any{
		map[string]any{"ref": "request_file:images:0", "field": "images", "filename": "one.png", "mimeType": "image/png", "size": 12},
		map[string]any{"ref": "request_file:images:1", "field": "images", "filename": "two.png", "mimeType": "image/png", "size": 12},
		map[string]any{"ref": "request_file:file:0", "field": "file", "filename": "voice.mp3", "mimeType": "audio/mpeg", "size": 12},
	}
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model":         "grok-imagine-video-1.5-preview",
		"upstreamModel": "grok-imagine-video-1.5-preview",
		"body": map[string]any{
			"kind": "multipart",
			"fields": map[string]any{
				"prompt": []any{"animate"},
				"images": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
				"audios": []any{"https://assets.example/one.mp3", "https://assets.example/two.mp3"},
			},
			"files": files,
		},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	assert.Equal(t, "image_to_video", resolved["action"])
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, []any{"https://assets.example/one.png", "https://assets.example/two.png"}, body["images"])
	assert.Equal(t, []any{"https://assets.example/one.mp3", "https://assets.example/two.mp3"}, body["audios"])

	submitValue, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "grok-imagine-video-1.5-preview",
		"requestBody":   body,
		"files":         files,
	})
	require.NoError(t, err)
	parts := decodeSoraPluginMap(t, submitValue)["parts"].([]any)
	assert.Contains(t, parts, map[string]any{"name": "images", "value": "https://assets.example/one.png"})
	assert.Contains(t, parts, map[string]any{"name": "images", "value": "https://assets.example/two.png"})
	assert.Contains(t, parts, map[string]any{"name": "audios", "value": "https://assets.example/one.mp3"})
	assert.Contains(t, parts, map[string]any{"name": "audios", "value": "https://assets.example/two.mp3"})
	assert.Contains(t, parts, map[string]any{"name": "images", "fileRef": "request_file:images:0", "filename": "one.png"})
	assert.Contains(t, parts, map[string]any{"name": "images", "fileRef": "request_file:images:1", "filename": "two.png"})
	assert.Contains(t, parts, map[string]any{"name": "file", "fileRef": "request_file:file:0", "filename": "voice.mp3"})
}

func TestSoraOpenAIVideoPreservesDurationFieldNames(t *testing.T) {
	plugin := compileSoraPlugin(t)
	decode := func(t *testing.T, model string, request map[string]any) map[string]any {
		t.Helper()
		value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
			"model": model,
			"body":  map[string]any{"kind": "json", "value": request},
		})
		require.NoError(t, err)
		return decodeSoraPluginMap(t, value)["requestBody"].(map[string]any)
	}

	native := decode(t, "sora-2", map[string]any{"prompt": "native", "duration": 8})
	assert.Equal(t, float64(8), native["duration"])
	assert.NotContains(t, native, "seconds")

	vendor := decode(t, "seedance-2-5-1080p", map[string]any{"prompt": "vendor", "seconds": 10})
	assert.Equal(t, float64(10), vendor["seconds"])
	assert.NotContains(t, vendor, "duration")

	for _, model := range []string{"grok-imagine-video-1.5-preview", "minimax-h3-f", "seedance2.5-720p", "vendor-video"} {
		preserved := decode(t, model, map[string]any{"prompt": "vendor", "seconds": 6})
		assert.Equal(t, float64(6), preserved["seconds"], model)
		assert.NotContains(t, preserved, "duration", model)
	}
}

func TestSoraResponsesPreservesVideoExtensions(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_responses", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "seedance-2-5-1080p",
		"body": map[string]any{"kind": "json", "value": map[string]any{
			"model":           "public-video-alias",
			"input":           "animate",
			"seconds":         10,
			"ratio":           "16:9",
			"resolution":      "1080p",
			"referenceImages": []any{"https://assets.example/frame.png"},
			"referenceVideos": []any{"https://assets.example/guide.mp4"},
			"referenceAudios": []any{"https://assets.example/audio.mp3"},
		}},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, "image_to_video", resolved["action"])
	assert.Equal(t, float64(10), body["seconds"])
	assert.NotContains(t, body, "duration")
	assert.Equal(t, "16:9", body["ratio"])
	assert.Equal(t, "1080p", body["resolution"])
	assert.Equal(t, []any{"https://assets.example/frame.png"}, body["referenceImages"])
	assert.Equal(t, []any{"https://assets.example/guide.mp4"}, body["referenceVideos"])
	assert.Equal(t, []any{"https://assets.example/audio.mp3"}, body["referenceAudios"])
}

func TestSoraResponsesPreservesAllInputImages(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_responses", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "alibaba/wan-3.0",
		"body": map[string]any{"kind": "json", "value": map[string]any{
			"model":  "public-video-alias",
			"input":  "animate",
			"images": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
		}},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, "image_to_video", resolved["action"])
	assert.NotContains(t, body, "input_reference")
	assert.Equal(t, []any{"https://assets.example/one.png", "https://assets.example/two.png"}, body["images"])
}

func TestSoraMultipartKeepsInputReferenceFile(t *testing.T) {
	plugin := compileSoraPlugin(t)
	file := map[string]any{
		"ref":      "request_file:input_reference",
		"field":    "input_reference",
		"filename": "frame.png",
		"mimeType": "image/png",
		"size":     12,
	}
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model":         "vendor-video",
		"upstreamModel": "vendor-video",
		"body": map[string]any{
			"kind":   "multipart",
			"fields": map[string]any{"prompt": []any{"animate"}, "seconds": []any{"5"}},
			"files":  []any{file},
		},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, "5", body["seconds"])
	assert.NotContains(t, body, "duration")

	submitValue, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody":   body,
		"files":         []any{file},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, submitValue)
	assert.Equal(t, "multipart", submit["bodyType"])
	parts := submit["parts"].([]any)
	require.NotEmpty(t, parts)
	assert.Contains(t, parts, map[string]any{
		"name":     "input_reference",
		"fileRef":  "request_file:input_reference",
		"filename": "frame.png",
	})
}
