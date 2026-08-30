/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { describe, expect, test } from 'vitest'

import type { TaskLog } from '../../types'
import { getTaskLogDetails } from '../task-details'

const baseLog: TaskLog = {
  id: 1,
  user_id: 2,
  platform: '1',
  task_id: 'task_public',
  action: 'textGenerate',
  channel_id: 3,
  submit_time: 1,
  status: 'SUCCESS',
}

describe('getTaskLogDetails', () => {
  test('extracts a recorded video request and generated media', () => {
    const details = getTaskLogDetails({
      ...baseLog,
      result_url: 'https://example.com/result.mp4',
      properties: {
        origin_model_name: 'video-origin',
        upstream_model_name: 'video-upstream',
        input: JSON.stringify({
          prompt: 'city at night',
          model: 'video-origin',
          duration: 6,
          size: '1280x720',
          images: ['https://example.com/reference.png'],
          metadata: { aspect_ratio: '16:9', resolution: '720p' },
        }),
      },
    })

    expect(details.hasRecordedRequest).toBe(true)
    expect(details.prompt).toBe('city at night')
    expect(details.model).toBe('video-origin')
    expect(details.duration).toBe('6')
    expect(details.size).toBe('1280x720')
    expect(details.resolution).toBe('720p')
    expect(details.aspectRatio).toBe('16:9')
    expect(details.referenceImages).toEqual([
      'https://example.com/reference.png',
    ])
    expect(details.rawResponse).toBe('')
    expect(details.proxyResultUrl).toBe('/v1/videos/task_public/content')
    expect(details.resultUrl).toBe('/v1/videos/task_public/content')
    expect(details.isVideoTask).toBe(true)
  })

  test('recovers useful fields from a historical response', () => {
    const details = getTaskLogDetails({
      ...baseLog,
      properties: {
        input: '',
        origin_model_name: 'historical-video',
      },
      data: {
        detail: {
          duration: 10,
          resolution: '1080p',
          aspect_ratio: '16:9',
          url: 'https://example.com/historical.mp4',
        },
      },
    })

    expect(details.hasRecordedRequest).toBe(false)
    expect(details.model).toBe('historical-video')
    expect(details.duration).toBe('10')
    expect(details.resolution).toBe('1080p')
    expect(details.resultUrl).toBe('/v1/videos/task_public/content')
  })

  test('uses the local proxy and hides upstream URLs in the response', () => {
    const details = getTaskLogDetails({
      ...baseLog,
      result_url: '/v1/videos/task_public/content',
      properties: {
        input: JSON.stringify({
          prompt: 'cinematic city',
          aspect_ratio: '16:9',
          custom_parameter: { camera_motion: 'dolly-in' },
        }),
      },
      data: {
        status: 'succeeded',
        data: {
          model: 'upstream-video-model',
          result: {
            video_url: 'https://cdn.example.com/upstream.mp4',
            download_url: 'cdn.internal/private.mp4',
          },
          diagnostic: 'download from https://origin.example.com/private/path',
          duration: 10,
        },
      },
    })

    expect(details.prompt).toBe('cinematic city')
    expect(details.rawRequest).toContain('custom_parameter')
    expect(details.rawResponse).toContain('upstream-video-model')
    expect(details.rawResponse).not.toContain('cdn.example.com')
    expect(details.rawResponse).not.toContain('origin.example.com')
    expect(details.rawResponse).not.toContain('cdn.internal')
    expect(details.rawResponse).toContain('[hidden upstream URL]')
    expect(details.proxyResultUrl).toBe('/v1/videos/task_public/content')
    expect(details.resultUrl).toBe('/v1/videos/task_public/content')
  })
})
