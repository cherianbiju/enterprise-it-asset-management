import { Routes } from '@angular/router';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { AssetListComponent } from './components/asset-list/asset-list.component';
import { AssetFormComponent } from './components/asset-form/asset-form.component';
import { AssetHistoryComponent } from './components/asset-history/asset-history.component';

export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  { path: 'dashboard', component: DashboardComponent, title: 'Dashboard - IT Asset Management' },
  { path: 'assets', component: AssetListComponent, title: 'Assets Inventory - IT Asset Management' },
  { path: 'assets/new', component: AssetFormComponent, title: 'Register Asset - IT Asset Management' },
  { path: 'assets/edit/:id', component: AssetFormComponent, title: 'Edit Asset - IT Asset Management' },
  { path: 'assets/:id/history', component: AssetHistoryComponent, title: 'Asset History - IT Asset Management' },
  { path: 'history', component: AssetHistoryComponent, title: 'Audit Trail - IT Asset Management' },
  { path: '**', redirectTo: 'dashboard' }
];
