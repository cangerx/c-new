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
import {
  Check,
  Clock3,
  Copy,
  FileJson2,
  Film,
  Image as ImageIcon,
  Info,
  Loader2,
  ReceiptText,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Dialog } from '@/components/dialog'
import { StatusBadge } from '@/components/status-badge'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { useCopyToClipboard } from '@/hooks/use-copy-to-clipboard'
import { api } from '@/lib/api'
import { formatLogQuota, formatTimestampToDate } from '@/lib/format'

import { taskActionMapper, taskStatusMapper } from '../../lib/mappers'
import { getTaskLogDetails } from '../../lib/task-details'
import type { TaskLog } from '../../types'

type TaskDetailsDialogProps = {
  log: TaskLog
  isAdmin: boolean
  isRoot?: boolean
  open: boolean
  onOpenChange: (open: boolean) => void
}

function DetailRow({
  label,
  value,
}: {
  label: string
  value: React.ReactNode
}) {
  if (value == null || value === '') return null
  return (
    <div className='grid min-w-0 grid-cols-[6rem_minmax(0,1fr)] gap-3 text-xs'>
      <span className='text-muted-foreground'>{label}</span>
      <span className='min-w-0 break-all'>{value}</span>
    </div>
  )
}

function DetailSection({
  icon,
  title,
  children,
}: {
  icon: React.ReactNode
  title: string
  children: React.ReactNode
}) {
  return (
    <section className='space-y-2.5'>
      <Label className='flex items-center gap-2 text-xs font-semibold'>
        {icon}
        {title}
      </Label>
      <div className='bg-muted/25 space-y-2.5 rounded-md border p-3'>
        {children}
      </div>
    </section>
  )
}

function AuthenticatedVideo({ src }: { src: string }) {
  const { t } = useTranslation()
  const [objectUrl, setObjectUrl] = useState('')
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    if (!src) {
      setObjectUrl('')
      setFailed(false)
      return
    }

    const controller = new AbortController()
    let createdUrl = ''
    setObjectUrl('')
    setFailed(false)

    void api
      .get<Blob>(src, {
        responseType: 'blob',
        signal: controller.signal,
        disableDuplicate: true,
        skipErrorHandler: true,
      })
      .then((response) => {
        if (controller.signal.aborted) return
        createdUrl = URL.createObjectURL(response.data)
        setObjectUrl(createdUrl)
      })
      .catch(() => {
        if (!controller.signal.aborted) setFailed(true)
      })

    return () => {
      controller.abort()
      if (createdUrl) URL.revokeObjectURL(createdUrl)
    }
  }, [src])

  if (failed) {
    return (
      <div className='bg-background text-muted-foreground flex aspect-video w-full items-center justify-center rounded-md border text-xs'>
        {t('Failed to load')}
      </div>
    )
  }
  if (!objectUrl) {
    return (
      <div className='bg-background flex aspect-video w-full items-center justify-center rounded-md border'>
        <Loader2 className='text-muted-foreground size-5 animate-spin' />
      </div>
    )
  }
  return (
    <video
      src={objectUrl}
      controls
      preload='metadata'
      className='bg-background aspect-video w-full rounded-md border object-contain'
    />
  )
}

export function TaskDetailsDialog({
  log,
  isAdmin,
  open,
  onOpenChange,
}: TaskDetailsDialogProps) {
  const { t } = useTranslation()
  const { copiedText, copyToClipboard } = useCopyToClipboard({ notify: false })
  const details = getTaskLogDetails(log)
  const statusLabel = t(
    taskStatusMapper.getLabel(log.status, log.status || 'Submitting')
  )
  const actionLabel = t(taskActionMapper.getLabel(log.action, log.action))
  const startedAt = log.start_time || log.submit_time
  const finishedAt = log.finish_time || 0
  const elapsedSeconds =
    finishedAt > startedAt ? Math.max(0, finishedAt - startedAt) : 0
  const metadataText = details.metadata
    ? JSON.stringify(details.metadata, null, 2)
    : ''

  return (
    <Dialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('Task Details')}
      description={t(
        'Inspect the request, generated media, and execution information.'
      )}
      contentClassName='sm:max-w-4xl'
      contentHeight='min(75vh, 760px)'
      bodyClassName='space-y-5 py-1'
    >
      <div className='flex min-w-0 flex-wrap items-center gap-2'>
        <StatusBadge
          label={statusLabel}
          variant={taskStatusMapper.getVariant(log.status)}
          size='sm'
          copyable={false}
        />
        <StatusBadge
          label={log.task_id}
          copyText={log.task_id}
          variant='neutral'
          size='sm'
          className='max-w-full font-mono'
        />
      </div>

      {(details.resultUrl || details.audioUrls.length > 0) && (
        <DetailSection
          icon={<Film className='size-4' />}
          title={t('Generated Media')}
        >
          {details.isVideoTask && details.resultUrl && (
            <AuthenticatedVideo src={open ? details.resultUrl : ''} />
          )}
          {details.audioUrls.map((url) => (
            <audio
              key={url}
              src={url}
              controls
              preload='metadata'
              className='w-full'
            />
          ))}
        </DetailSection>
      )}

      <DetailSection
        icon={<Info className='size-4' />}
        title={t('Request Content')}
      >
        {!details.hasRecordedRequest && (
          <p className='text-muted-foreground text-xs'>
            {t('This historical task did not record request details.')}
          </p>
        )}
        {details.prompt && (
          <div className='relative'>
            <Button
              type='button'
              variant='ghost'
              size='icon-sm'
              className='absolute top-0 right-0'
              onClick={() => copyToClipboard(details.prompt)}
              title={t('Copy to clipboard')}
            >
              {copiedText === details.prompt ? (
                <Check className='size-4 text-emerald-600' />
              ) : (
                <Copy className='size-4' />
              )}
            </Button>
            <p className='pr-9 text-sm leading-relaxed break-words whitespace-pre-wrap'>
              {details.prompt}
            </p>
          </div>
        )}
        <div className='grid gap-2 sm:grid-cols-2 sm:gap-x-6'>
          <DetailRow label={t('Model')} value={details.model} />
          <DetailRow label={t('Origin Model')} value={details.originModel} />
          <DetailRow
            label={t('Upstream Model')}
            value={details.upstreamModel}
          />
          <DetailRow label={t('Mode')} value={details.mode} />
          <DetailRow label={t('Size')} value={details.size} />
          <DetailRow label={t('Duration')} value={details.duration} />
          <DetailRow label={t('Resolution')} value={details.resolution} />
          <DetailRow label={t('Aspect Ratio')} value={details.aspectRatio} />
        </div>
      </DetailSection>

      {details.referenceImages.length > 0 && (
        <DetailSection
          icon={<ImageIcon className='size-4' />}
          title={t('Reference Images')}
        >
          <div className='grid grid-cols-2 gap-2 sm:grid-cols-3'>
            {details.referenceImages.map((url) => (
              <a
                key={url}
                href={url}
                target='_blank'
                rel='noopener noreferrer'
                className='bg-background aspect-video overflow-hidden rounded-md border'
              >
                <img
                  src={url}
                  alt={t('Reference image')}
                  loading='lazy'
                  className='size-full object-cover'
                />
              </a>
            ))}
          </div>
        </DetailSection>
      )}

      {(details.rawRequest || metadataText) && (
        <DetailSection
          icon={<FileJson2 className='size-4' />}
          title={t('Raw Request')}
        >
          <div className='relative'>
            <Button
              type='button'
              variant='ghost'
              size='icon-sm'
              className='absolute top-0 right-0'
              onClick={() =>
                copyToClipboard(details.rawRequest || metadataText)
              }
              title={t('Copy to clipboard')}
            >
              {copiedText === (details.rawRequest || metadataText) ? (
                <Check className='size-4 text-emerald-600' />
              ) : (
                <Copy className='size-4' />
              )}
            </Button>
            <pre className='max-h-72 overflow-auto pr-9 font-mono text-xs leading-relaxed whitespace-pre-wrap'>
              {details.rawRequest || metadataText}
            </pre>
          </div>
        </DetailSection>
      )}

      {details.rawResponse && (
        <DetailSection
          icon={<FileJson2 className='size-4' />}
          title={t('Upstream Response')}
        >
          <div className='relative'>
            <Button
              type='button'
              variant='ghost'
              size='icon-sm'
              className='absolute top-0 right-0'
              onClick={() => copyToClipboard(details.rawResponse)}
              title={t('Copy to clipboard')}
            >
              {copiedText === details.rawResponse ? (
                <Check className='size-4 text-emerald-600' />
              ) : (
                <Copy className='size-4' />
              )}
            </Button>
            <pre className='max-h-96 overflow-auto pr-9 font-mono text-xs leading-relaxed whitespace-pre-wrap'>
              {details.rawResponse}
            </pre>
          </div>
        </DetailSection>
      )}

      <DetailSection
        icon={<Clock3 className='size-4' />}
        title={t('Execution Details')}
      >
        <div className='grid gap-2 sm:grid-cols-2 sm:gap-x-6'>
          <DetailRow label={t('Platform')} value={t(log.platform)} />
          <DetailRow label={t('Action')} value={actionLabel} />
          <DetailRow label={t('Progress')} value={log.progress || '-'} />
          <DetailRow
            label={t('Submit Time')}
            value={formatTimestampToDate(log.submit_time, 'seconds')}
          />
          <DetailRow
            label={t('Finish Time')}
            value={
              log.finish_time
                ? formatTimestampToDate(log.finish_time, 'seconds')
                : '-'
            }
          />
          <DetailRow
            label={t('Duration')}
            value={elapsedSeconds ? `${elapsedSeconds}s` : '-'}
          />
          {isAdmin && (
            <DetailRow label={t('Channel')} value={String(log.channel_id)} />
          )}
          {isAdmin && (
            <DetailRow
              label={t('User')}
              value={log.username || String(log.user_id)}
            />
          )}
        </div>
      </DetailSection>

      {(log.quota != null || log.group) && (
        <DetailSection
          icon={<ReceiptText className='size-4' />}
          title={t('Billing Details')}
        >
          <DetailRow
            label={t('Fee')}
            value={log.quota != null ? formatLogQuota(log.quota) : '-'}
          />
          <DetailRow label={t('Group')} value={log.group || '-'} />
        </DetailSection>
      )}

      {log.fail_reason && !/^https?:\/\//i.test(log.fail_reason) && (
        <div className='space-y-2'>
          <Label className='text-xs font-semibold text-red-600 dark:text-red-400'>
            {t('Fail Reason')}
          </Label>
          <pre className='border-destructive/30 bg-destructive/5 max-h-56 overflow-auto rounded-md border p-3 text-xs break-words whitespace-pre-wrap text-red-700 dark:text-red-300'>
            {log.fail_reason}
          </pre>
        </div>
      )}
    </Dialog>
  )
}
