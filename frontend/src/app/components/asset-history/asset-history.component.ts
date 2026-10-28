import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { AssetService } from '../../services/asset.service';
import { Asset, AssetHistory } from '../../models/asset.model';

@Component({
  selector: 'app-asset-history',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './asset-history.component.html',
  styleUrls: ['./asset-history.component.css']
})
export class AssetHistoryComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private assetService = inject(AssetService);

  assetId: number | null = null;
  asset: Asset | null = null;
  historyList: AssetHistory[] = [];
  assetsList: Asset[] = [];

  isLoading = false;
  errorMessage = '';

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      const idParam = params.get('id');
      if (idParam) {
        this.assetId = Number(idParam);
        this.loadAssetAndHistory(this.assetId);
      } else {
        // Load asset inventory so user can select which asset to audit
        this.loadAvailableAssets();
      }
    });
  }

  loadAvailableAssets(): void {
    this.isLoading = true;
    this.assetService.getAssets({ limit: 100 }).subscribe({
      next: (res) => {
        this.assetsList = res.data || [];
        this.isLoading = false;
        if (this.assetsList.length > 0 && !this.assetId) {
          this.onSelectAsset(this.assetsList[0].id);
        }
      },
      error: (err) => {
        console.error('Failed to load asset selector', err);
        this.errorMessage = 'Unable to load asset list for history inspection.';
        this.isLoading = false;
      }
    });
  }

  onSelectAsset(id: number): void {
    this.assetId = id;
    this.loadAssetAndHistory(id);
  }

  loadAssetAndHistory(id: number): void {
    this.isLoading = true;
    this.errorMessage = '';

    this.assetService.getAsset(id).subscribe({
      next: (asset) => {
        this.asset = asset;
        this.loadHistory(id);
      },
      error: (err) => {
        console.error('Failed to fetch asset details', err);
        this.errorMessage = 'Asset not found or failed to load.';
        this.isLoading = false;
      }
    });
  }

  loadHistory(id: number): void {
    this.assetService.getAssetHistory(id).subscribe({
      next: (history) => {
        this.historyList = history || [];
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to fetch history logs', err);
        this.errorMessage = 'Failed to load audit history records.';
        this.isLoading = false;
      }
    });
  }

  getActionBadgeClass(action: string): string {
    switch (action.toUpperCase()) {
      case 'CREATED': return 'badge-created';
      case 'ASSIGNED': return 'badge-assigned';
      case 'RETURNED': return 'badge-returned';
      case 'REPAIR': return 'badge-repair';
      case 'RETIRED': return 'badge-retired';
      case 'UPDATED': return 'badge-updated';
      default: return 'badge-secondary';
    }
  }

  getActionIcon(action: string): string {
    switch (action.toUpperCase()) {
      case 'CREATED': return '✨';
      case 'ASSIGNED': return '👤';
      case 'RETURNED': return '↩️';
      case 'REPAIR': return '🔧';
      case 'RETIRED': return '🗑️';
      case 'UPDATED': return '✏️';
      default: return '📌';
    }
  }
}
