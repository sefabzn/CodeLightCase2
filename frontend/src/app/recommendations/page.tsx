'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useWizard } from '@/context/WizardContext';
import { useRecommendation } from '@/lib/hooks';
import { RecommendationCard } from '@/components/RecommendationCard';
import { SummaryModal } from '@/components/SummaryModal';
import type { RecommendationCandidateDTO, RecommendationRequest } from '@/types/api';

export default function RecommendationsPage() {
  const router = useRouter();
  const { state } = useWizard();
  const [selectedRecommendation, setSelectedRecommendation] = useState<RecommendationCandidateDTO | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [modalRecommendation, setModalRecommendation] = useState<RecommendationCandidateDTO | null>(null);

  const recommendationMutation = useRecommendation();

  // Check if we have the required data
  const hasRequiredData = state.household.length > 0 && state.addressId;

  // Create the request object
  const recommendationRequest: RecommendationRequest | null = hasRequiredData ? {
    user_id: state.userId || 1, // Default user ID for now
    address_id: state.addressId,
    household: state.household,
    prefer_tech: state.preferTech,
  } : null;

  // Fetch recommendations on mount
  useEffect(() => {
    if (recommendationRequest && !recommendationMutation.data && !recommendationMutation.isPending) {
      recommendationMutation.mutate(recommendationRequest);
    }
  }, [recommendationRequest]);

  // Redirect if no data
  useEffect(() => {
    if (!hasRequiredData) {
      router.push('/');
    }
  }, [hasRequiredData, router]);

  const handleShowDetails = (recommendation: RecommendationCandidateDTO) => {
    setModalRecommendation(recommendation);
    setShowModal(true);
  };

  const handleSelectRecommendation = (recommendation: RecommendationCandidateDTO) => {
    setSelectedRecommendation(recommendation);
    
    // Navigate to checkout with the selected recommendation
    // For now, we'll store it in sessionStorage and navigate
    sessionStorage.setItem('selectedRecommendation', JSON.stringify(recommendation));
    router.push('/checkout');
  };

  const handleRetry = () => {
    if (recommendationRequest) {
      recommendationMutation.mutate(recommendationRequest);
    }
  };

  const handleBackToSetup = () => {
    router.push('/');
  };

  if (!hasRequiredData) {
    return null; // Will redirect
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                Your Recommendations
              </h1>
              <p className="mt-2 text-gray-600">
                We've analyzed your needs and found the best package combinations for you
              </p>
            </div>
            
            <button
              onClick={handleBackToSetup}
              className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md text-gray-700 bg-white hover:bg-gray-50 transition-colors"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Setup
            </button>
          </div>
          
          {/* Progress indicator */}
          <div className="mt-8">
            <nav aria-label="Progress">
              <ol className="flex items-center justify-center space-x-5">
                <li className="flex items-center">
                  <div className="flex items-center text-green-600">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-green-600 bg-green-600 text-white">
                      ✓
                    </span>
                    <span className="ml-4 text-sm font-medium">Setup</span>
                  </div>
                </li>
                <li className="flex items-center">
                  <div className="flex items-center text-blue-600">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-blue-600 bg-blue-600 text-white">
                      2
                    </span>
                    <span className="ml-4 text-sm font-medium">Recommendations</span>
                  </div>
                </li>
                <li className="flex items-center">
                  <div className="flex items-center text-gray-500">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-gray-300">
                      3
                    </span>
                    <span className="ml-4 text-sm font-medium">Checkout</span>
                  </div>
                </li>
              </ol>
            </nav>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Household Summary */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Your Requirements</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <div className="font-medium text-gray-700">Address</div>
              <div className="text-gray-600">{state.addressId}</div>
            </div>
            <div>
              <div className="font-medium text-gray-700">Household Lines</div>
              <div className="text-gray-600">{state.household.length} mobile line{state.household.length > 1 ? 's' : ''}</div>
            </div>
            <div>
              <div className="font-medium text-gray-700">Total Usage</div>
              <div className="text-gray-600">
                {state.household.reduce((sum, line) => sum + line.expected_gb, 0).toFixed(1)} GB, {' '}
                {state.household.reduce((sum, line) => sum + line.expected_min, 0).toLocaleString()} min
              </div>
            </div>
          </div>
        </div>

        {/* Loading State */}
        {recommendationMutation.isPending && (
          <div className="bg-white rounded-lg shadow-sm p-12">
            <div className="text-center">
              <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Finding Your Perfect Package</h3>
              <p className="text-gray-600">
                We're analyzing thousands of combinations to find the best deals for you...
              </p>
            </div>
          </div>
        )}

        {/* Error State */}
        {recommendationMutation.isError && (
          <div className="bg-white rounded-lg shadow-sm p-8">
            <div className="text-center">
              <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Unable to Load Recommendations</h3>
              <p className="text-gray-600 mb-4">
                There was an error while fetching your personalized recommendations. Please try again.
              </p>
              <div className="space-x-3">
                <button
                  onClick={handleRetry}
                  className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                >
                  Try Again
                </button>
                <button
                  onClick={handleBackToSetup}
                  className="inline-flex items-center px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors"
                >
                  Back to Setup
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Recommendations */}
        {recommendationMutation.data?.top3 && (
          <div className="space-y-6">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                {recommendationMutation.data.top3.length} Package{recommendationMutation.data.top3.length > 1 ? 's' : ''} Found
              </h2>
              <p className="text-gray-600">
                Sorted by best value for your specific needs
              </p>
            </div>

            <div className="grid gap-6 lg:grid-cols-1 xl:grid-cols-2 2xl:grid-cols-3">
              {recommendationMutation.data.top3.map((recommendation, index) => (
                <RecommendationCard
                  key={`${recommendation.combo_label}-${index}`}
                  recommendation={recommendation}
                  rank={index + 1}
                  onSelect={handleSelectRecommendation}
                  onShowDetails={handleShowDetails}
                />
              ))}
            </div>

            {/* Help Text */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mt-8">
              <h3 className="text-sm font-medium text-blue-900 mb-2">💡 How to choose:</h3>
              <div className="text-sm text-blue-700 space-y-1">
                <p>• <strong>Best Value:</strong> Optimized for your budget with maximum savings</p>
                <p>• <strong>Popular Choice:</strong> Most commonly selected by similar households</p>
                <p>• <strong>Premium Option:</strong> Higher-tier services with premium features</p>
                <p>• Click "View Details" to see the complete cost breakdown</p>
                <p>• All packages include bundle discounts where applicable</p>
              </div>
            </div>
          </div>
        )}

        {/* Empty State */}
        {recommendationMutation.data?.top3.length === 0 && (
          <div className="bg-white rounded-lg shadow-sm p-8">
            <div className="text-center">
              <div className="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">No Recommendations Available</h3>
              <p className="text-gray-600 mb-4">
                We couldn't find suitable packages for your requirements. Please try adjusting your preferences.
              </p>
              <button
                onClick={handleBackToSetup}
                className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
              >
                Modify Requirements
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Summary Modal */}
      <SummaryModal
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        recommendation={modalRecommendation}
      />
    </div>
  );
}
