export type AssetStatus = 'AVAILABLE' | 'ASSIGNED' | 'UNDER_REPAIR' | 'RETIRED';

export interface Asset {
  id: number;
  asset_code: string;
  asset_name: string;
  asset_type: string;
  serial_number: string;
  purchase_date: string;
  status: AssetStatus;
  assigned_to?: string | null;
  location: string;
  created_at?: string;
  updated_at?: string;
}

export interface AssetHistory {
  id: number;
  asset_id: number;
  action: string;
  performed_by: string;
  details: string;
  created_at: string;
}

export interface CreateAssetRequest {
  asset_code: string;
  asset_name: string;
  asset_type: string;
  serial_number: string;
  purchase_date: string;
  status?: AssetStatus;
  assigned_to?: string | null;
  location: string;
}

export interface UpdateAssetRequest {
  asset_name: string;
  asset_type: string;
  purchase_date: string;
  location: string;
}

export interface AssignAssetRequest {
  assigned_to: string;
  location?: string;
  notes?: string;
  performed_by?: string;
}

export interface ReturnAssetRequest {
  location?: string;
  notes?: string;
  performed_by?: string;
}

export interface StatusChangeRequest {
  location?: string;
  notes?: string;
  performed_by?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface DashboardStats {
  total: number;
  available: number;
  assigned: number;
  under_repair: number;
  retired: number;
}

export interface AssetFilterParams {
  search?: string;
  status?: string;
  type?: string;
  page?: number;
  limit?: number;
}
