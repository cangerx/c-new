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
import { AlertTriangle, Check, Copy, ExternalLink } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Dialog } from '@/components/dialog'
import { StatusBadge } from '@/components/status-badge'
import { Button } from '@/components/ui/button'
import { useCopyToClipboard } from '@/hooks/use-copy-to-clipboard'
import { formatLogQuota, formatTimestampToDate } from '@/lib/format'

import { DetailRow, DetailSection } from './detail-section'
import {
  taskActionMapper,
  taskPlatformMapper,
  taskStatusMapper,
} from '../../lib/mappers'
import type { TaskLog } from '../../types'

interface TaskLogDetailsDialogProps {
  log: TaskLog
  open: boolean
  onOpenChange: (open: boolean) => void
}

/**
 * Full details for a task log row. Unlike UsageLog's DetailsDialog, task logs
 * carry their billing group, quota, prompt and result URL directly on the
 * record, so all of it is surfaced here without extra requests.
 */
export function TaskLogDetailsDialog({
  log,
  open,
  onOpenChange,
}: TaskLogDetailsDialogProps) {
  const { t } = useTranslation()
  const { copiedText, copyToClipboard } = useCopyToClipboard({ notify: false })

  const platform = t(taskPlatformMapper.getLabel(log.platform))
  const action = t(taskActionMapper.getLabel(log.action))
  const statusLabel = t(
    taskStatusMapper.getLabel(log.status, log.status || 'Submitting')
  )

  const queueSeconds =
    log.start_time && log.start_time > 0
      ? Math.max(0, log.start_time - log.submit_time)
      : undefined
  const runSeconds =
    log.start_time && log.finish_time && log.start_time > 0
      ? Math.max(0, log.finish_time - log.start_time)
      : undefined

  const dataJson = log.data ? parseTaskDataRaw(log.data) : undefined
  // Upstream-returned video URL (e.g. R2 signed URL) lives in task.data;
  // result_url is the local proxy endpoint, so prefer the real upstream URL.
  const videoUrl = log.data ? parseTaskDataVideoUrl(log.data) : undefined

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('Task Details')}
      description={t('View the complete details for this task')}
      contentClassName='max-sm:max-h-[calc(100dvh-1.5rem)] max-sm:w-[calc(100vw-1.5rem)] max-sm:max-w-[calc(100vw-1.5rem)] max-sm:p-4 sm:max-w-lg'
      contentHeight='min(72dvh, 720px)'
      bodyClassName='pr-2 sm:pr-4'
    >
      <div className='w-full max-w-full min-w-0 space-y-2.5 overflow-x-hidden py-1 sm:space-y-3'>
        {/* Header: status badge + platform/action */}
        <div className='flex flex-wrap items-center gap-2'>
          <StatusBadge
            label={statusLabel}
            variant={taskStatusMapper.getVariant(log.status)}
            size='sm'
            copyable={false}
          />
          <span className='text-muted-foreground text-xs'>
            {platform} · {action}
          </span>
          {log.task_id && (
            <StatusBadge
              label={log.task_id}
              copyText={log.task_id}
              variant='neutral'
              size='sm'
              className='ml-auto max-w-full border-border/60 bg-muted/30 !text-foreground truncate rounded-md border px-1.5 py-0.5 font-mono'
            />
          )}
        </div>

        {/* Identity */}
        <DetailSection label={t('Identity')}>
          <DetailRow label={t('Task ID')} value={log.task_id || '-'} mono />
          <DetailRow label={t('Platform')} value={platform} />
          <DetailRow label={t('Action')} value={action} />
          <DetailRow label={t('Status')} value={statusLabel} />
          {log.group && <DetailRow label={t('Group')} value={log.group} mono />}
          {log.username && (
            <DetailRow label={t('User')} value={log.username} mono />
          )}
        </DetailSection>

        {/* Billing */}
        {log.quota != null && (
          <DetailSection label={t('Billing Details')}>
            <DetailRow
              label={t('Quota')}
              value={formatLogQuota(log.quota)}
              mono
            />
            {log.channel_id > 0 && (
              <DetailRow
                label={t('Channel')}
                value={String(log.channel_id)}
                mono
              />
            )}
          </DetailSection>
        )}

        {/* Timing */}
        <DetailSection label={t('Timing')}>
          <DetailRow
            label={t('Submit Time')}
            value={formatTimestampToDate(log.submit_time, 'seconds')}
            mono
          />
          {log.start_time && log.start_time > 0 && (
            <>
              <DetailRow
                label={t('Start Time')}
                value={formatTimestampToDate(log.start_time, 'seconds')}
                mono
              />
              {queueSeconds != null && (
                <DetailRow
                  label={t('Queue Wait')}
                  value={`${queueSeconds}s`}
                  mono
                />
              )}
            </>
          )}
          {log.finish_time && log.finish_time > 0 && (
            <>
              <DetailRow
                label={t('Finish Time')}
                value={formatTimestampToDate(log.finish_time, 'seconds')}
                mono
              />
              {runSeconds != null && (
                <DetailRow label={t('Run Time')} value={`${runSeconds}s`} mono />
              )}
            </>
          )}
          {log.progress && (
            <DetailRow label={t('Progress')} value={log.progress} mono />
          )}
        </DetailSection>

        {/* Prompt / model */}
        {log.properties?.input && (
          <DetailSection label={t('Prompt')}>
            <DetailRow
              label={t('Prompt')}
              value={log.properties.input}
              muted={false}
            />
          </DetailSection>
        )}
        {(log.properties?.origin_model_name ||
          log.properties?.upstream_model_name) && (
          <DetailSection label={t('Model')}>
            {log.properties?.origin_model_name && (
              <DetailRow
                label={t('Model')}
                value={log.properties.origin_model_name}
                mono
              />
            )}
            {log.properties?.upstream_model_name &&
              log.properties.upstream_model_name !==
                log.properties.origin_model_name && (
                <DetailRow
                  label={t('Upstream Model')}
                  value={log.properties.upstream_model_name}
                  mono
                />
              )}
          </DetailSection>
        )}

        {/* Result: upstream video URL (from task data) with inline player */}
        {(videoUrl || log.result_url) && (
          <DetailSection label={t('Result')}>
            {videoUrl ? (
              <>
                <video
                  controls
                  preload='metadata'
                  src={videoUrl}
                  className='w-full rounded-md border bg-black'
                />
                <a
                  href={videoUrl}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='flex items-center gap-1 text-xs break-all text-blue-600 hover:underline dark:text-blue-400'
                >
                  {videoUrl}
                  <ExternalLink className='size-3 shrink-0' />
                </a>
              </>
            ) : (
              <a
                href={log.result_url}
                target='_blank'
                rel='noopener noreferrer'
                className='flex items-center gap-1 text-xs text-blue-600 hover:underline dark:text-blue-400'
              >
                {log.result_url}
                <ExternalLink className='size-3 shrink-0' />
              </a>
            )}
          </DetailSection>
        )}

        {/* Fail reason */}
        {log.fail_reason && (
          <DetailSection
            label={t('Fail Reason')}
            variant='danger'
            icon={
              <AlertTriangle className='size-3.5' aria-hidden='true' />
            }
          >
            <div className='relative'>
              <Button
                variant='ghost'
                size='sm'
                className='absolute top-0 right-0 h-7 w-7 p-0'
                onClick={() => copyToClipboard(log.fail_reason ?? '')}
                title={t('Copy to clipboard')}
              >
                {copiedText === log.fail_reason ? (
                  <Check className='size-3.5 text-green-600' />
                ) : (
                  <Copy className='size-3.5' />
                )}
              </Button>
              <p className='overflow-wrap-anywhere pr-8 text-xs leading-relaxed break-all whitespace-pre-wrap text-red-600'>
                {log.fail_reason}
              </p>
            </div>
          </DetailSection>
        )}

        {/* Raw data */}
        {dataJson && (
          <DetailSection label={t('Raw Data')}>
            <pre className='text-muted-foreground max-h-48 overflow-auto font-mono text-[11px] leading-relaxed whitespace-pre-wrap'>
              {dataJson}
            </pre>
          </DetailSection>
        )}
      </div>
    </Dialog>
  )
}

function parseTaskDataVideoUrl(data: unknown): string | undefined {
  if (data == null) return undefined
  try {
    // dto.TaskDto forwards json.RawMessage: data may be a raw JSON string or
    // an already-parsed object depending on transport.
    const parsed = typeof data === 'string' ? JSON.parse(data) : data
    // Video responses are objects; suno audio responses are arrays.
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return undefined
    }
    const obj = parsed as Record<string, unknown>
    const metadata = obj.metadata as Record<string, unknown> | undefined
    const candidates = [
      obj.url,
      obj.video_url,
      metadata?.url,
      metadata?.origin_video_url,
    ]
    for (const candidate of candidates) {
      if (typeof candidate === 'string' && candidate.trim()) {
        return candidate
      }
    }
    return undefined
  } catch {
    return undefined
  }
}

function parseTaskDataRaw(data: unknown): string | undefined {
  if (data == null) return undefined
  try {
    // dto.TaskDto forwards json.RawMessage, so data may arrive as a raw
    // string or as an already-parsed array/object depending on transport.
    const parsed = typeof data === 'string' ? JSON.parse(data) : data
    return JSON.stringify(parsed, null, 2)
  } catch {
    return undefined
  }
}
