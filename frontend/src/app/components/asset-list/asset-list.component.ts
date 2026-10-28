import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink, ActivatedRoute } from '@angular/router';
import { AssetService } from '../../services/asset.service';
import { Asset, AssetHistory, AssetStatus } from '../../models/asset.model';

@Component({
  selector: 'app-asset-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './asset-list.component.html',
  styleUrls: ['./asset-list.component.css']
})
export class AssetListComponent implements OnInit {
  private assetService = inject(AssetService);
  private route = inject(ActivatedRoute);

  assets: Asset[] = [];
  totalAssets = 0;
  totalPages = 1;
  currentPage = 1;
  pageSize = 10;

  // Filters
  searchTerm = '';
  selectedStatus = '';
  selectedType = '';

  isLoading = false;
  errorMessage = '';
  successMessage = '';

  // Modal State
  activeModal: 'ASSIGN' | 'RETURN' | 'REPAIR' | 'RETIRE' | 'HISTORY' | null = null;
  selectedAsset: Asset | null = null;
  selectedAssetHistory: AssetHistory[] = [];
  isHistoryLoading = false;

  // Modal Form Inputs
  assigneeName = '';
  actionLocation = '';
  actionNotes = '';

  assetTypes = ['LAPTOP', 'DESKTOP', 'MONITOR', 'SERVER', 'PHONE', 'NETWORK', 'PERIPHERAL', 'OTHER'];
  statusOptions: AssetStatus[] = ['AVAILABLE', 'ASSIGNED', 'UNDER_REPAIR', 'RETIRED'];

  ngOnInit(): void {
    // Check if query params pre-select filters (e.g. from dashboard cards)
    this.route.queryParams.subscribe(params => {
      if (params['status']) {
        this.selectedStatus = params['status'];
      }
      if (params['search']) {
        this.searchTerm = params['search'];
      }
      this.loadAssets();
    });
  }

  loadAssets(page: number = this.currentPage): void {
    this.isLoading = true;
    this.errorMessage = '';
    this.currentPage = page;

    this.assetService.getAssets({
      search: this.searchTerm,
      status: this.selectedStatus,
      type: this.selectedType,
      page: this.currentPage,
      limit: this.pageSize
    }).subscribe({
      next: (res) => {
        this.assets = res.data || [];
        this.totalAssets = res.total;
        this.totalPages = res.total_pages || 1;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to load assets', err);
        this.errorMessage = 'Failed to load assets. Verify backend server is running.';
        this.isLoading = false;
      }
    });
  }

  onFilterChange(): void {
    this.currentPage = 1;
    this.loadAssets();
  }

  resetFilters(): void {
    this.searchTerm = '';
    this.selectedStatus = '';
    this.selectedType = '';
    this.currentPage = 1;
    this.loadAssets();
  }

  deleteAsset(asset: Asset): void {
    if (confirm(`Are you sure you want to permanently delete asset ${asset.asset_code} (${asset.asset_name})?`)) {
      this.assetService.deleteAsset(asset.id).subscribe({
        next: () => {
          this.setSuccess(`Asset ${asset.asset_code} deleted successfully.`);
          this.loadAssets();
        },
        error: (err) => {
          this.setError(err.error?.error || 'Failed to delete asset.');
        }
      });
    }
  }

  // --- Modal Handlers ---

  openAssignModal(asset: Asset): void {
    this.selectedAsset = asset;
    this.assigneeName = '';
    this.actionLocation = asset.location || '';
    this.actionNotes = '';
    this.activeModal = 'ASSIGN';
  }

  openReturnModal(asset: Asset): void {
    this.selectedAsset = asset;
    this.actionLocation = 'IT Storage Room';
    this.actionNotes = '';
    this.activeModal = 'RETURN';
  }

  openRepairModal(asset: Asset): void {
    this.selectedAsset = asset;
    this.actionLocation = 'Service Center / TechFix';
    this.actionNotes = '';
    this.activeModal = 'REPAIR';
  }

  openRetireModal(asset: Asset): void {
    this.selectedAsset = asset;
    this.actionLocation = 'E-Waste Decommission';
    this.actionNotes = '';
    this.activeModal = 'RETIRE';
  }

  openHistoryModal(asset: Asset): void {
    this.selectedAsset = asset;
    this.activeModal = 'HISTORY';
    this.isHistoryLoading = true;
    this.selectedAssetHistory = [];

    this.assetService.getAssetHistory(asset.id).subscribe({
      next: (history) => {
        this.selectedAssetHistory = history || [];
        this.isHistoryLoading = false;
      },
      error: (err) => {
        console.error('Failed to load history', err);
        this.setError('Could not load audit history for asset.');
        this.isHistoryLoading = false;
      }
    });
  }

  closeModal(): void {
    this.activeModal = null;
    this.selectedAsset = null;
    this.assigneeName = '';
    this.actionLocation = '';
    this.actionNotes = '';
  }

  confirmAssign(): void {
    if (!this.selectedAsset || !this.assigneeName.trim()) return;

    this.assetService.assignAsset(this.selectedAsset.id, {
      assigned_to: this.assigneeName.trim(),
      location: this.actionLocation.trim(),
      notes: this.actionNotes.trim(),
      performed_by: 'Admin'
    }).subscribe({
      next: (updated) => {
        this.setSuccess(`Asset ${updated.asset_code} assigned to ${this.assigneeName}.`);
        this.closeModal();
        this.loadAssets();
      },
      error: (err) => {
        this.setError(err.error?.error || 'Failed to assign asset.');
      }
    });
  }

  confirmReturn(): void {
    if (!this.selectedAsset) return;

    this.assetService.returnAsset(this.selectedAsset.id, {
      location: this.actionLocation.trim(),
      notes: this.actionNotes.trim(),
      performed_by: 'Admin'
    }).subscribe({
      next: (updated) => {
        this.setSuccess(`Asset ${updated.asset_code} returned to available inventory.`);
        this.closeModal();
        this.loadAssets();
      },
      error: (err) => {
        this.setError(err.error?.error || 'Failed to return asset.');
      }
    });
  }

  confirmRepair(): void {
    if (!this.selectedAsset) return;

    this.assetService.repairAsset(this.selectedAsset.id, {
      location: this.actionLocation.trim(),
      notes: this.actionNotes.trim(),
      performed_by: 'Admin'
    }).subscribe({
      next: (updated) => {
        this.setSuccess(`Asset ${updated.asset_code} marked as UNDER_REPAIR.`);
        this.closeModal();
        this.loadAssets();
      },
      error: (err) => {
        this.setError(err.error?.error || 'Failed to mark asset for repair.');
      }
    });
  }

  confirmRetire(): void {
    if (!this.selectedAsset) return;

    this.assetService.retireAsset(this.selectedAsset.id, {
      location: this.actionLocation.trim(),
      notes: this.actionNotes.trim(),
      performed_by: 'Admin'
    }).subscribe({
      next: (updated) => {
        this.setSuccess(`Asset ${updated.asset_code} marked as RETIRED.`);
        this.closeModal();
        this.loadAssets();
      },
      error: (err) => {
        this.setError(err.error?.error || 'Failed to retire asset.');
      }
    });
  }

  getStatusBadgeClass(status: string): string {
    switch (status) {
      case 'AVAILABLE': return 'badge-available';
      case 'ASSIGNED': return 'badge-assigned';
      case 'UNDER_REPAIR': return 'badge-repair';
      case 'RETIRED': return 'badge-retired';
      default: return '';
    }
  }

  private setSuccess(msg: string): void {
    this.successMessage = msg;
    setTimeout(() => this.successMessage = '', 5000);
  }

  private setError(msg: string): void {
    this.errorMessage = msg;
    setTimeout(() => this.errorMessage = '', 6000);
  }
}
