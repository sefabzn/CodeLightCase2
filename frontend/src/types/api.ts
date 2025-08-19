// API Types for Turkcell Recommendation System

export interface RecommendationRequest {
  user_id: number;
  address_id: string;
  household: HouseholdLineDTO[];
  prefer_tech?: string[];
}

export interface HouseholdLineDTO {
  line_id: string;
  expected_gb: number;
  expected_min: number;
  tv_hd_hours: number;
}

export interface RecommendationResponse {
  top3: RecommendationCandidateDTO[];
}

export interface RecommendationCandidateDTO {
  combo_label: string;
  items: RecommendationItemsDTO;
  monthly_total: number;
  savings: number;
  reasoning: string;
  discounts: RecommendationDiscountsDTO;
}

export interface RecommendationItemsDTO {
  mobile: MobilePlanAssignmentDTO[];
  home?: HomePlanDTO | null;
  tv?: TVPlanDTO | null;
}

export interface MobilePlanAssignmentDTO {
  line_id: string;
  plan: MobilePlanDTO;
  line_cost: number;
  overage_gb: number;
  overage_min: number;
}

export interface RecommendationDiscountsDTO {
  line_discount: number;
  bundle_discount: number;
  total_discount: number;
}

export interface MobilePlanDTO {
  plan_id: number;
  plan_name: string;
  quota_gb: number;
  quota_min: number;
  monthly_price: number;
  overage_gb: number;
  overage_min: number;
}

export interface HomePlanDTO {
  home_id: number;
  name: string;
  tech: string;
  down_mbps: number;
  monthly_price: number;
  install_fee: number;
}

export interface TVPlanDTO {
  tv_id: number;
  name: string;
  hd_hours_included: number;
  monthly_price: number;
}

export interface CoverageInfo {
  address_id: string;
  city: string;
  district: string;
  fiber: boolean;
  vdsl: boolean;
  fwa: boolean;
  available_tech: string[];
}

export interface InstallSlot {
  slot_id: string;
  address_id: string;
  slot_start: string;
  slot_end: string;
  tech: string;
  available: boolean;
}

export interface InstallSlotsResponse {
  address_id: string;
  tech: string;
  slots: InstallSlot[];
}

export interface CheckoutRequest {
  user_id: number;
  selected_combo: RecommendationCandidateDTO;
  slot_id: string;
  address_id: string;
}

export interface CheckoutResponse {
  status: string;
  order_id: string;
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: string[];
  };
}

export interface HealthResponse {
  status: string;
  database: string;
  service: string;
  version: string;
}
