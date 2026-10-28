import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { AssetService } from '../../services/asset.service';
import { AssetStatus, CreateAssetRequest, UpdateAssetRequest } from '../../models/asset.model';

@Component({
  selector: 'app-asset-form',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './asset-form.component.html',
  styleUrls: ['./asset-form.component.css']
})
export class AssetFormComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private assetService = inject(AssetService);

  isEditMode = false;
  assetId?: number;
  isLoading = false;
  isSubmitting = false;
  errorMessage = '';

  // Form Fields
  assetCode = '';
  assetName = '';
  assetType = 'LAPTOP';
  serialNumber = '';
  purchaseDate = new Date().toISOString().substring(0, 10);
  status: AssetStatus = 'AVAILABLE';
  assignedTo = '';
  location = '';

  assetTypes = ['LAPTOP', 'DESKTOP', 'MONITOR', 'SERVER', 'PHONE', 'NETWORK', 'PERIPHERAL', 'OTHER'];
  statusOptions: AssetStatus[] = ['AVAILABLE', 'ASSIGNED', 'UNDER_REPAIR', 'RETIRED'];

  ngOnInit(): void {
    const idParam = this.route.snapshot.paramMap.get('id');
    if (idParam) {
      this.isEditMode = true;
      this.assetId = Number(idParam);
      this.loadAssetData(this.assetId);
    }
  }

  loadAssetData(id: number): void {
    this.isLoading = true;
    this.errorMessage = '';
    this.assetService.getAsset(id).subscribe({
      next: (asset) => {
        this.assetCode = asset.asset_code;
        this.assetName = asset.asset_name;
        this.assetType = asset.asset_type;
        this.serialNumber = asset.serial_number;
        this.purchaseDate = asset.purchase_date;
        this.status = asset.status;
        this.assignedTo = asset.assigned_to || '';
        this.location = asset.location;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to load asset', err);
        this.errorMessage = 'Could not retrieve asset details.';
        this.isLoading = false;
      }
    });
  }

  onSubmit(): void {
    this.errorMessage = '';

    if (!this.assetName.trim() || !this.location.trim() || !this.purchaseDate) {
      this.errorMessage = 'Please complete all required fields.';
      return;
    }

    if (!this.isEditMode) {
      if (!this.assetCode.trim() || !this.serialNumber.trim()) {
        this.errorMessage = 'Asset Code and Serial Number are required.';
        return;
      }

      this.isSubmitting = true;
      const createPayload: CreateAssetRequest = {
        asset_code: this.assetCode.trim(),
        asset_name: this.assetName.trim(),
        asset_type: this.assetType,
        serial_number: this.serialNumber.trim(),
        purchase_date: this.purchaseDate,
        status: this.status,
        assigned_to: this.assignedTo.trim() ? this.assignedTo.trim() : null,
        location: this.location.trim()
      };

      this.assetService.createAsset(createPayload).subscribe({
        next: () => {
          this.isSubmitting = false;
          this.router.navigate(['/assets']);
        },
        error: (err) => {
          this.isSubmitting = false;
          this.errorMessage = err.error?.error || 'Failed to create asset. Check if asset code or serial number is unique.';
        }
      });
    } else {
      if (!this.assetId) return;

      this.isSubmitting = true;
      const updatePayload: UpdateAssetRequest = {
        asset_name: this.assetName.trim(),
        asset_type: this.assetType,
        purchase_date: this.purchaseDate,
        location: this.location.trim()
      };

      this.assetService.updateAsset(this.assetId, updatePayload).subscribe({
        next: () => {
          this.isSubmitting = false;
          this.router.navigate(['/assets']);
        },
        error: (err) => {
          this.isSubmitting = false;
          this.errorMessage = err.error?.error || 'Failed to update asset.';
        }
      });
    }
  }
}
