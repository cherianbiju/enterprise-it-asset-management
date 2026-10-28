import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import {
  Asset,
  AssetHistory,
  CreateAssetRequest,
  UpdateAssetRequest,
  AssignAssetRequest,
  ReturnAssetRequest,
  StatusChangeRequest,
  PaginatedResponse,
  DashboardStats,
  AssetFilterParams
} from '../models/asset.model';

@Injectable({
  providedIn: 'root'
})
export class AssetService {
  private http = inject(HttpClient);
  private apiUrl = environment.apiUrl;

  getAssets(filter?: AssetFilterParams): Observable<PaginatedResponse<Asset>> {
    let params = new HttpParams();
    if (filter) {
      if (filter.search && filter.search.trim()) {
        params = params.set('search', filter.search.trim());
      }
      if (filter.status && filter.status.trim()) {
        params = params.set('status', filter.status.trim());
      }
      if (filter.type && filter.type.trim()) {
        params = params.set('type', filter.type.trim());
      }
      if (filter.page) {
        params = params.set('page', filter.page.toString());
      }
      if (filter.limit) {
        params = params.set('limit', filter.limit.toString());
      }
    }
    return this.http.get<PaginatedResponse<Asset>>(`${this.apiUrl}/assets`, { params });
  }

  getStats(): Observable<DashboardStats> {
    return this.http.get<DashboardStats>(`${this.apiUrl}/assets/stats`);
  }

  getAsset(id: number): Observable<Asset> {
    return this.http.get<Asset>(`${this.apiUrl}/assets/${id}`);
  }

  createAsset(payload: CreateAssetRequest): Observable<Asset> {
    return this.http.post<Asset>(`${this.apiUrl}/assets`, payload);
  }

  updateAsset(id: number, payload: UpdateAssetRequest): Observable<Asset> {
    return this.http.put<Asset>(`${this.apiUrl}/assets/${id}`, payload);
  }

  deleteAsset(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/assets/${id}`);
  }

  assignAsset(id: number, payload: AssignAssetRequest): Observable<Asset> {
    return this.http.post<Asset>(`${this.apiUrl}/assets/${id}/assign`, payload);
  }

  returnAsset(id: number, payload: ReturnAssetRequest): Observable<Asset> {
    return this.http.post<Asset>(`${this.apiUrl}/assets/${id}/return`, payload);
  }

  repairAsset(id: number, payload: StatusChangeRequest): Observable<Asset> {
    return this.http.post<Asset>(`${this.apiUrl}/assets/${id}/repair`, payload);
  }

  retireAsset(id: number, payload: StatusChangeRequest): Observable<Asset> {
    return this.http.post<Asset>(`${this.apiUrl}/assets/${id}/retire`, payload);
  }

  getAssetHistory(id: number): Observable<AssetHistory[]> {
    return this.http.get<AssetHistory[]>(`${this.apiUrl}/assets/${id}/history`);
  }
}
