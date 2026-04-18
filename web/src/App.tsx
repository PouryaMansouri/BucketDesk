import {
  Check,
  ChevronLeft,
  ChevronRight,
  Cloud,
  Copy,
  Eye,
  EyeOff,
  File,
  Folder,
  Grid3X3,
  Languages,
  Loader2,
  Plus,
  RefreshCw,
  Save,
  Search,
  Server,
  Settings,
  Trash2,
  Upload,
  Wifi,
} from 'lucide-react'
import { ChangeEvent, FormEvent, ReactNode, useCallback, useEffect, useMemo, useRef, useState } from 'react'

type Locale = 'fa' | 'en'

type Profile = {
  id: string
  name: string
  endpoint: string
  region: string
  bucket: string
  accessKey: string
  secretKey?: string
  hasSecret?: boolean
  cdnUrl: string
  pathStyle: boolean
  publicBase?: string
}

type ObjectItem = {
  key: string
  name: string
  size: number
  lastModified?: string
  url: string
  mimeType: string
}

type FolderItem = {
  name: string
  prefix: string
}

type BrowseResult = {
  prefix: string
  folders: FolderItem[]
  objects: ObjectItem[]
  nextToken?: string
  limit: number
}

const messages = {
  fa: {
    appName: 'BucketDesk',
    tagline: 'مدیریت bucket بدون دسترسی به MinIO Console',
    newProfile: 'پروفایل جدید',
    createConnection: 'اتصال MinIO را بسازید',
    scopedAccess: 'دسترسی محدود S3، بدون نیاز به MinIO Console',
    upload: 'آپلود',
    connectionSettings: 'تنظیمات اتصال',
    profileName: 'نام پروفایل',
    endpoint: 'Endpoint URL',
    accessKey: 'Access Key ID',
    secretKey: 'Secret Access Key',
    savedSecretPlaceholder: 'ذخیره شده، برای تغییر پر کنید',
    region: 'Region',
    bucket: 'Bucket',
    cdnUrl: 'CDN URL',
    pathStyle: 'Use Path-Style Endpoint',
    save: 'ذخیره',
    testConnection: 'تست اتصال',
    root: 'root',
    searchObjects: 'جستجو در objectها...',
    delete: 'حذف',
    readingBucket: 'در حال خواندن bucket...',
    emptyPath: 'این مسیر خالی است یا هنوز پروفایلی انتخاب نشده.',
    page: 'صفحه',
    refresh: 'تازه‌سازی',
    grid: 'نمای شبکه‌ای',
    copyUrl: 'کپی آدرس',
    copied: 'آدرس کپی شد',
    profileSaved: 'پروفایل ذخیره شد',
    saveFailed: 'خطا در ذخیره پروفایل',
    browseFailed: 'خطا در مرور bucket',
    uploadFailed: 'خطا در آپلود',
    deleteFailed: 'خطا در حذف',
    connectionFailed: 'اتصال برقرار نشد',
    connectionOk: (bucket: string) => `اتصال موفق بود: ${bucket}`,
    uploaded: (count: string) => `${count} فایل آپلود شد`,
    confirmDelete: (count: string) => `حذف ${count} object انتخاب‌شده؟`,
    deleted: 'objectها حذف شدند',
  },
  en: {
    appName: 'BucketDesk',
    tagline: 'Manage buckets without exposing the MinIO Console',
    newProfile: 'New profile',
    createConnection: 'Create a MinIO connection',
    scopedAccess: 'Scoped S3 credentials, no MinIO Console required',
    upload: 'Upload',
    connectionSettings: 'Connection settings',
    profileName: 'Profile name',
    endpoint: 'Endpoint URL',
    accessKey: 'Access Key ID',
    secretKey: 'Secret Access Key',
    savedSecretPlaceholder: 'Saved. Fill only to replace it',
    region: 'Region',
    bucket: 'Bucket',
    cdnUrl: 'CDN URL',
    pathStyle: 'Use Path-Style Endpoint',
    save: 'Save',
    testConnection: 'Test connection',
    root: 'root',
    searchObjects: 'Search objects...',
    delete: 'Delete',
    readingBucket: 'Reading bucket...',
    emptyPath: 'This path is empty, or no profile is selected yet.',
    page: 'Page',
    refresh: 'Refresh',
    grid: 'Grid view',
    copyUrl: 'Copy URL',
    copied: 'URL copied',
    profileSaved: 'Profile saved',
    saveFailed: 'Failed to save profile',
    browseFailed: 'Failed to browse bucket',
    uploadFailed: 'Upload failed',
    deleteFailed: 'Delete failed',
    connectionFailed: 'Connection failed',
    connectionOk: (bucket: string) => `Connection succeeded: ${bucket}`,
    uploaded: (count: string) => `${count} files uploaded`,
    confirmDelete: (count: string) => `Delete ${count} selected object(s)?`,
    deleted: 'Objects deleted',
  },
}

const emptyProfile: Profile = {
  id: '',
  name: 'Production media',
  endpoint: '',
  region: 'us-east-1',
  bucket: '',
  accessKey: '',
  secretKey: '',
  cdnUrl: '',
  pathStyle: true,
}

export function App() {
  const [locale, setLocale] = useState<Locale>(() => {
    const saved = window.localStorage.getItem('bucketdesk.locale')
    return saved === 'en' || saved === 'fa' ? saved : 'fa'
  })
  const [profiles, setProfiles] = useState<Profile[]>([])
  const [activeProfile, setActiveProfile] = useState('')
  const [editing, setEditing] = useState<Profile>(emptyProfile)
  const [showSecret, setShowSecret] = useState(false)
  const [prefix, setPrefix] = useState('')
  const [query, setQuery] = useState('')
  const [view, setView] = useState<'grid' | 'list'>('grid')
  const [browse, setBrowse] = useState<BrowseResult>({ prefix: '', folders: [], objects: [], limit: 100 })
  const [tokens, setTokens] = useState<string[]>([''])
  const [page, setPage] = useState(0)
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [busy, setBusy] = useState(false)
  const [message, setMessage] = useState('')
  const fileInputRef = useRef<HTMLInputElement | null>(null)

  const t = messages[locale]
  const isRTL = locale === 'fa'

  const currentProfile = useMemo(
    () => profiles.find((profile) => profile.id === activeProfile),
    [activeProfile, profiles],
  )

  const visibleObjects = useMemo(() => {
    const clean = query.trim().toLowerCase()
    if (!clean) return browse.objects
    return browse.objects.filter((object) => object.name.toLowerCase().includes(clean) || object.key.toLowerCase().includes(clean))
  }, [browse.objects, query])

  const formatNumber = useCallback(
    (value: number | string) => (locale === 'fa' ? toPersian(value) : String(value)),
    [locale],
  )

  const loadProfiles = useCallback(async () => {
    const response = await fetch('/api/profiles')
    const data = await response.json()
    setProfiles(data.profiles || [])
    if (!activeProfile && data.profiles?.length) {
      setActiveProfile(data.profiles[0].id)
      setEditing({ ...data.profiles[0], secretKey: '' })
    }
  }, [activeProfile])

  const loadObjects = useCallback(async (nextPrefix = prefix, nextPage = page) => {
    if (!activeProfile) return
    setBusy(true)
    try {
      const params = new URLSearchParams({
        profile: activeProfile,
        prefix: nextPrefix,
        limit: '120',
      })
      const token = tokens[nextPage]
      if (token) params.set('token', token)
      const response = await fetch(`/api/objects?${params}`)
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Browse failed')
      setBrowse(data)
      if (data.nextToken) {
        setTokens((prev) => {
          const next = [...prev]
          next[nextPage + 1] = data.nextToken
          return next
        })
      }
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.browseFailed)
    } finally {
      setBusy(false)
    }
  }, [activeProfile, page, prefix, t.browseFailed, tokens])

  useEffect(() => {
    document.documentElement.lang = locale
    document.documentElement.dir = isRTL ? 'rtl' : 'ltr'
    window.localStorage.setItem('bucketdesk.locale', locale)
  }, [isRTL, locale])

  useEffect(() => {
    void loadProfiles()
  }, [loadProfiles])

  useEffect(() => {
    setSelected(new Set())
    setTokens([''])
    setPage(0)
    if (activeProfile) void loadObjects(prefix, 0)
  }, [activeProfile, prefix])

  async function saveProfile(event: FormEvent) {
    event.preventDefault()
    setBusy(true)
    try {
      const response = await fetch('/api/profiles', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editing),
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Save failed')
      setMessage(t.profileSaved)
      await loadProfiles()
      setActiveProfile(data.id)
      setEditing({ ...data, secretKey: '' })
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.saveFailed)
    } finally {
      setBusy(false)
    }
  }

  async function testConnection() {
    setBusy(true)
    try {
      const response = await fetch('/api/test', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editing),
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Test failed')
      setMessage(t.connectionOk(data.bucket))
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.connectionFailed)
    } finally {
      setBusy(false)
    }
  }

  async function uploadFiles(event: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(event.target.files || [])
    event.target.value = ''
    if (!activeProfile || files.length === 0) return

    const form = new FormData()
    for (const file of files) form.append('files', file)

    setBusy(true)
    try {
      const params = new URLSearchParams({ profile: activeProfile, prefix })
      const response = await fetch(`/api/upload?${params}`, { method: 'POST', body: form })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Upload failed')
      setMessage(t.uploaded(formatNumber(files.length)))
      await loadObjects(prefix, page)
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.uploadFailed)
    } finally {
      setBusy(false)
    }
  }

  async function deleteSelected() {
    if (!activeProfile || selected.size === 0) return
    if (!confirm(t.confirmDelete(formatNumber(selected.size)))) return
    setBusy(true)
    try {
      const params = new URLSearchParams({ profile: activeProfile })
      const response = await fetch(`/api/objects?${params}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ keys: Array.from(selected) }),
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || 'Delete failed')
      setSelected(new Set())
      setMessage(t.deleted)
      await loadObjects(prefix, page)
    } catch (error) {
      setMessage(error instanceof Error ? error.message : t.deleteFailed)
    } finally {
      setBusy(false)
    }
  }

  function switchLocale() {
    setLocale((current) => (current === 'fa' ? 'en' : 'fa'))
  }

  function openFolder(nextPrefix: string) {
    setPrefix(nextPrefix)
    setPage(0)
  }

  function goUp() {
    const parts = prefix.replace(/\/$/, '').split('/').filter(Boolean)
    parts.pop()
    setPrefix(parts.length ? `${parts.join('/')}/` : '')
    setPage(0)
  }

  function toggleKey(key: string) {
    setSelected((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key)
      else next.add(key)
      return next
    })
  }

  function nextPage() {
    if (!browse.nextToken) return
    const next = page + 1
    setPage(next)
    void loadObjects(prefix, next)
  }

  function previousPage() {
    const next = Math.max(0, page - 1)
    setPage(next)
    void loadObjects(prefix, next)
  }

  return (
    <main className="shell" dir={isRTL ? 'rtl' : 'ltr'}>
      <aside className="sidebar">
        <div className="brand">
          <span className="brandIcon"><Cloud size={24} /></span>
          <div>
            <strong>{t.appName}</strong>
            <small>{t.tagline}</small>
          </div>
        </div>

        <button className="newButton" onClick={() => setEditing(emptyProfile)}>
          <Plus size={16} />
          {t.newProfile}
        </button>

        <div className="profileList">
          {profiles.map((profile) => (
            <button
              key={profile.id}
              className={profile.id === activeProfile ? 'profile active' : 'profile'}
              onClick={() => {
                setActiveProfile(profile.id)
                setEditing({ ...profile, secretKey: '' })
              }}
            >
              <Server size={16} />
              <span>{profile.name}</span>
              <small dir="ltr">{profile.bucket}</small>
            </button>
          ))}
        </div>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <h1>{currentProfile?.name || t.createConnection}</h1>
            <p dir="ltr">{currentProfile?.publicBase || t.scopedAccess}</p>
          </div>
          <div className="actions">
            <button className="iconButton languageButton" title={locale === 'fa' ? 'English' : 'فارسی'} onClick={switchLocale}>
              <Languages size={18} />
              <span>{locale === 'fa' ? 'EN' : 'فا'}</span>
            </button>
            <button className="iconButton" title={t.refresh} onClick={() => loadObjects(prefix, page)} disabled={!activeProfile || busy}>
              {busy ? <Loader2 className="spin" size={18} /> : <RefreshCw size={18} />}
            </button>
            <button className="primaryButton" onClick={() => fileInputRef.current?.click()} disabled={!activeProfile || busy}>
              <Upload size={16} />
              {t.upload}
            </button>
            <input ref={fileInputRef} className="hidden" type="file" multiple onChange={uploadFiles} />
          </div>
        </header>

        {message && (
          <div className="notice">
            <Check size={16} />
            {message}
            <button aria-label="Dismiss" onClick={() => setMessage('')}>×</button>
          </div>
        )}

        <section className="settingsPanel">
          <div className="sectionTitle">
            <Settings size={18} />
            <h2>{t.connectionSettings}</h2>
          </div>
          <form className="settingsGrid" onSubmit={saveProfile}>
            <Field label={t.profileName}>
              <input value={editing.name} onChange={(event) => setEditing({ ...editing, name: event.target.value })} />
            </Field>
            <Field label={t.endpoint}>
              <input dir="ltr" placeholder="https://s3.example.com" value={editing.endpoint} onChange={(event) => setEditing({ ...editing, endpoint: event.target.value })} />
            </Field>
            <Field label={t.accessKey}>
              <input dir="ltr" value={editing.accessKey} onChange={(event) => setEditing({ ...editing, accessKey: event.target.value })} />
            </Field>
            <Field label={t.secretKey}>
              <div className="secretInput">
                <input dir="ltr" type={showSecret ? 'text' : 'password'} placeholder={editing.hasSecret ? t.savedSecretPlaceholder : ''} value={editing.secretKey || ''} onChange={(event) => setEditing({ ...editing, secretKey: event.target.value })} />
                <button type="button" onClick={() => setShowSecret(!showSecret)}>
                  {showSecret ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
            </Field>
            <Field label={t.region}>
              <input dir="ltr" value={editing.region} onChange={(event) => setEditing({ ...editing, region: event.target.value })} />
            </Field>
            <Field label={t.bucket}>
              <input dir="ltr" value={editing.bucket} onChange={(event) => setEditing({ ...editing, bucket: event.target.value })} />
            </Field>
            <Field label={t.cdnUrl}>
              <input dir="ltr" placeholder="https://cdn.example.com" value={editing.cdnUrl} onChange={(event) => setEditing({ ...editing, cdnUrl: event.target.value })} />
            </Field>
            <label className="toggle">
              <input type="checkbox" checked={editing.pathStyle} onChange={(event) => setEditing({ ...editing, pathStyle: event.target.checked })} />
              <span />
              {t.pathStyle}
            </label>
            <div className="formActions">
              <button type="submit" className="primaryButton" disabled={busy}>
                <Save size={16} />
                {t.save}
              </button>
              <button type="button" className="successButton" onClick={testConnection} disabled={busy}>
                <Wifi size={16} />
                {t.testConnection}
              </button>
            </div>
          </form>
        </section>

        <section className="browserPanel">
          <div className="toolbar">
            <div className="breadcrumb">
              <button onClick={() => openFolder('')}>{t.root}</button>
              {prefix && prefix.replace(/\/$/, '').split('/').map((part, index, parts) => {
                const nextPrefix = `${parts.slice(0, index + 1).join('/')}/`
                return (
                  <span key={nextPrefix}>
                    <ChevronLeft size={14} />
                    <button onClick={() => openFolder(nextPrefix)}>{part}</button>
                  </span>
                )
              })}
            </div>
            <div className="tools">
              <div className="search">
                <Search size={16} />
                <input placeholder={t.searchObjects} value={query} onChange={(event) => setQuery(event.target.value)} />
              </div>
              <button className="iconButton" title={t.grid} onClick={() => setView(view === 'grid' ? 'list' : 'grid')} data-active={view === 'grid'}>
                <Grid3X3 size={18} />
              </button>
              <button className="dangerButton" onClick={deleteSelected} disabled={selected.size === 0 || busy}>
                <Trash2 size={16} />
                {t.delete}
              </button>
            </div>
          </div>

          {busy && <div className="loading"><Loader2 className="spin" /> {t.readingBucket}</div>}

          <div className="folderGrid">
            {prefix && (
              <button className="folderCard" onClick={goUp}>
                <ChevronRight size={28} />
                <span>..</span>
              </button>
            )}
            {browse.folders.map((folder) => (
              <button className="folderCard" key={folder.prefix} onClick={() => openFolder(folder.prefix)}>
                <Folder size={32} />
                <span>{folder.name}</span>
              </button>
            ))}
          </div>

          {view === 'grid' ? (
            <div className="objectGrid">
              {visibleObjects.map((object) => (
                <ObjectCard
                  key={object.key}
                  copiedLabel={t.copied}
                  copyLabel={t.copyUrl}
                  formatNumber={formatNumber}
                  object={object}
                  selected={selected.has(object.key)}
                  onCopied={setMessage}
                  onToggle={() => toggleKey(object.key)}
                />
              ))}
            </div>
          ) : (
            <div className="objectList">
              {visibleObjects.map((object) => (
                <button key={object.key} className="objectRow" onClick={() => toggleKey(object.key)}>
                  <span className={selected.has(object.key) ? 'check selected' : 'check'}>{selected.has(object.key) && <Check size={14} />}</span>
                  <File size={18} />
                  <strong dir="ltr">{object.name}</strong>
                  <small>{formatSize(object.size, formatNumber)}</small>
                </button>
              ))}
            </div>
          )}

          {!busy && browse.folders.length === 0 && visibleObjects.length === 0 && (
            <div className="empty">{t.emptyPath}</div>
          )}

          <div className="pagination">
            <button className="iconButton" onClick={previousPage} disabled={page === 0 || busy}><ChevronRight size={18} /></button>
            <span>{t.page} {formatNumber(page + 1)}</span>
            <button className="iconButton" onClick={nextPage} disabled={!browse.nextToken || busy}><ChevronLeft size={18} /></button>
          </div>
        </section>
      </section>
    </main>
  )
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="field">
      <span>{label}</span>
      {children}
    </label>
  )
}

function ObjectCard({
  copiedLabel,
  copyLabel,
  formatNumber,
  object,
  selected,
  onCopied,
  onToggle,
}: {
  copiedLabel: string
  copyLabel: string
  formatNumber: (value: number | string) => string
  object: ObjectItem
  selected: boolean
  onCopied: (message: string) => void
  onToggle: () => void
}) {
  const isImage = object.mimeType.startsWith('image/')
  return (
    <article className={selected ? 'objectCard selected' : 'objectCard'}>
      <button className="selectButton" onClick={onToggle}>{selected && <Check size={14} />}</button>
      <div className="preview">
        {isImage ? <img src={object.url} alt="" loading="lazy" /> : <File size={40} />}
        <button
          className="copyButton"
          title={copyLabel}
          onClick={() => {
            void navigator.clipboard.writeText(object.url).then(() => onCopied(copiedLabel))
          }}
        >
          <Copy size={14} />
        </button>
      </div>
      <div className="objectMeta">
        <strong dir="ltr" title={object.name}>{object.name}</strong>
        <small>{formatSize(object.size, formatNumber)}</small>
      </div>
    </article>
  )
}

function formatSize(size: number, formatNumber: (value: number | string) => string) {
  if (size < 1024) return `${formatNumber(size)} B`
  if (size < 1024 * 1024) return `${formatNumber(Math.round(size / 1024))} KB`
  return `${formatNumber((size / 1024 / 1024).toFixed(1))} MB`
}

function toPersian(value: number | string) {
  return String(value).replace(/\d/g, (digit) => '۰۱۲۳۴۵۶۷۸۹'[Number(digit)])
}
