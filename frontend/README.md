# IT Asset Management - Angular Frontend

Modern enterprise Angular web application for tracking IT hardware assets, device assignments, service repairs, and audit trails.

## Features
- **Dashboard**: Real-time asset metrics (Total, Available, Assigned, Under Repair, Retired) with progress bars.
- **Inventory List**: Search, filter by status and device category, pagination, and quick action buttons.
- **Asset Actions**:
  - Modal to assign available hardware to staff.
  - Modal to return assigned equipment into storage.
  - Modal to initiate warranty/vendor repairs.
  - Modal to decommission and retire end-of-life hardware.
- **Audit Trail**: Visual chronological timeline tracing every lifecycle event.
- **Asset Registration / Editing**: Complete input validation and duplicate prevention.

## Getting Started

### 1. Install Dependencies
```bash
npm install
```

### 2. Run the Development Server
```bash
npm start
```
Navigate to `http://localhost:4200` in your web browser. The application will automatically reload if you change any of the source files.

### 3. Build for Production
```bash
npm run build
```
The build artifacts will be stored in the `dist/` directory.
