'use client';

import React from 'react';
import type { RecommendationCandidateDTO } from '@/types/api';

interface SummaryModalProps {
  isOpen: boolean;
  onClose: () => void;
  recommendation: RecommendationCandidateDTO | null;
}

export function SummaryModal({ isOpen, onClose, recommendation }: SummaryModalProps) {
  if (!isOpen || !recommendation) return null;

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('tr-TR', {
      style: 'currency',
      currency: 'TRY',
      minimumFractionDigits: 2,
    }).format(price);
  };

  // Calculate component costs
  const mobileTotal = recommendation.items.mobile.reduce((sum, mobile) => sum + mobile.line_cost, 0);
  const homeTotal = recommendation.items.home?.monthly_price || 0;
  const tvTotal = recommendation.items.tv?.monthly_price || 0;
  const subtotal = mobileTotal + homeTotal + tvTotal;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onClose}></div>
      
      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
          {/* Header */}
          <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center">
            <h2 className="text-2xl font-bold text-gray-900">
              Package Details: {recommendation.combo_label}
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div className="p-6 space-y-6">
            {/* Cost Summary */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h3 className="text-lg font-semibold text-blue-900 mb-3">Monthly Cost Breakdown</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Subtotal (before discounts)</span>
                  <span className="font-medium">{formatPrice(subtotal)}</span>
                </div>
                {recommendation.discounts.line_discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Multi-line discount</span>
                    <span>-{formatPrice(recommendation.discounts.line_discount)}</span>
                  </div>
                )}
                {recommendation.discounts.bundle_discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Bundle discount</span>
                    <span>-{formatPrice(recommendation.discounts.bundle_discount)}</span>
                  </div>
                )}
                <hr className="border-blue-300" />
                <div className="flex justify-between text-lg font-bold text-blue-900">
                  <span>Total Monthly Cost</span>
                  <span>{formatPrice(recommendation.monthly_total)}</span>
                </div>
                {recommendation.savings > 0 && (
                  <div className="text-center text-green-600 font-medium">
                    You save {formatPrice(recommendation.savings)} compared to individual plans
                  </div>
                )}
              </div>
            </div>

            {/* Mobile Plans Details */}
            {recommendation.items.mobile.length > 0 && (
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <span className="mr-2">📱</span>
                  Mobile Plans ({recommendation.items.mobile.length} line{recommendation.items.mobile.length > 1 ? 's' : ''})
                </h3>
                <div className="space-y-4">
                  {recommendation.items.mobile.map((mobile, index) => (
                    <div key={mobile.line_id} className="bg-gray-50 rounded-md p-3">
                      <div className="flex justify-between items-start mb-2">
                        <div>
                          <h4 className="font-medium text-gray-900">Line {index + 1}: {mobile.plan.plan_name}</h4>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>• {mobile.plan.quota_gb} GB data, {mobile.plan.quota_min.toLocaleString()} minutes</div>
                            {mobile.overage_gb > 0 && (
                              <div className="text-orange-600">• {mobile.overage_gb.toFixed(1)} GB overage</div>
                            )}
                            {mobile.overage_min > 0 && (
                              <div className="text-orange-600">• {mobile.overage_min} minute overage</div>
                            )}
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-semibold">{formatPrice(mobile.line_cost)}</div>
                          <div className="text-xs text-gray-500">per month</div>
                        </div>
                      </div>
                      
                      {(mobile.overage_gb > 0 || mobile.overage_min > 0) && (
                        <div className="text-xs text-gray-600 bg-yellow-50 p-2 rounded border-l-4 border-yellow-400">
                          <div className="font-medium text-yellow-800">Overage charges included:</div>
                          {mobile.overage_gb > 0 && (
                            <div>• Data: {mobile.overage_gb.toFixed(1)} GB × {formatPrice(mobile.plan.overage_gb)}/GB</div>
                          )}
                          {mobile.overage_min > 0 && (
                            <div>• Minutes: {mobile.overage_min} min × {formatPrice(mobile.plan.overage_min)}/min</div>
                          )}
                        </div>
                      )}
                    </div>
                  ))}
                  
                  <div className="bg-blue-50 p-3 rounded-md">
                    <div className="flex justify-between items-center">
                      <span className="font-medium text-blue-900">Mobile Subtotal</span>
                      <span className="font-bold text-blue-900">{formatPrice(mobileTotal)}</span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Home Internet Details */}
            {recommendation.items.home && (
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <span className="mr-2">🏠</span>
                  Home Internet
                </h3>
                <div className="bg-gray-50 rounded-md p-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-gray-900">{recommendation.items.home.name}</h4>
                      <div className="text-sm text-gray-600 space-y-1">
                        <div>• {recommendation.items.home.down_mbps} Mbps download speed</div>
                        <div>• Technology: {recommendation.items.home.tech.toUpperCase()}</div>
                        <div>• Installation fee: {formatPrice(recommendation.items.home.install_fee)}</div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold">{formatPrice(recommendation.items.home.monthly_price)}</div>
                      <div className="text-xs text-gray-500">per month</div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* TV Package Details */}
            {recommendation.items.tv && (
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <span className="mr-2">📺</span>
                  TV Package
                </h3>
                <div className="bg-gray-50 rounded-md p-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-medium text-gray-900">{recommendation.items.tv.name}</h4>
                      <div className="text-sm text-gray-600">
                        <div>• {recommendation.items.tv.hd_hours_included} HD hours included per month</div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold">{formatPrice(recommendation.items.tv.monthly_price)}</div>
                      <div className="text-xs text-gray-500">per month</div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Reasoning */}
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <h3 className="text-lg font-semibold text-yellow-900 mb-3 flex items-center">
                <span className="mr-2">💡</span>
                Why We Recommend This Package
              </h3>
              <p className="text-yellow-800 leading-relaxed">
                {recommendation.reasoning}
              </p>
            </div>

            {/* Actions */}
            <div className="flex justify-end space-x-3 pt-4 border-t border-gray-200">
              <button
                onClick={onClose}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Close
              </button>
              <button
                onClick={() => {
                  onClose();
                  // This would typically trigger the selection flow
                  console.log('Selected recommendation:', recommendation);
                }}
                className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
              >
                Select This Package
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
