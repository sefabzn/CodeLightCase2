'use client';

import React, { useState, useEffect } from 'react';
import { useInstallSlots } from '@/lib/hooks';
import type { InstallSlot } from '@/types/api';

interface SlotPickerProps {
  addressId: string;
  tech: string;
  selectedSlotId?: string;
  onChange: (slotId: string) => void;
  className?: string;
}

export function SlotPicker({ 
  addressId, 
  tech, 
  selectedSlotId, 
  onChange, 
  className = "" 
}: SlotPickerProps) {
  const [selectedSlot, setSelectedSlot] = useState<string>(selectedSlotId || '');
  
  const { data: slotsData, isLoading, isError, error } = useInstallSlots(addressId, tech);

  // Update local state when prop changes
  useEffect(() => {
    setSelectedSlot(selectedSlotId || '');
  }, [selectedSlotId]);

  const handleSlotSelect = (slotId: string) => {
    setSelectedSlot(slotId);
    onChange(slotId);
  };

  const formatDateTime = (dateTimeString: string) => {
    const date = new Date(dateTimeString);
    return {
      date: date.toLocaleDateString('tr-TR', { 
        weekday: 'long', 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric' 
      }),
      time: date.toLocaleTimeString('tr-TR', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    };
  };

  const getSlotDuration = (startTime: string, endTime: string) => {
    const start = new Date(startTime);
    const end = new Date(endTime);
    const durationMs = end.getTime() - start.getTime();
    const hours = Math.floor(durationMs / (1000 * 60 * 60));
    const minutes = Math.floor((durationMs % (1000 * 60 * 60)) / (1000 * 60));
    
    if (hours > 0 && minutes > 0) {
      return `${hours}h ${minutes}m`;
    } else if (hours > 0) {
      return `${hours}h`;
    } else {
      return `${minutes}m`;
    }
  };

  const groupSlotsByDate = (slots: InstallSlot[]) => {
    const grouped = slots.reduce((acc, slot) => {
      const date = new Date(slot.slot_start).toDateString();
      if (!acc[date]) {
        acc[date] = [];
      }
      acc[date].push(slot);
      return acc;
    }, {} as Record<string, InstallSlot[]>);

    // Sort dates and slots within each date
    Object.keys(grouped).forEach(date => {
      grouped[date].sort((a, b) => new Date(a.slot_start).getTime() - new Date(b.slot_start).getTime());
    });

    return Object.entries(grouped).sort(([dateA], [dateB]) => 
      new Date(dateA).getTime() - new Date(dateB).getTime()
    );
  };

  if (isLoading) {
    return (
      <div className={`space-y-4 ${className}`}>
        <h3 className="text-lg font-semibold text-gray-900">Installation Slots</h3>
        <div className="flex items-center space-x-3">
          <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
          <span className="text-gray-600">Loading available installation slots...</span>
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className={`space-y-4 ${className}`}>
        <h3 className="text-lg font-semibold text-gray-900">Installation Slots</h3>
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <div className="flex items-center">
            <svg className="w-5 h-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <h4 className="text-red-800 font-medium">Unable to load installation slots</h4>
              <p className="text-red-700 text-sm mt-1">
                {error instanceof Error ? error.message : 'Please try again later'}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const availableSlots = slotsData?.slots.filter(slot => slot.available) || [];
  const groupedSlots = groupSlotsByDate(availableSlots);

  if (availableSlots.length === 0) {
    return (
      <div className={`space-y-4 ${className}`}>
        <h3 className="text-lg font-semibold text-gray-900">Installation Slots</h3>
        <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4">
          <div className="flex items-center">
            <svg className="w-5 h-5 text-yellow-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <div>
              <h4 className="text-yellow-800 font-medium">No slots available</h4>
              <p className="text-yellow-700 text-sm mt-1">
                No installation slots are currently available for {tech.toUpperCase()} technology at this address.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      <div>
        <h3 className="text-lg font-semibold text-gray-900">Choose Installation Slot</h3>
        <p className="text-gray-600 text-sm mt-1">
          Select a convenient time for {tech.toUpperCase()} installation at {addressId}
        </p>
      </div>

      <div className="space-y-6">
        {groupedSlots.map(([dateStr, slots]) => {
          const sampleDate = formatDateTime(slots[0].slot_start);
          
          return (
            <div key={dateStr} className="space-y-3">
              <h4 className="font-medium text-gray-900 border-b border-gray-200 pb-2">
                {sampleDate.date}
              </h4>
              
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {slots.map((slot) => {
                  const startTime = formatDateTime(slot.slot_start);
                  const endTime = formatDateTime(slot.slot_end);
                  const duration = getSlotDuration(slot.slot_start, slot.slot_end);
                  const isSelected = selectedSlot === slot.slot_id;
                  
                  return (
                    <button
                      key={slot.slot_id}
                      onClick={() => handleSlotSelect(slot.slot_id)}
                      className={`relative p-4 border-2 rounded-lg text-left transition-all ${
                        isSelected
                          ? 'border-blue-600 bg-blue-50 shadow-md'
                          : 'border-gray-200 bg-white hover:border-gray-300 hover:shadow-sm'
                      }`}
                    >
                      {/* Selected indicator */}
                      {isSelected && (
                        <div className="absolute top-2 right-2">
                          <div className="w-5 h-5 bg-blue-600 rounded-full flex items-center justify-center">
                            <svg className="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                          </div>
                        </div>
                      )}

                      <div className="space-y-2">
                        <div className={`font-medium ${isSelected ? 'text-blue-900' : 'text-gray-900'}`}>
                          {startTime.time} - {endTime.time}
                        </div>
                        
                        <div className={`text-sm ${isSelected ? 'text-blue-700' : 'text-gray-600'}`}>
                          Duration: {duration}
                        </div>

                        <div className="flex items-center space-x-2">
                          <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                            isSelected 
                              ? 'bg-blue-100 text-blue-800' 
                              : 'bg-green-100 text-green-800'
                          }`}>
                            Available
                          </span>
                          
                          <span className={`text-xs ${isSelected ? 'text-blue-600' : 'text-gray-500'}`}>
                            {tech.toUpperCase()}
                          </span>
                        </div>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          );
        })}
      </div>

      {/* Selected slot summary */}
      {selectedSlot && (
        <div className="bg-green-50 border border-green-200 rounded-md p-4">
          <div className="flex items-center">
            <svg className="w-5 h-5 text-green-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            <div>
              <h4 className="text-green-800 font-medium">Installation slot selected</h4>
              <p className="text-green-700 text-sm mt-1">
                {(() => {
                  const slot = availableSlots.find(s => s.slot_id === selectedSlot);
                  if (!slot) return '';
                  
                  const start = formatDateTime(slot.slot_start);
                  const end = formatDateTime(slot.slot_end);
                  const duration = getSlotDuration(slot.slot_start, slot.slot_end);
                  
                  return `${start.date} from ${start.time} to ${end.time} (${duration})`;
                })()}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Help text */}
      <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
        <h4 className="text-sm font-medium text-blue-900 mb-2">💡 Installation Information</h4>
        <div className="text-sm text-blue-700 space-y-1">
          <p>• Installation requires someone to be present at the address</p>
          <p>• {tech.toUpperCase()} installation typically takes {tech === 'fiber' ? '2-4 hours' : '1-2 hours'}</p>
          <p>• You'll receive a confirmation SMS and email with installer contact details</p>
          <p>• Installation can be rescheduled up to 24 hours before the appointment</p>
        </div>
      </div>
    </div>
  );
}
