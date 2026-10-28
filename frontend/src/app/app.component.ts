import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <div class="app-layout">
      <!-- Enterprise Top Navigation Bar -->
      <header class="navbar">
        <div class="navbar-container">
          <div class="brand">
            <span class="logo-icon">💻</span>
            <div class="brand-info">
              <span class="brand-title">IT Asset Management</span>
              <span class="brand-badge">Enterprise Edition</span>
            </div>
          </div>

          <nav class="nav-links">
            <a routerLink="/dashboard" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">
              Dashboard
            </a>
            <a routerLink="/assets" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">
              Asset Inventory
            </a>
            <a routerLink="/assets/new" routerLinkActive="active">
              + New Asset
            </a>
            <a routerLink="/history" routerLinkActive="active">
              Audit Logs
            </a>
          </nav>

          <div class="header-right">
            <div class="status-indicator-box">
              <span class="live-dot"></span>
              <span class="status-text">API Online</span>
            </div>
          </div>
        </div>
      </header>

      <!-- Main Application Body -->
      <main class="main-content">
        <div class="container">
          <router-outlet></router-outlet>
        </div>
      </main>

      <!-- Footer -->
      <footer class="footer">
        <div class="container footer-content">
          <p>&copy; 2026 IT Asset Management Platform. Full-Stack Go & Angular Architecture.</p>
        </div>
      </footer>
    </div>
  `,
  styles: [`
    .app-layout {
      display: flex;
      flex-direction: column;
      min-height: 100vh;
      background-color: #f8fafc;
    }
    .navbar {
      background-color: #0f172a;
      border-bottom: 1px solid #1e293b;
      position: sticky;
      top: 0;
      z-index: 50;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    }
    .navbar-container {
      max-width: 1280px;
      margin: 0 auto;
      padding: 0 1.5rem;
      height: 68px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }
    .logo-icon {
      font-size: 1.5rem;
      background: #1e293b;
      padding: 0.4rem;
      border-radius: 8px;
      border: 1px solid #334155;
    }
    .brand-info {
      display: flex;
      flex-direction: column;
    }
    .brand-title {
      color: #ffffff;
      font-size: 1.1rem;
      font-weight: 700;
      letter-spacing: -0.01em;
    }
    .brand-badge {
      color: #94a3b8;
      font-size: 0.7rem;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    .nav-links {
      display: flex;
      gap: 0.5rem;
    }
    .nav-links a {
      color: #cbd5e1;
      font-size: 0.875rem;
      font-weight: 500;
      padding: 0.5rem 0.85rem;
      border-radius: 6px;
      transition: all 0.15s ease;
      text-decoration: none;
    }
    .nav-links a:hover {
      color: #ffffff;
      background-color: #1e293b;
    }
    .nav-links a.active {
      color: #ffffff;
      background-color: #2563eb;
    }
    .header-right {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    .status-indicator-box {
      display: flex;
      align-items: center;
      gap: 0.4rem;
      background: #1e293b;
      padding: 0.3rem 0.65rem;
      border-radius: 9999px;
      border: 1px solid #334155;
    }
    .live-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background-color: #10b981;
      box-shadow: 0 0 6px #10b981;
    }
    .status-text {
      color: #e2e8f0;
      font-size: 0.75rem;
      font-weight: 500;
    }
    .main-content {
      flex: 1;
      padding: 2rem 0 3rem;
    }
    .footer {
      background-color: #ffffff;
      border-top: 1px solid #e2e8f0;
      padding: 1.25rem 0;
      color: #64748b;
      font-size: 0.825rem;
      text-align: center;
    }
  `]
})
export class AppComponent {
  title = 'IT Asset Management';
}
