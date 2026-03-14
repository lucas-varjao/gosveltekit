import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

const frontendDir = fileURLToPath(new URL('.', import.meta.url))

function readProjectMetadata(): Record<string, string> {
    try {
        const raw = readFileSync(resolve(frontendDir, '..', 'project.env'), 'utf8')

        return raw.split('\n').reduce<Record<string, string>>((acc, line) => {
            const trimmed = line.trim()
            if (!trimmed || trimmed.startsWith('#')) {
                return acc
            }

            const separatorIndex = trimmed.indexOf('=')
            if (separatorIndex === -1) {
                return acc
            }

            const key = trimmed.slice(0, separatorIndex).trim()
            const value = trimmed
                .slice(separatorIndex + 1)
                .trim()
                .replace(/^['"]|['"]$/g, '')
            acc[key] = value
            return acc
        }, {})
    } catch {
        return {}
    }
}

function resolveAppVersion(): string {
    const envVersion = process.env.APP_VERSION?.trim()

    if (envVersion) {
        return envVersion
    }

    try {
        return readFileSync(resolve(frontendDir, '..', 'VERSION'), 'utf8').trim() || 'dev'
    } catch {
        return 'dev'
    }
}

const appVersion = resolveAppVersion()
const projectMetadata = readProjectMetadata()

export default defineConfig({
    plugins: [sveltekit()],
    define: {
        __APP_VERSION__: JSON.stringify(appVersion),
        __APP_DISPLAY_NAME__: JSON.stringify(projectMetadata.APP_DISPLAY_NAME || 'Starter App'),
        __APP_DESCRIPTION__: JSON.stringify(
            projectMetadata.APP_DESCRIPTION ||
                'A fullstack starter with SvelteKit, Go, PostgreSQL, and session-based auth.'
        )
    },
    server: {
        proxy: {
            // Proxy API requests to the backend during development
            '/auth': {
                target: 'http://localhost:8080',
                changeOrigin: true
            },
            '/api': {
                target: 'http://localhost:8080',
                changeOrigin: true
            }
        }
    }
})
