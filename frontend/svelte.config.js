import nodeAdapter from '@sveltejs/adapter-node';
import vercelAdapter from '@sveltejs/adapter-vercel';

const adapter = process.env.VERCEL ? vercelAdapter() : nodeAdapter();

/** @type {import('@sveltejs/kit').Config} */
const config = {
	compilerOptions: {
		runes: ({ filename }) =>
			filename.split(/[/\\]/).includes('node_modules') ? undefined : true
	},
	kit: {
		adapter,
		paths: {
			relative: false
		}
	}
};

export default config;
