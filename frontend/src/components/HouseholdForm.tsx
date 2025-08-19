'use client';

import React, { useState } from 'react';
import { useWizard } from '@/context/WizardContext';
import type { HouseholdLineDTO } from '@/types/api';

interface HouseholdFormProps {
  onValidationChange?: (isValid: boolean) => void;
}

export function HouseholdForm({ onValidationChange }: HouseholdFormProps) {
  const { state, updateHousehold, addHouseholdLine, removeHouseholdLine, updateHouseholdLine } = useWizard();
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Generate unique line ID
  const generateLineId = () => `line_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

  // Validate individual line
  const validateLine = (line: HouseholdLineDTO): Record<string, string> => {
    const lineErrors: Record<string, string> = {};
    
    if (!line.expected_gb || line.expected_gb < 0) {
      lineErrors[`${line.line_id}_gb`] = 'Expected GB must be a positive number';
    }
    
    if (!line.expected_min || line.expected_min < 0) {
      lineErrors[`${line.line_id}_min`] = 'Expected minutes must be a positive number';
    }
    
    if (line.tv_hd_hours < 0) {
      lineErrors[`${line.line_id}_tv`] = 'TV HD hours cannot be negative';
    }
    
    return lineErrors;
  };

  // Validate all lines and update validation state
  const validateAllLines = (household: HouseholdLineDTO[]) => {
    const allErrors: Record<string, string> = {};
    
    if (household.length === 0) {
      allErrors.general = 'At least one household line is required';
    }
    
    household.forEach(line => {
      const lineErrors = validateLine(line);
      Object.assign(allErrors, lineErrors);
    });
    
    setErrors(allErrors);
    const isValid = Object.keys(allErrors).length === 0;
    onValidationChange?.(isValid);
    return isValid;
  };

  // Add new line
  const handleAddLine = () => {
    const newLine: HouseholdLineDTO = {
      line_id: generateLineId(),
      expected_gb: 10,
      expected_min: 300,
      tv_hd_hours: 0,
    };
    
    addHouseholdLine(newLine);
    
    // Validate after adding
    setTimeout(() => validateAllLines([...state.household, newLine]), 0);
  };

  // Remove line
  const handleRemoveLine = (lineId: string) => {
    removeHouseholdLine(lineId);
    
    const updatedHousehold = state.household.filter(line => line.line_id !== lineId);
    setTimeout(() => validateAllLines(updatedHousehold), 0);
  };

  // Update line field
  const handleLineUpdate = (lineId: string, field: keyof HouseholdLineDTO, value: number) => {
    updateHouseholdLine(lineId, { [field]: value });
    
    // Validate after update
    const updatedHousehold = state.household.map(line =>
      line.line_id === lineId ? { ...line, [field]: value } : line
    );
    setTimeout(() => validateAllLines(updatedHousehold), 0);
  };

  // Initialize with one line if empty
  React.useEffect(() => {
    if (state.household.length === 0) {
      handleAddLine();
    } else {
      validateAllLines(state.household);
    }
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Household Information</h2>
        <p className="text-gray-600 mt-2">
          Add information about each mobile line in your household.
        </p>
      </div>

      {errors.general && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-md">
          <p className="text-red-700 text-sm">{errors.general}</p>
        </div>
      )}

      <div className="space-y-4">
        {state.household.map((line, index) => (
          <div key={line.line_id} className="p-4 border border-gray-200 rounded-lg bg-gray-50">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-medium text-gray-900">
                Line {index + 1}
              </h3>
              {state.household.length > 1 && (
                <button
                  type="button"
                  onClick={() => handleRemoveLine(line.line_id)}
                  className="px-3 py-1 text-sm text-red-600 hover:text-red-800 hover:bg-red-50 rounded-md transition-colors"
                >
                  Remove
                </button>
              )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* Expected GB */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Expected GB per month
                </label>
                <input
                  type="number"
                  min="0"
                  step="0.1"
                  value={line.expected_gb}
                  onChange={(e) => handleLineUpdate(line.line_id, 'expected_gb', parseFloat(e.target.value) || 0)}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors[`${line.line_id}_gb`] ? 'border-red-300' : 'border-gray-300'
                  }`}
                  placeholder="10"
                />
                {errors[`${line.line_id}_gb`] && (
                  <p className="text-red-600 text-sm mt-1">{errors[`${line.line_id}_gb`]}</p>
                )}
              </div>

              {/* Expected Minutes */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Expected minutes per month
                </label>
                <input
                  type="number"
                  min="0"
                  value={line.expected_min}
                  onChange={(e) => handleLineUpdate(line.line_id, 'expected_min', parseInt(e.target.value) || 0)}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors[`${line.line_id}_min`] ? 'border-red-300' : 'border-gray-300'
                  }`}
                  placeholder="300"
                />
                {errors[`${line.line_id}_min`] && (
                  <p className="text-red-600 text-sm mt-1">{errors[`${line.line_id}_min`]}</p>
                )}
              </div>

              {/* TV HD Hours */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  TV HD hours per month
                </label>
                <input
                  type="number"
                  min="0"
                  value={line.tv_hd_hours}
                  onChange={(e) => handleLineUpdate(line.line_id, 'tv_hd_hours', parseInt(e.target.value) || 0)}
                  className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${
                    errors[`${line.line_id}_tv`] ? 'border-red-300' : 'border-gray-300'
                  }`}
                  placeholder="0"
                />
                {errors[`${line.line_id}_tv`] && (
                  <p className="text-red-600 text-sm mt-1">{errors[`${line.line_id}_tv`]}</p>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Add Line Button */}
      <button
        type="button"
        onClick={handleAddLine}
        className="w-full py-2 px-4 border-2 border-dashed border-gray-300 rounded-md text-gray-600 hover:border-gray-400 hover:text-gray-700 transition-colors"
      >
        + Add Another Line
      </button>

      {/* Summary */}
      <div className="p-4 bg-blue-50 border border-blue-200 rounded-md">
        <h4 className="text-sm font-medium text-blue-900 mb-2">Summary</h4>
        <div className="text-sm text-blue-700 space-y-1">
          <p>Total lines: {state.household.length}</p>
          <p>Total expected GB: {state.household.reduce((sum, line) => sum + line.expected_gb, 0).toFixed(1)} GB</p>
          <p>Total expected minutes: {state.household.reduce((sum, line) => sum + line.expected_min, 0).toLocaleString()}</p>
          <p>Total TV HD hours: {state.household.reduce((sum, line) => sum + line.tv_hd_hours, 0)}</p>
        </div>
      </div>
    </div>
  );
}
