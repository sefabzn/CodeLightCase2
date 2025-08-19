'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from './api-client';
import type {
  RecommendationRequest,
  RecommendationResponse,
  CoverageInfo,
  InstallSlotsResponse,
  CheckoutRequest,
  CheckoutResponse,
} from '@/types/api';

// Query Keys
export const queryKeys = {
  health: () => ['health'] as const,
  coverage: (addressId: string) => ['coverage', addressId] as const,
  installSlots: (addressId: string, tech: string) => ['installSlots', addressId, tech] as const,
  recommendation: (input: RecommendationRequest) => ['recommendation', input] as const,
} as const;

// J1: useCatalog hook - combines coverage and install slots for an address
export function useCatalog(addressId: string) {
  // Get coverage information
  const coverageQuery = useQuery({
    queryKey: queryKeys.coverage(addressId),
    queryFn: () => apiClient.getCoverage(addressId),
    enabled: !!addressId,
    staleTime: 5 * 60 * 1000, // 5 minutes - coverage doesn't change often
  });

  // Get install slots for the first available tech (if coverage is available)
  const availableTech = coverageQuery.data?.available_tech?.[0] || 'fiber';
  const installSlotsQuery = useQuery({
    queryKey: queryKeys.installSlots(addressId, availableTech),
    queryFn: () => apiClient.getInstallSlots(addressId, availableTech),
    enabled: !!addressId && !!coverageQuery.data,
    staleTime: 2 * 60 * 1000, // 2 minutes - slots change more frequently
  });

  return {
    coverage: coverageQuery.data,
    slots: installSlotsQuery.data,
    isLoading: coverageQuery.isLoading || installSlotsQuery.isLoading,
    isError: coverageQuery.isError || installSlotsQuery.isError,
    error: coverageQuery.error || installSlotsQuery.error,
    refetch: () => {
      coverageQuery.refetch();
      installSlotsQuery.refetch();
    },
  };
}

// J2: useRecommendation hook - posts to API and returns top3
export function useRecommendation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: RecommendationRequest) => apiClient.getRecommendations(input),
    onSuccess: (data, variables) => {
      // Cache the result for potential re-use
      queryClient.setQueryData(queryKeys.recommendation(variables), data);
    },
    onError: (error) => {
      console.error('Recommendation request failed:', error);
    },
  });
}

// Helper hook to get cached recommendation data
export function useRecommendationData(input?: RecommendationRequest) {
  return useQuery({
    queryKey: input ? queryKeys.recommendation(input) : ['recommendation', 'empty'],
    queryFn: () => apiClient.getRecommendations(input!),
    enabled: false, // Only used for reading cached data
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

// J3: useCheckout hook - posts checkout
export function useCheckout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: CheckoutRequest) => apiClient.checkout(request),
    onSuccess: (data) => {
      // Invalidate relevant queries after successful checkout
      queryClient.invalidateQueries({ queryKey: ['installSlots'] });
      console.log('Checkout successful:', data.order_id);
    },
    onError: (error) => {
      console.error('Checkout failed:', error);
    },
  });
}

// Health check hook for monitoring
export function useHealth() {
  return useQuery({
    queryKey: queryKeys.health(),
    queryFn: () => apiClient.health(),
    staleTime: 30 * 1000, // 30 seconds
    refetchInterval: 60 * 1000, // Refetch every minute
    retry: 2,
  });
}

// Hook to get install slots for a specific tech
export function useInstallSlots(addressId: string, tech: string) {
  return useQuery({
    queryKey: queryKeys.installSlots(addressId, tech),
    queryFn: () => apiClient.getInstallSlots(addressId, tech),
    enabled: !!addressId && !!tech,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

// Hook specifically for coverage information
export function useCoverage(addressId: string) {
  return useQuery({
    queryKey: queryKeys.coverage(addressId),
    queryFn: () => apiClient.getCoverage(addressId),
    enabled: !!addressId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
