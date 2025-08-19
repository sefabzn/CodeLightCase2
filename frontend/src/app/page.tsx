'use client';

import { useState, useEffect } from 'react';
import { apiClient, ApiError } from '@/lib/api-client';
import type { HealthResponse, CoverageInfo } from '@/types/api';

export default function Home() {
  const [healthStatus, setHealthStatus] = useState<HealthResponse | null>(null);
  const [coverage, setCoverage] = useState<CoverageInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const testAPI = async () => {
    setLoading(true);
    setError(null);
    
    try {
      // Test health endpoint
      const health = await apiClient.health();
      setHealthStatus(health);
      
      // Test coverage endpoint with sample data
      const coverageData = await apiClient.getCoverage('A1001');
      setCoverage(coverageData);
      
    } catch (err) {
      if (err instanceof ApiError) {
        setError(`API Error: ${err.message} (${err.code})`);
      } else {
        setError(`Network Error: ${err instanceof Error ? err.message : 'Unknown error'}`);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    testAPI();
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-red-50 to-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <header className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            🚀 Turkcell Recommendation System
          </h1>
          <p className="text-lg text-gray-600">
            AI-powered telecom package recommendations
          </p>
        </header>

        <div className="grid md:grid-cols-2 gap-8">
          {/* API Health Status */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-800">
              🔗 API Connection Status
            </h2>
            
            {loading && (
              <div className="flex items-center space-x-2">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-red-600"></div>
                <span>Testing API connection...</span>
              </div>
            )}
            
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-md p-4">
                <p className="text-red-800">❌ {error}</p>
                <button 
                  onClick={testAPI}
                  className="mt-2 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                >
                  Retry
                </button>
              </div>
            )}
            
            {healthStatus && !loading && (
              <div className="bg-green-50 border border-green-200 rounded-md p-4">
                <p className="text-green-800 font-medium">✅ API Connected!</p>
                <div className="mt-2 text-sm text-gray-600">
                  <p>Status: {healthStatus.status}</p>
                  <p>Database: {healthStatus.database}</p>
                  <p>Service: {healthStatus.service}</p>
                  <p>Version: {healthStatus.version}</p>
                </div>
              </div>
            )}
          </div>

          {/* Sample Data */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-4 text-gray-800">
              📍 Sample Coverage Data
            </h2>
            
            {coverage && (
              <div className="space-y-3">
                <div>
                  <p className="font-medium">Address: {coverage.address_id}</p>
                  <p className="text-gray-600">{coverage.city}, {coverage.district}</p>
                </div>
                
                <div>
                  <p className="font-medium">Available Technologies:</p>
                  <div className="flex space-x-2 mt-1">
                    {coverage.available_tech.map((tech) => (
                      <span 
                        key={tech}
                        className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm"
                      >
                        {tech.toUpperCase()}
                      </span>
                    ))}
                  </div>
                </div>
                
                <div className="grid grid-cols-3 gap-2 text-sm">
                  <div className={`p-2 rounded ${coverage.fiber ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'}`}>
                    Fiber: {coverage.fiber ? '✅' : '❌'}
                  </div>
                  <div className={`p-2 rounded ${coverage.vdsl ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'}`}>
                    VDSL: {coverage.vdsl ? '✅' : '❌'}
                  </div>
                  <div className={`p-2 rounded ${coverage.fwa ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-500'}`}>
                    FWA: {coverage.fwa ? '✅' : '❌'}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="text-center mt-12">
          <p className="text-gray-600">
            Frontend is connected to backend API running on <code className="bg-gray-200 px-2 py-1 rounded">localhost:8000</code>
          </p>
        </div>
      </div>
    </div>
  );
}
