import {
  RecommendationRequest,
  RecommendationResponse,
  CoverageInfo,
  InstallSlotsResponse,
  CheckoutRequest,
  CheckoutResponse,
  HealthResponse,
  ErrorResponse,
} from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
    public details?: string[]
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    let errorResponse: ErrorResponse;
    try {
      errorResponse = await response.json();
    } catch {
      throw new ApiError(
        response.status,
        'NETWORK_ERROR',
        `Network error: ${response.status} ${response.statusText}`
      );
    }

    throw new ApiError(
      response.status,
      errorResponse.error.code,
      errorResponse.error.message,
      errorResponse.error.details
    );
  }

  return response.json();
}

export const apiClient = {
  // Health check
  async health(): Promise<HealthResponse> {
    return fetchApi<HealthResponse>('/health');
  },

  // Get coverage information for an address
  async getCoverage(addressId: string): Promise<CoverageInfo> {
    return fetchApi<CoverageInfo>(`/api/coverage/${addressId}`);
  },

  // Get install slots for an address and technology
  async getInstallSlots(addressId: string, tech = 'fiber'): Promise<InstallSlotsResponse> {
    return fetchApi<InstallSlotsResponse>(`/api/install-slots/${addressId}?tech=${tech}`);
  },

  // Get recommendations
  async getRecommendations(request: RecommendationRequest): Promise<RecommendationResponse> {
    return fetchApi<RecommendationResponse>('/api/recommendation', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  },

  // Process checkout
  async checkout(request: CheckoutRequest): Promise<CheckoutResponse> {
    return fetchApi<CheckoutResponse>('/api/checkout', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  },
};

export { ApiError };
export type { ErrorResponse };
