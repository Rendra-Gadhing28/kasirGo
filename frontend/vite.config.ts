import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit()
	],
	server: {
		port: 5173,
		proxy: {
			'/api': {
				target: process.env.BACKEND_URL || 'https://kasirgo-production-ceb2.up.railway.app',
				changeOrigin: true
			}
		}
	}
});
