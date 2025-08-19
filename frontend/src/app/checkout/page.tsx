'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useWizard } from '@/context/WizardContext';
import { useCheckout } from '@/lib/hooks';
import { SlotPicker } from '@/components/SlotPicker';
import type { RecommendationCandidateDTO, CheckoutRequest } from '@/types/api';

export default function CheckoutPage() {
  const router = useRouter();
  const { state } = useWizard();
  const [selectedRecommendation, setSelectedRecommendation] = useState<RecommendationCandidateDTO | null>(null);
  const [selectedSlotId, setSelectedSlotId] = useState<string>('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [orderCompleted, setOrderCompleted] = useState<string | null>(null);

  const checkoutMutation = useCheckout();

  // Load selected recommendation from session storage
  useEffect(() => {
    const savedRecommendation = sessionStorage.getItem('selectedRecommendation');
    if (savedRecommendation) {
      try {
        const recommendation = JSON.parse(savedRecommendation);
        setSelectedRecommendation(recommendation);
      } catch (error) {
        console.error('Failed to parse saved recommendation:', error);
        router.push('/recommendations');
      }
    } else {
      router.push('/recommendations');
    }
  }, [router]);

  // Redirect if no required data
  useEffect(() => {
    if (!state.household.length || !state.addressId) {
      router.push('/');
    }
  }, [state, router]);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('tr-TR', {
      style: 'currency',
      currency: 'TRY',
      minimumFractionDigits: 2,
    }).format(price);
  };

  const handleSlotChange = (slotId: string) => {
    setSelectedSlotId(slotId);
  };

  const handleConfirmOrder = async () => {
    if (!selectedRecommendation || !selectedSlotId || !state.addressId) {
      return;
    }

    setIsProcessing(true);

    const checkoutRequest: CheckoutRequest = {
      user_id: state.userId || 1,
      selected_combo: selectedRecommendation,
      slot_id: selectedSlotId,
      address_id: state.addressId,
    };

    try {
      const result = await checkoutMutation.mutateAsync(checkoutRequest);
      setOrderCompleted(result.order_id);
      
      // Clear the selected recommendation from session storage
      sessionStorage.removeItem('selectedRecommendation');
    } catch (error) {
      console.error('Checkout failed:', error);
    } finally {
      setIsProcessing(false);
    }
  };

  const handleBackToRecommendations = () => {
    router.push('/recommendations');
  };

  const handleStartOver = () => {
    sessionStorage.removeItem('selectedRecommendation');
    router.push('/');
  };

  // Show success page if order is completed
  if (orderCompleted) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Order Confirmed!</h1>
          
          <div className="space-y-4 mb-8">
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="text-sm text-gray-600 mb-1">Order ID</div>
              <div className="text-lg font-mono font-bold text-gray-900">{orderCompleted}</div>
            </div>
            
            <p className="text-gray-600">
              Your Turkcell package has been successfully ordered. You'll receive a confirmation email and SMS shortly.
            </p>
            
            <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
              <div className="text-sm text-blue-800">
                <div className="font-medium mb-1">What's next?</div>
                <ul className="text-left space-y-1">
                  <li>• You'll receive installation details via SMS/email</li>
                  <li>• Our technician will contact you before the appointment</li>
                  <li>• Installation will be completed on your selected date</li>
                  <li>• Services will be activated immediately after installation</li>
                </ul>
              </div>
            </div>
          </div>
          
          <div className="space-y-3">
            <button
              onClick={() => router.push('/')}
              className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Get Another Recommendation
            </button>
            
            <button
              onClick={() => window.print()}
              className="w-full px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors"
            >
              Print Confirmation
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!selectedRecommendation) {
    return null; // Will redirect
  }

  const canProceed = selectedSlotId && !isProcessing;
  const homeTech = selectedRecommendation.items.home?.tech || 'fiber';

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Checkout</h1>
              <p className="mt-2 text-gray-600">
                Complete your order and schedule installation
              </p>
            </div>
            
            <button
              onClick={handleBackToRecommendations}
              className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md text-gray-700 bg-white hover:bg-gray-50 transition-colors"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Recommendations
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
                  <div className="flex items-center text-green-600">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-green-600 bg-green-600 text-white">
                      ✓
                    </span>
                    <span className="ml-4 text-sm font-medium">Recommendations</span>
                  </div>
                </li>
                <li className="flex items-center">
                  <div className="flex items-center text-blue-600">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-blue-600 bg-blue-600 text-white">
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
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Order Summary */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow-sm p-6 sticky top-8">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Order Summary</h2>
              
              <div className="space-y-4">
                <div>
                  <h3 className="font-medium text-gray-900">{selectedRecommendation.combo_label}</h3>
                  <div className="text-sm text-gray-600 mt-1">
                    {selectedRecommendation.items.mobile.length} mobile line{selectedRecommendation.items.mobile.length > 1 ? 's' : ''}
                    {selectedRecommendation.items.home && ', home internet'}
                    {selectedRecommendation.items.tv && ', TV package'}
                  </div>
                </div>

                <div className="border-t border-gray-200 pt-4">
                  <div className="flex justify-between text-sm">
                    <span>Subtotal</span>
                    <span>{formatPrice(selectedRecommendation.monthly_total + selectedRecommendation.discounts.total_discount)}</span>
                  </div>
                  
                  {selectedRecommendation.discounts.total_discount > 0 && (
                    <div className="flex justify-between text-sm text-green-600 mt-1">
                      <span>Total Discounts</span>
                      <span>-{formatPrice(selectedRecommendation.discounts.total_discount)}</span>
                    </div>
                  )}
                </div>

                <div className="border-t border-gray-200 pt-4">
                  <div className="flex justify-between text-lg font-bold">
                    <span>Monthly Total</span>
                    <span>{formatPrice(selectedRecommendation.monthly_total)}</span>
                  </div>
                </div>

                {selectedRecommendation.items.home && (
                  <div className="bg-yellow-50 border border-yellow-200 rounded-md p-3">
                    <div className="text-sm text-yellow-800">
                      <div className="font-medium">One-time Installation Fee</div>
                      <div>{formatPrice(selectedRecommendation.items.home.install_fee)}</div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Installation Scheduling */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <SlotPicker
                addressId={state.addressId}
                tech={homeTech}
                selectedSlotId={selectedSlotId}
                onChange={handleSlotChange}
              />
            </div>

            {/* Confirm Order */}
            <div className="bg-white rounded-lg shadow-sm p-6 mt-8">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div className="mb-4 sm:mb-0">
                  <h3 className="text-lg font-medium text-gray-900">Ready to confirm?</h3>
                  <p className="text-sm text-gray-600">
                    Your order will be processed and installation will be scheduled.
                  </p>
                </div>

                <button
                  onClick={handleConfirmOrder}
                  disabled={!canProceed}
                  className={`w-full sm:w-auto inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md transition-colors ${
                    canProceed
                      ? 'text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500'
                      : 'text-gray-400 bg-gray-100 cursor-not-allowed'
                  }`}
                >
                  {isProcessing ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                      Processing Order...
                    </>
                  ) : (
                    <>
                      Confirm Order
                      <svg className="ml-2 w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                    </>
                  )}
                </button>
              </div>

              {/* Validation message */}
              {!selectedSlotId && (
                <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                  <p className="text-sm text-yellow-700">
                    Please select an installation slot to continue.
                  </p>
                </div>
              )}

              {/* Error message */}
              {checkoutMutation.isError && (
                <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
                  <p className="text-sm text-red-700">
                    Order failed. Please try again or contact support.
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
