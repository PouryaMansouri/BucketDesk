import { cp, rm } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')

await rm(resolve(root, 'internal/server/web/dist'), { recursive: true, force: true })
await cp(resolve(root, 'web/dist'), resolve(root, 'internal/server/web/dist'), { recursive: true })
