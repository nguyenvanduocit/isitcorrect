import { defineConfig } from 'vite'
import { crx } from '@crxjs/vite-plugin'

import manifest from './manifest.json'
import { version } from './package.json'

manifest.version = version

export default defineConfig({
  plugins: [crx({ manifest })],
  build: {
    rollupOptions: {
      output: {
        entryFileNames: `assets/[name].js`,
        chunkFileNames: `assets/[name].js`,
        assetFileNames: `assets/[name].[ext]`
      }
    }
  }
})
