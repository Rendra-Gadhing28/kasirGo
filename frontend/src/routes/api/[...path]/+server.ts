import type { RequestHandler } from './$types';

const DEFAULT_HOST_URL = 'https://kasir-go-be.vercel.app';
const BACKEND_INTERNAL_URL = process.env.BACKEND_URL || (process.env.VERCEL ? DEFAULT_HOST_URL : 'http://backend:8080');
const BACKEND_LOCAL_URL = process.env.BACKEND_URL || DEFAULT_HOST_URL;

async function proxyRequest(event: Parameters<RequestHandler>[0]) {
  const { request, params, url } = event;
  
  // Decide backend URL: use BACKEND_URL if set, otherwise fallback depending on env
  const baseBackend = process.env.BACKEND_URL || (process.env.NODE_ENV === 'production' ? BACKEND_INTERNAL_URL : BACKEND_LOCAL_URL);
  const cleanBase = baseBackend.replace(/\/+$/, '');
  const targetUrl = `${cleanBase}/api/${params.path}${url.search}`;

  const headers = new Headers(request.headers);
  // Remove host header to let fetch set target host
  headers.delete('host');

  let body: ArrayBuffer | undefined = undefined;
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    body = await request.arrayBuffer();
  }

  try {
    const res = await fetch(targetUrl, {
      method: request.method,
      headers,
      body,
      redirect: 'manual'
    });

    const responseHeaders = new Headers(res.headers);
    return new Response(res.body, {
      status: res.status,
      statusText: res.statusText,
      headers: responseHeaders
    });
  } catch (err: any) {
    // If docker internal failed, try local fallback
    if (baseBackend !== BACKEND_LOCAL_URL) {
      try {
        const fallbackUrl = `${BACKEND_LOCAL_URL}/api/${params.path}${url.search}`;
        const res = await fetch(fallbackUrl, {
          method: request.method,
          headers,
          body,
          redirect: 'manual'
        });
        return new Response(res.body, {
          status: res.status,
          statusText: res.statusText,
          headers: new Headers(res.headers)
        });
      } catch (e2) {
        // ignore
      }
    }

    return new Response(
      JSON.stringify({
        success: false,
        message: `Gagal terhubung ke backend server (${err.message})`
      }),
      {
        status: 502,
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}

export const GET: RequestHandler = proxyRequest;
export const POST: RequestHandler = proxyRequest;
export const PUT: RequestHandler = proxyRequest;
export const DELETE: RequestHandler = proxyRequest;
export const PATCH: RequestHandler = proxyRequest;
export const fallback: RequestHandler = proxyRequest;
