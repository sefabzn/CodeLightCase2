'use client';

import React, { useState, useEffect } from 'react';
import { useWizard } from '@/context/WizardContext';
import { useCoverage } from '@/lib/hooks';
import { CoverageBadge } from './CoverageBadge';

interface AddressFormProps {
  onValidationChange?: (isValid: boolean) => void;
}

export function AddressForm({ onValidationChange }: AddressFormProps) {
  const { state, updateAddressId } = useWizard();
  const [city, setCity] = useState('');
  const [district, setDistrict] = useState('');
  const [addressInput, setAddressInput] = useState(state.addressId);
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Get coverage data when address is provided
  const { data: coverage, isLoading: coverageLoading, isError: coverageError } = useCoverage(state.addressId);

  // Validation
  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!city.trim()) {
      newErrors.city = 'City is required';
    }

    if (!district.trim()) {
      newErrors.district = 'District is required';
    }

    if (!state.addressId.trim()) {
      newErrors.address_id = 'Address ID is required';
    } else if (coverageError) {
      newErrors.address_id = 'Invalid address ID - no coverage information found';
    }

    setErrors(newErrors);
    const isValid = Object.keys(newErrors).length === 0 && !coverageError;
    onValidationChange?.(isValid);
    return isValid;
  };

  // Handle address ID blur - trigger coverage lookup
  const handleAddressBlur = () => {
    if (addressInput.trim() && addressInput !== state.addressId) {
      updateAddressId(addressInput.trim());
    }
  };

  // Update validation when coverage data changes
  useEffect(() => {
    if (city && district && state.addressId) {
      validateForm();
    }
  }, [coverage, coverageError, city, district, state.addressId]);

  // Sample address IDs for dropdown
  const sampleAddresses = [
    { id: 'A1001', description: 'Istanbul, Kadıköy - Full Coverage (Fiber, VDSL, FWA)' },
    { id: 'A1002', description: 'Ankara, Çankaya - VDSL & FWA Coverage' },
    { id: 'A1003', description: 'Izmir, Konak - FWA Coverage Only' }
  ];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Address Information</h2>
        <p className="text-gray-600 mt-2">
          Provide your address details to check available services and coverage.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* City */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            City
          </label>
          <input
            type="text"
            value={city}
            onChange={(e) => setCity(e.target.value)}
            onBlur={validateForm}
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.city ? 'border-red-300' : 'border-gray-300'
            }`}
            placeholder="e.g., Istanbul"
          />
          {errors.city && (
            <p className="text-red-600 text-sm mt-1">{errors.city}</p>
          )}
        </div>

        {/* District */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            District
          </label>
          <input
            type="text"
            value={district}
            onChange={(e) => setDistrict(e.target.value)}
            onBlur={validateForm}
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.district ? 'border-red-300' : 'border-gray-300'
            }`}
            placeholder="e.g., Kadıköy"
          />
          {errors.district && (
            <p className="text-red-600 text-sm mt-1">{errors.district}</p>
          )}
        </div>
      </div>

      {/* Address ID */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Address ID
        </label>
        <div className="space-y-3">
          {/* Sample Addresses Dropdown */}
          <select
            value={addressInput}
            onChange={(e) => {
              setAddressInput(e.target.value);
              updateAddressId(e.target.value);
            }}
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.address_id ? 'border-red-300' : 'border-gray-300'
            }`}
          >
            <option value="">Select a sample address...</option>
            {sampleAddresses.map((addr) => (
              <option key={addr.id} value={addr.id}>
                {addr.id} - {addr.description}
              </option>
            ))}
          </select>

          {/* Or custom input */}
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="bg-white px-2 text-gray-500">or enter custom address ID</span>
            </div>
          </div>

          <input
            type="text"
            value={addressInput}
            onChange={(e) => setAddressInput(e.target.value)}
            onBlur={handleAddressBlur}
            className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
              errors.address_id ? 'border-red-300' : 'border-gray-300'
            }`}
            placeholder="e.g., A1001"
          />
        </div>
        
        {errors.address_id && (
          <p className="text-red-600 text-sm mt-1">{errors.address_id}</p>
        )}
      </div>

      {/* Coverage Information */}
      {state.addressId && (
        <div className="p-4 bg-gray-50 border border-gray-200 rounded-md">
          <h4 className="text-sm font-medium text-gray-700 mb-3">Coverage Information</h4>
          
          {coverageLoading && (
            <div className="flex items-center space-x-2">
              <div className="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
              <span className="text-sm text-gray-600">Checking coverage...</span>
            </div>
          )}

          {coverageError && (
            <div className="text-red-600 text-sm">
              Unable to find coverage information for this address ID.
            </div>
          )}

          {coverage && (
            <div className="space-y-3">
              <div className="text-sm text-gray-600">
                <p><strong>Address:</strong> {coverage.city}, {coverage.district}</p>
                <p><strong>Address ID:</strong> {coverage.address_id}</p>
              </div>
              
              <div>
                <p className="text-sm font-medium text-gray-700 mb-2">Available Technologies:</p>
                <CoverageBadge coverage={coverage} />
              </div>

              {coverage.available_tech.length === 0 && (
                <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                  <p className="text-yellow-800 text-sm">
                    No internet technologies are available at this address.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Help Text */}
      <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
        <h4 className="text-sm font-medium text-blue-900 mb-2">Address ID Help</h4>
        <div className="text-sm text-blue-700 space-y-1">
          <p>• Use the dropdown to select a sample address for testing</p>
          <p>• Address IDs are in format A#### (e.g., A1001, A1002)</p>
          <p>• Each address has different technology coverage (Fiber, VDSL, FWA)</p>
          <p>• Coverage information will be automatically loaded when you select an address</p>
        </div>
      </div>
    </div>
  );
}
