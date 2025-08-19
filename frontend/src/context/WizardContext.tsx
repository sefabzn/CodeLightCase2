'use client';

import React, { createContext, useContext, useState, ReactNode } from 'react';
import type { HouseholdLineDTO } from '@/types/api';

export interface WizardState {
  userId: number | null;
  addressId: string;
  household: HouseholdLineDTO[];
  preferTech: string[];
}

export interface WizardContextType {
  state: WizardState;
  updateUserId: (userId: number) => void;
  updateAddressId: (addressId: string) => void;
  updateHousehold: (household: HouseholdLineDTO[]) => void;
  updatePreferTech: (tech: string[]) => void;
  addHouseholdLine: (line: HouseholdLineDTO) => void;
  removeHouseholdLine: (lineId: string) => void;
  updateHouseholdLine: (lineId: string, updates: Partial<HouseholdLineDTO>) => void;
  reset: () => void;
}

const defaultState: WizardState = {
  userId: null,
  addressId: '',
  household: [],
  preferTech: ['fiber', 'vdsl', 'fwa'], // Default preference order
};

const WizardContext = createContext<WizardContextType | undefined>(undefined);

export function WizardProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<WizardState>(defaultState);

  const updateUserId = (userId: number) => {
    setState(prev => ({ ...prev, userId }));
  };

  const updateAddressId = (addressId: string) => {
    setState(prev => ({ ...prev, addressId }));
  };

  const updateHousehold = (household: HouseholdLineDTO[]) => {
    setState(prev => ({ ...prev, household }));
  };

  const updatePreferTech = (preferTech: string[]) => {
    setState(prev => ({ ...prev, preferTech }));
  };

  const addHouseholdLine = (line: HouseholdLineDTO) => {
    setState(prev => ({
      ...prev,
      household: [...prev.household, line]
    }));
  };

  const removeHouseholdLine = (lineId: string) => {
    setState(prev => ({
      ...prev,
      household: prev.household.filter(line => line.line_id !== lineId)
    }));
  };

  const updateHouseholdLine = (lineId: string, updates: Partial<HouseholdLineDTO>) => {
    setState(prev => ({
      ...prev,
      household: prev.household.map(line =>
        line.line_id === lineId ? { ...line, ...updates } : line
      )
    }));
  };

  const reset = () => {
    setState(defaultState);
  };

  const value: WizardContextType = {
    state,
    updateUserId,
    updateAddressId,
    updateHousehold,
    updatePreferTech,
    addHouseholdLine,
    removeHouseholdLine,
    updateHouseholdLine,
    reset,
  };

  return (
    <WizardContext.Provider value={value}>
      {children}
    </WizardContext.Provider>
  );
}

export function useWizard(): WizardContextType {
  const context = useContext(WizardContext);
  if (context === undefined) {
    throw new Error('useWizard must be used within a WizardProvider');
  }
  return context;
}
