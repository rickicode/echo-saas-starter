const BASE_URL = '/api/v1'

interface RequestOptions {
  headers?: Record<string, string>
  signal?: AbortSignal
}

function getAuthHeaders(): Record<string, string> {
  const token = localStorage.getItem('access_token')
  if (token) {
    return { Authorization: `Bearer ${token}` }
  }
  return {}
}

async function handleResponse<T>(response: Response, retryFn?: () => Promise<Response>): Promise<T> {
  if (response.status === 401 && retryFn) {
    const refreshed = await attemptTokenRefresh()
    if (!refreshed) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/login'
      throw new Error('Unauthorized')
    }
    // Retry the original request with the new token
    const retryResponse = await retryFn()
    return handleResponse<T>(retryResponse)
  }

  if (response.status === 401) {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}

async function attemptTokenRefresh(): Promise<boolean> {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) {
    return false
  }

  try {
    const response = await fetch(`${BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })

    if (!response.ok) {
      return false
    }

    const data = await response.json()
    localStorage.setItem('access_token', data.access_token)
    if (data.refresh_token) {
      localStorage.setItem('refresh_token', data.refresh_token)
    }
    return true
  } catch {
    return false
  }
}

export async function get<T>(path: string, options?: RequestOptions): Promise<T> {
  const makeRequest = () =>
    fetch(`${BASE_URL}${path}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
        ...options?.headers,
      },
      signal: options?.signal,
    })

  const response = await makeRequest()
  return handleResponse<T>(response, makeRequest)
}

export async function post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
  const makeRequest = () =>
    fetch(`${BASE_URL}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
        ...options?.headers,
      },
      body: body ? JSON.stringify(body) : undefined,
      signal: options?.signal,
    })

  const response = await makeRequest()
  return handleResponse<T>(response, makeRequest)
}

export async function put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
  const makeRequest = () =>
    fetch(`${BASE_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
        ...options?.headers,
      },
      body: body ? JSON.stringify(body) : undefined,
      signal: options?.signal,
    })

  const response = await makeRequest()
  return handleResponse<T>(response, makeRequest)
}

export async function del<T>(path: string, options?: RequestOptions): Promise<T> {
  const makeRequest = () =>
    fetch(`${BASE_URL}${path}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
        ...options?.headers,
      },
      signal: options?.signal,
    })

  const response = await makeRequest()
  return handleResponse<T>(response, makeRequest)
}
