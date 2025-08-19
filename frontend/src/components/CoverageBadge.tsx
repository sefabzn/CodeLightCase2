'use client';

import React from 'react';
import type { CoverageInfo } from '@/types/api';

interface CoverageBadgeProps {
  coverage: CoverageInfo;
  className?: string;
}

interface TechBadgeProps {
  tech: string;
  available: boolean;
  isPrimary?: boolean;
}

function TechBadge({ tech, available, isPrimary = false }: TechBadgeProps) {
  const techConfig = {
    fiber: {
      label: 'Fiber',
      description: 'High-speed fiber optic internet',
      icon: '🚀',
    },
    vdsl: {
      label: 'VDSL',
      description: 'Very-high-bit-rate digital subscriber line',
      icon: '📡',
    },
    fwa: {
      label: 'FWA',
      description: 'Fixed wireless access',
      icon: '📶',
    },
  };

  const config = techConfig[tech as keyof typeof techConfig] || {
    label: tech.toUpperCase(),
    description: `${tech} technology`,
    icon: '📡',
  };

  const baseClasses = "inline-flex items-center px-3 py-1 rounded-full text-sm font-medium transition-colors";
  
  let badgeClasses = baseClasses;
  if (available) {
    if (isPrimary) {
      badgeClasses += " bg-green-100 text-green-800 border-2 border-green-300";
    } else {
      badgeClasses += " bg-green-100 text-green-800";
    }
  } else {
    badgeClasses += " bg-gray-100 text-gray-500";
  }

  return (
    <div className="group relative">
      <span className={badgeClasses}>
        <span className="mr-1">{config.icon}</span>
        {config.label}
        {available && isPrimary && (
          <span className="ml-1 text-xs">(Recommended)</span>
        )}
      </span>
      
      {/* Tooltip */}
      <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-md opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none z-10 whitespace-nowrap">
        <div className="font-medium">{config.label}</div>
        <div className="text-gray-300">{config.description}</div>
        <div className="text-gray-400">
          {available ? 'Available' : 'Not available'}
        </div>
        {/* Arrow */}
        <div className="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
      </div>
    </div>
  );
}

export function CoverageBadge({ coverage, className = "" }: CoverageBadgeProps) {
  const technologies = [
    { name: 'fiber', available: coverage.fiber },
    { name: 'vdsl', available: coverage.vdsl },
    { name: 'fwa', available: coverage.fwa },
  ];

  // Determine primary (recommended) technology
  const primaryTech = coverage.available_tech[0]; // First in available_tech array is the preferred one

  return (
    <div className={`space-y-3 ${className}`}>
      {/* Technology badges */}
      <div className="flex flex-wrap gap-2">
        {technologies.map((tech) => (
          <TechBadge
            key={tech.name}
            tech={tech.name}
            available={tech.available}
            isPrimary={tech.available && tech.name === primaryTech}
          />
        ))}
      </div>

      {/* Available technologies summary */}
      {coverage.available_tech.length > 0 && (
        <div className="text-xs text-gray-600">
          <span className="font-medium">Technology priority:</span>{' '}
          {coverage.available_tech.map((tech, index) => (
            <span key={tech}>
              {tech.toUpperCase()}
              {index < coverage.available_tech.length - 1 ? ' → ' : ''}
            </span>
          ))}
        </div>
      )}

      {/* No coverage message */}
      {coverage.available_tech.length === 0 && (
        <div className="text-xs text-red-600 font-medium">
          No internet technologies available at this address
        </div>
      )}

      {/* Speed expectations */}
      {coverage.available_tech.length > 0 && (
        <div className="text-xs text-gray-500 bg-gray-50 p-2 rounded border">
          <div className="font-medium mb-1">Expected speeds:</div>
          <div className="space-y-1">
            {coverage.fiber && (
              <div>• Fiber: Up to 1000 Mbps download</div>
            )}
            {coverage.vdsl && (
              <div>• VDSL: Up to 100 Mbps download</div>
            )}
            {coverage.fwa && (
              <div>• FWA: Up to 50 Mbps download</div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
