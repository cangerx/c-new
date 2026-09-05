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
import { TASK_ACTIONS } from '../constants'
import type { TaskLog, TaskPluginInfo, TaskPluginRuntimeInfo } from '../types'

export interface TaskDetailAccess {
  plugin?: TaskPluginInfo
  runtime?: TaskPluginRuntimeInfo
  upstreamTaskId?: string
  nodeName?: string
}

/** Resolve elevated task metadata without exposing it to ordinary users. */
export function resolveTaskDetailAccess(
  log: TaskLog,
  isAdmin: boolean,
  isRoot: boolean
): TaskDetailAccess {
  const access: TaskDetailAccess = {}
  if (isAdmin && log.admin_info?.task_plugin) {
    access.plugin = log.admin_info.task_plugin
  }
  if (isRoot && log.root_info) {
    if (log.root_info.task_plugin) access.runtime = log.root_info.task_plugin
    if (log.root_info.upstream_task_id) {
      access.upstreamTaskId = log.root_info.upstream_task_id
    }
    if (log.root_info.node_name) access.nodeName = log.root_info.node_name
  }
  return access
}

type UnknownRecord = Record<string, unknown>

export type TaskLogDetails = {
  hasRecordedRequest: boolean
  prompt: string
  model: string
  originModel: string
  upstreamModel: string
  mode: string
  size: string
  duration: string
  resolution: string
  aspectRatio: string
  referenceImages: string[]
  metadata: UnknownRecord | null
  rawRequest: string
  rawResponse: string
  resultUrl: string
  proxyResultUrl: string
  audioUrls: string[]
  isVideoTask: boolean
}

function parseJsonValue(value: unknown): unknown {
  if (typeof value !== 'string') return value
  const trimmed = value.trim()
  if (!trimmed) return null
  try {
    return JSON.parse(trimmed)
  } catch {
    return value
  }
}

function asRecord(value: unknown): UnknownRecord | null {
  const parsed = parseJsonValue(value)
  if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
    return null
  }
  return parsed as UnknownRecord
}

function asDisplayString(value: unknown): string {
  if (typeof value === 'string') return value.trim()
  if (typeof value === 'number' && Number.isFinite(value)) return String(value)
  return ''
}

function findNestedValue(
  value: unknown,
  keys: ReadonlySet<string>,
  depth = 0
): unknown {
  if (depth > 6 || value == null) return undefined
  const parsed = parseJsonValue(value)
  if (Array.isArray(parsed)) {
    for (const item of parsed) {
      const found = findNestedValue(item, keys, depth + 1)
      if (found !== undefined) return found
    }
    return undefined
  }
  const record = asRecord(parsed)
  if (!record) return undefined

  for (const [key, child] of Object.entries(record)) {
    if (keys.has(key) && child != null && child !== '') return child
  }
  for (const child of Object.values(record)) {
    const found = findNestedValue(child, keys, depth + 1)
    if (found !== undefined) return found
  }
  return undefined
}

function collectUrls(
  value: unknown,
  keys: ReadonlySet<string>,
  result: Set<string>,
  depth = 0
) {
  if (depth > 6 || value == null) return
  const parsed = parseJsonValue(value)
  if (Array.isArray(parsed)) {
    parsed.forEach((item) => collectUrls(item, keys, result, depth + 1))
    return
  }
  const record = asRecord(parsed)
  if (!record) return

  for (const [key, child] of Object.entries(record)) {
    if (keys.has(key)) {
      const values = Array.isArray(child) ? child : [child]
      values.forEach((item) => {
        const url = asDisplayString(item)
        if (/^https?:\/\//i.test(url)) result.add(url)
      })
    }
    collectUrls(child, keys, result, depth + 1)
  }
}

function formatRawPayload(value: unknown): string {
  if (value == null || value === '') return ''
  const parsed = parseJsonValue(value)
  if (typeof parsed === 'string') return parsed
  try {
    return JSON.stringify(parsed, null, 2)
  } catch {
    return ''
  }
}

const HTTP_URL_PATTERN = /https?:\/\/[^\s"'<>\\\])}]+/gi
const URL_VALUE_KEYS = new Set([
  'url',
  'uri',
  'video_url',
  'videoUrl',
  'result_url',
  'content_url',
  'download_url',
])

const SENSITIVE_RESPONSE_KEY_PARTS = new Set([
  'amount',
  'authorization',
  'balance',
  'billing',
  'charge',
  'cost',
  'credential',
  'currency',
  'fee',
  'password',
  'payment',
  'price',
  'quota',
  'secret',
  'signature',
  'token',
])

const SENSITIVE_RESPONSE_IDS = new Set([
  'api_key',
  'apikey',
  'id',
  'private_key',
  'request_id',
  'task_id',
  'trace_id',
  'upstream_request_id',
  'upstream_task_id',
])

function isSensitiveResponseKey(key: string): boolean {
  const normalized = key
    .trim()
    .replaceAll(/([a-z0-9])([A-Z])/g, '$1_$2')
    .toLowerCase()
    .replaceAll(/[-.]/g, '_')
  if (SENSITIVE_RESPONSE_IDS.has(normalized)) return true
  return normalized
    .split('_')
    .some((part) => SENSITIVE_RESPONSE_KEY_PARTS.has(part))
}

function sanitizeResponseData(value: unknown, key = ''): unknown {
  if (key && isSensitiveResponseKey(key)) return undefined
  if (typeof value === 'string') {
    if (value && URL_VALUE_KEYS.has(key)) return '[hidden upstream URL]'
    return value.replaceAll(HTTP_URL_PATTERN, '[hidden upstream URL]')
  }
  if (Array.isArray(value)) {
    return value
      .map((item) => sanitizeResponseData(item))
      .filter((item) => item !== undefined)
  }
  const record = asRecord(value)
  if (!record) return value
  return Object.fromEntries(
    Object.entries(record).flatMap(([childKey, child]) => {
      const sanitized = sanitizeResponseData(child, childKey)
      return sanitized === undefined ? [] : [[childKey, sanitized]]
    })
  )
}

function formatSanitizedResponse(value: unknown): string {
  if (value == null || value === '') return ''
  const parsed = parseJsonValue(value)
  if (typeof parsed === 'string') return ''
  const sanitized = sanitizeResponseData(parsed)
  try {
    return JSON.stringify(sanitized, null, 2)
  } catch {
    return ''
  }
}

const VIDEO_ACTIONS = new Set<string>([
  TASK_ACTIONS.GENERATE,
  TASK_ACTIONS.TEXT_GENERATE,
  TASK_ACTIONS.FIRST_TAIL_GENERATE,
  TASK_ACTIONS.REFERENCE_GENERATE,
  TASK_ACTIONS.REMIX_GENERATE,
])

const REFERENCE_IMAGE_KEYS = new Set([
  'image',
  'images',
  'image_url',
  'input_reference',
  'reference_image',
  'reference_images',
  'first_frame_image',
  'last_frame_image',
])

export function getTaskLogDetails(log: TaskLog): TaskLogDetails {
  const properties = asRecord(log.properties) ?? {}
  const inputValue = properties.input
  const parsedInput = parseJsonValue(inputValue)
  const request = asRecord(inputValue) ?? {}
  const responseData = parseJsonValue(log.data)
  const metadata = asRecord(request.metadata)
  const hasRecordedRequest =
    inputValue != null &&
    inputValue !== '' &&
    (Object.keys(request).length > 0 || typeof inputValue === 'string')

  const originModel = asDisplayString(properties.origin_model_name)
  const upstreamModel = asDisplayString(properties.upstream_model_name)
  const responseModel = asDisplayString(
    findNestedValue(responseData, new Set(['model']))
  )
  const model = asDisplayString(request.model) || originModel || responseModel
  const duration =
    asDisplayString(request.duration) ||
    asDisplayString(request.seconds) ||
    asDisplayString(metadata?.duration) ||
    asDisplayString(metadata?.duration_seconds) ||
    asDisplayString(
      findNestedValue(responseData, new Set(['duration', 'seconds']))
    )
  const resolution =
    asDisplayString(metadata?.resolution) ||
    asDisplayString(findNestedValue(responseData, new Set(['resolution'])))
  const size = asDisplayString(request.size)
  const aspectRatio =
    asDisplayString(metadata?.aspect_ratio) ||
    asDisplayString(findNestedValue(responseData, new Set(['aspect_ratio'])))

  const referenceImages = new Set<string>()
  collectUrls(request, REFERENCE_IMAGE_KEYS, referenceImages)
  if (referenceImages.size === 0) {
    collectUrls(responseData, REFERENCE_IMAGE_KEYS, referenceImages)
  }
  const audioUrls = new Set<string>()
  collectUrls(responseData, new Set(['audio_url']), audioUrls)

  const isVideoTask = VIDEO_ACTIONS.has(log.action)
  const proxyResultUrl =
    isVideoTask && log.status === 'SUCCESS' && log.task_id
      ? `/v1/videos/${encodeURIComponent(log.task_id)}/content`
      : ''

  return {
    hasRecordedRequest,
    prompt:
      asDisplayString(request.prompt) ||
      (typeof parsedInput === 'string' ? parsedInput : '') ||
      asDisplayString(findNestedValue(responseData, new Set(['prompt']))),
    model,
    originModel,
    upstreamModel,
    mode: asDisplayString(request.mode),
    size,
    duration,
    resolution,
    aspectRatio,
    referenceImages: [...referenceImages],
    metadata,
    rawRequest: formatRawPayload(inputValue),
    rawResponse: formatSanitizedResponse(log.data),
    resultUrl: proxyResultUrl,
    proxyResultUrl,
    audioUrls: [...audioUrls],
    isVideoTask,
  }
}
