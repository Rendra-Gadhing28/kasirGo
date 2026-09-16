import type { APIResponse } from '$lib/types';

const BASE_URL = '/api/v1';

class APIClient {
  private getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('kasirpro_token');
    }
    return null;
  }

  async request<T = any>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    const url = endpoint.startsWith('http') ? endpoint : `${BASE_URL}${endpoint}`;
    const token = this.getToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...(options.headers as Record<string, string> || {})
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: 'include' // allows HttpOnly cookies to pass
    });

    const contentType = response.headers.get('content-type') || '';
    let data: APIResponse<T>;

    if (contentType.includes('application/json')) {
      data = await response.json().catch(() => ({
        success: false,
        message: `HTTP ${response.status}: Gagal memproses data JSON server.`
      }));
    } else {
      const text = await response.text();
      data = {
        success: response.ok,
        message: response.ok ? 'OK' : `HTTP ${response.status}: Layanan API tidak ditemukan (${text.slice(0, 80)}...)`
      };
    }

    if (!response.ok || !data.success) {
      let errorMsg = data.message || 'Terjadi kesalahan sistem.';
      if (data.errors && Array.isArray(data.errors) && data.errors.length > 0) {
        errorMsg += ': ' + data.errors.map((e: any) => e.message || e.field).join(', ');
      }
      throw new Error(errorMsg);
    }

    return data;
  }

  get<T = any>(endpoint: string) {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  post<T = any>(endpoint: string, body?: any) {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined
    });
  }

  put<T = any>(endpoint: string, body?: any) {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined
    });
  }

  delete<T = any>(endpoint: string) {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

export const api = new APIClient();
