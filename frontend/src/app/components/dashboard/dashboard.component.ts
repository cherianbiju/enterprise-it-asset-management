import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AssetService } from '../../services/asset.service';
import { DashboardStats } from '../../models/asset.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  private assetService = inject(AssetService);

  stats: DashboardStats = {
    total: 0,
    available: 0,
    assigned: 0,
    under_repair: 0,
    retired: 0
  };

  isLoading = true;
  errorMessage = '';

  ngOnInit(): void {
    this.loadStats();
  }

  loadStats(): void {
    this.isLoading = true;
    this.errorMessage = '';
    this.assetService.getStats().subscribe({
      next: (data) => {
        this.stats = data;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Error fetching dashboard statistics:', err);
        this.errorMessage = 'Unable to connect to backend service. Please check that the Go server is running.';
        this.isLoading = false;
      }
    });
  }

  getPercentage(value: number): number {
    if (!this.stats.total || this.stats.total === 0) return 0;
    return Math.round((value / this.stats.total) * 100);
  }
}
