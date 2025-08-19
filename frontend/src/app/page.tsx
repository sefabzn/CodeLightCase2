'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useWizard } from '@/context/WizardContext';
import { HouseholdForm } from '@/components/HouseholdForm';
import { AddressForm } from '@/components/AddressForm';

export default function Home() {
  const router = useRouter();
  const { state } = useWizard();
  const [isHouseholdValid, setIsHouseholdValid] = useState(false);
  const [isAddressValid, setIsAddressValid] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const canContinue = isHouseholdValid && isAddressValid && state.household.length > 0 && state.addressId;

  const handleContinue = async () => {
    if (!canContinue) return;

    setIsSubmitting(true);
    
    try {
      // Small delay to show loading state
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Navigate to recommendations page
      router.push('/recommendations');
    } catch (error) {
      console.error('Navigation error:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center">
            <h1 className="text-3xl font-bold text-gray-900">
              Turkcell Package Recommendation
            </h1>
            <p className="mt-2 text-gray-600">
              Find the perfect mobile, home internet, and TV package for your household
            </p>
          </div>
          
          {/* Progress indicator */}
          <div className="mt-8">
            <nav aria-label="Progress">
              <ol className="flex items-center justify-center space-x-5">
                <li className="flex items-center">
                  <div className="flex items-center text-blue-600">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-blue-600 bg-blue-600 text-white">
                      1
                    </span>
                    <span className="ml-4 text-sm font-medium">Setup</span>
                  </div>
                </li>
                <li className="flex items-center">
                  <div className="flex items-center text-gray-500">
                    <span className="flex h-10 w-10 items-center justify-center rounded-full border-2 border-gray-300">
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
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-8">
          {/* Household Form */}
          <div className="bg-white rounded-lg shadow-sm p-6">
            <HouseholdForm onValidationChange={setIsHouseholdValid} />
          </div>

          {/* Address Form */}
          <div className="bg-white rounded-lg shadow-sm p-6">
            <AddressForm onValidationChange={setIsAddressValid} />
          </div>

          {/* Continue Button */}
          <div className="bg-white rounded-lg shadow-sm p-6">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
              <div className="mb-4 sm:mb-0">
                <h3 className="text-lg font-medium text-gray-900">Ready to continue?</h3>
                <p className="text-sm text-gray-600">
                  We'll analyze your needs and show you the best package recommendations.
                </p>
              </div>

              <button
                onClick={handleContinue}
                disabled={!canContinue || isSubmitting}
                className={`w-full sm:w-auto inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md transition-colors ${
                  canContinue
                    ? 'text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500'
                    : 'text-gray-400 bg-gray-100 cursor-not-allowed'
                }`}
              >
                {isSubmitting ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                    Processing...
                  </>
                ) : (
                  <>
                    Continue to Recommendations
                    <svg className="ml-2 w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </>
                )}
              </button>
            </div>

            {/* Validation summary */}
            {(!isHouseholdValid || !isAddressValid) && (
              <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                <div className="text-sm text-yellow-700">
                  <p className="font-medium mb-1">Please complete the following:</p>
                  <ul className="list-disc list-inside space-y-1">
                    {!isHouseholdValid && (
                      <li>Fill in valid household information for all lines</li>
                    )}
                    {!isAddressValid && (
                      <li>Provide a valid address with coverage information</li>
                    )}
                  </ul>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="bg-white border-t mt-12">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center text-sm text-gray-500">
            <p>© 2024 Turkcell Package Recommendation System</p>
            <p className="mt-1">Get personalized recommendations for mobile, home internet, and TV services</p>
          </div>
        </div>
      </div>
    </div>
  );
}