import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

const frontendDir = fileURLToPath(new URL('.', import.meta.url))

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

export default defineConfig({
    plugins: [sveltekit()],
    define: {
        __APP_VERSION__: JSON.stringify(appVersion)
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
