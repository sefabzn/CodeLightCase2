'use client';

import React, { useState } from 'react';
import type { RecommendationCandidateDTO } from '@/types/api';

interface RecommendationCardProps {
  recommendation: RecommendationCandidateDTO;
  rank: number;
  onSelect: (recommendation: RecommendationCandidateDTO) => void;
  onShowDetails: (recommendation: RecommendationCandidateDTO) => void;
}

export function RecommendationCard({ 
  recommendation, 
  rank, 
  onSelect, 
  onShowDetails 
}: RecommendationCardProps) {
  const [isHovered, setIsHovered] = useState(false);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('tr-TR', {
      style: 'currency',
      currency: 'TRY',
      minimumFractionDigits: 2,
    }).format(price);
  };

  const getBadgeColor = (rank: number) => {
    switch (rank) {
      case 1:
        return 'bg-green-100 text-green-800 border-green-200';
      case 2:
        return 'bg-blue-100 text-blue-800 border-blue-200';
      case 3:
        return 'bg-purple-100 text-purple-800 border-purple-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getRankLabel = (rank: number) => {
    switch (rank) {
      case 1:
        return 'Best Value';
      case 2:
        return 'Popular Choice';
      case 3:
        return 'Premium Option';
      default:
        return `Option ${rank}`;
    }
  };

  return (
    <div
      className={`relative bg-white rounded-lg border-2 transition-all duration-200 ${
        isHovered ? 'border-blue-300 shadow-lg' : 'border-gray-200 shadow-sm'
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Rank Badge */}
      <div className="absolute -top-3 left-4 z-10">
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${getBadgeColor(rank)}`}>
          #{rank} {getRankLabel(rank)}
        </span>
      </div>

      <div className="p-6 pt-8">
        {/* Header */}
        <div className="flex justify-between items-start mb-4">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">
              {recommendation.combo_label}
            </h3>
            <div className="flex items-center space-x-2">
              <span className="text-3xl font-bold text-blue-600">
                {formatPrice(recommendation.monthly_total)}
              </span>
              <span className="text-gray-600">/month</span>
            </div>
          </div>
          
          {/* Savings */}
          {recommendation.savings > 0 && (
            <div className="text-right">
              <div className="bg-green-50 text-green-700 px-2 py-1 rounded-md text-sm font-medium">
                Save {formatPrice(recommendation.savings)}
              </div>
              <div className="text-xs text-gray-500 mt-1">vs individual plans</div>
            </div>
          )}
        </div>

        {/* Services Included */}
        <div className="space-y-3 mb-6">
          {/* Mobile Lines */}
          {recommendation.items.mobile.length > 0 && (
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                <span className="text-blue-600 text-sm font-medium">📱</span>
              </div>
              <div className="flex-1">
                <div className="font-medium text-gray-900">
                  Mobile Plans ({recommendation.items.mobile.length} line{recommendation.items.mobile.length > 1 ? 's' : ''})
                </div>
                <div className="text-sm text-gray-600">
                  {recommendation.items.mobile.map((mobile, index) => (
                    <span key={mobile.line_id}>
                      {mobile.plan.plan_name}
                      {index < recommendation.items.mobile.length - 1 ? ', ' : ''}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* Home Internet */}
          {recommendation.items.home && (
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                <span className="text-green-600 text-sm font-medium">🏠</span>
              </div>
              <div className="flex-1">
                <div className="font-medium text-gray-900">Home Internet</div>
                <div className="text-sm text-gray-600">
                  {recommendation.items.home.name} - {recommendation.items.home.down_mbps} Mbps ({recommendation.items.home.tech.toUpperCase()})
                </div>
              </div>
            </div>
          )}

          {/* TV */}
          {recommendation.items.tv && (
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                <span className="text-purple-600 text-sm font-medium">📺</span>
              </div>
              <div className="flex-1">
                <div className="font-medium text-gray-900">TV Package</div>
                <div className="text-sm text-gray-600">
                  {recommendation.items.tv.name} - {recommendation.items.tv.hd_hours_included} HD hours
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Discounts Applied */}
        {(recommendation.discounts.line_discount > 0 || recommendation.discounts.bundle_discount > 0) && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-md p-3 mb-6">
            <div className="text-sm text-yellow-800">
              <div className="font-medium mb-1">💰 Discounts Applied:</div>
              <ul className="space-y-1">
                {recommendation.discounts.line_discount > 0 && (
                  <li>• Multi-line discount: {formatPrice(recommendation.discounts.line_discount)}</li>
                )}
                {recommendation.discounts.bundle_discount > 0 && (
                  <li>• Bundle discount: {formatPrice(recommendation.discounts.bundle_discount)}</li>
                )}
              </ul>
              <div className="font-medium mt-2">
                Total savings: {formatPrice(recommendation.discounts.total_discount)}
              </div>
            </div>
          </div>
        )}

        {/* Key Features */}
        <div className="bg-gray-50 rounded-md p-3 mb-6">
          <div className="text-sm">
            <div className="font-medium text-gray-900 mb-2">✨ Why this package:</div>
            <div className="text-gray-700 leading-relaxed">
              {recommendation.reasoning}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3">
          <button
            onClick={() => onShowDetails(recommendation)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors text-sm font-medium"
          >
            View Details
          </button>
          <button
            onClick={() => onSelect(recommendation)}
            className={`flex-1 px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              rank === 1
                ? 'bg-blue-600 text-white hover:bg-blue-700'
                : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
            }`}
          >
            {rank === 1 ? '🏆 Select Best Option' : 'Select This Package'}
          </button>
        </div>
      </div>

      {/* Recommended ribbon for rank 1 */}
      {rank === 1 && (
        <div className="absolute top-4 right-4">
          <div className="bg-gradient-to-r from-green-400 to-green-600 text-white px-2 py-1 rounded-full text-xs font-bold transform rotate-12">
            RECOMMENDED
          </div>
        </div>
      )}
    </div>
  );
}
