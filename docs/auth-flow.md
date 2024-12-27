# Authentication and Multi-tenant System Documentation

## Overview
This document describes the authentication flow and multi-tenant system, written in a prompt-style format to guide frontend development.

## Tenant Types and Features

### Personal Tenant
```typescript
interface PersonalTenant {
  type: 'personal';
  maxUsers: 1;
  features: {
    personalFeatures: true;
    personalDashboard: true;
    personalStorage: true;
    basicIntegrations: true;
  };
  canUpgrade: false;
  limits: {
    storage: '5GB';
    integrations: 3;
    apiCalls: 1000;
  };
}
```
**When to use:**
- Automatically created for each new user
- For individual workspaces
- Cannot be upgraded to other types
- Limited to 1 user (the owner)

**UI States:**
- Default landing after first login
- Cannot be deleted if it's the only tenant
- Cannot transfer ownership
- Shows personal usage metrics

### Team Tenant
```typescript
interface TeamTenant {
  type: 'team';
  maxUsers: 10;
  features: {
    teamFeatures: true;
    collaborationTools: true;
    teamDashboard: true;
    sharedStorage: true;
    advancedIntegrations: true;
    basicAnalytics: true;
    inviteSystem: true;
  };
  canUpgrade: true;
  upgradePath: ['enterprise'];
  limits: {
    storage: '50GB';
    integrations: 10;
    apiCalls: 10000;
    invitesPerMonth: 20;
  };
}
```
**When to use:**
- Default choice for small teams
- When you need basic collaboration
- Up to 10 team members
- Can be upgraded to Enterprise

**UI States:**
- Show member management
- Display upgrade prompts near limits
- Team settings accessible to admins
- Invitation management interface

### Enterprise Tenant
```typescript
interface EnterpriseTenant {
  type: 'enterprise';
  maxUsers: null; // unlimited
  features: {
    enterpriseFeatures: true;
    sso: true;
    audit: true;
    advancedSecurity: true;
    enterpriseDashboard: true;
    unlimitedStorage: true;
    customIntegrations: true;
    advancedAnalytics: true;
    bulkInviteSystem: true;
    customBranding: true;
    apiAccess: true;
    prioritySupport: true;
  };
  domainVerification: true;
  canUpgrade: false;
  limits: {
    storage: 'unlimited';
    integrations: 'unlimited';
    apiCalls: 100000;
    invitesPerMonth: 'unlimited';
  };
}
```
**When to use:**
- For large organizations
- Need for advanced security features
- SSO requirement
- Domain verification needed
- Unlimited users

**UI States:**
- Domain verification status
- SSO configuration
- Advanced security settings
- Custom branding options
- API key management

## User Journey Scenarios

### 1. First-time User
```typescript
interface FirstTimeUserFlow {
  states: {
    initial: 'LOGGED_OUT',
    afterOAuth: 'PERSONAL_TENANT_CREATION',
    completion: 'PERSONAL_WORKSPACE'
  };
  actions: {
    login: () => void;
    completeProfile?: () => void;
    skipOnboarding?: () => void;
  };
  ui: {
    welcomeScreen: boolean;
    onboarding: boolean;
    profileCompletion: boolean;
  };
}
```

### 2. Tenant Creation Scenarios
```typescript
interface TenantCreationScenarios {
  personal: {
    automatic: true;
    nameGeneration: `${user.name}'s Workspace`;
    slugGeneration: `personal-${uuid}`;
  };
  team: {
    validation: {
      name: string; // 3-50 chars
      slug: string; // 3-30 chars, lowercase, no spaces
      initialMembers?: string[]; // email addresses
    };
    steps: [
      'BASIC_INFO',
      'INVITE_MEMBERS?',
      'PAYMENT_SETUP?',
      'WORKSPACE_CUSTOMIZATION'
    ];
  };
  enterprise: {
    validation: {
      name: string;
      slug: string;
      domain: string;
      contractDetails?: {
        contactPerson: string;
        billingEmail: string;
        phone?: string;
      };
    };
    steps: [
      'BASIC_INFO',
      'DOMAIN_SETUP',
      'CONTRACT_DETAILS',
      'SSO_CONFIGURATION?',
      'BRANDING_SETUP?'
    ];
  };
}
```

### 3. Tenant Switching States
```typescript
interface TenantSwitchingStates {
  loading: boolean;
  currentTenant: Tenant;
  availableTenants: Tenant[];
  recentTenants: Tenant[];
  errors?: {
    accessDenied: boolean;
    tenantSuspended: boolean;
    networkError: boolean;
  };
  ui: {
    showTenantBadge: boolean;
    showQuickSwitch: boolean;
    showTenantSearch: boolean;
    showCreateNew: boolean;
  };
}
```

### 4. Feature Access Control
```typescript
interface FeatureAccessControl {
  checkAccess: (feature: string) => {
    hasAccess: boolean;
    requiresUpgrade: boolean;
    upgradeType?: 'team' | 'enterprise';
    reason?: string;
  };
  
  upgradePrompt: {
    show: boolean;
    feature: string;
    currentPlan: string;
    requiredPlan: string;
    benefits: string[];
    action: 'upgrade' | 'contact_sales';
  };
}
```

## UI Components Guide

### 1. Tenant Switcher
```typescript
interface TenantSwitcherProps {
  currentTenant: Tenant;
  tenants: Tenant[];
  layout: 'dropdown' | 'sidebar' | 'modal';
  features: {
    search: boolean;
    recentlyUsed: boolean;
    createNew: boolean;
    tenantIcons: boolean;
  };
  actions: {
    onSwitch: (tenantId: string) => Promise<void>;
    onCreate: () => void;
    onSearch: (query: string) => void;
  };
}
```

### 2. Tenant Settings Panel
```typescript
interface TenantSettingsProps {
  sections: {
    general: {
      name: boolean;
      slug: boolean;
      icon: boolean;
    };
    members: {
      invite: boolean;
      remove: boolean;
      roles: boolean;
    };
    security: {
      sso: boolean;
      2fa: boolean;
      apiKeys: boolean;
    };
    billing: {
      plan: boolean;
      usage: boolean;
      invoices: boolean;
    };
  };
  permissions: string[];
  onUpdate: (section: string, data: any) => Promise<void>;
}
```

### 3. User Management
```typescript
interface UserManagementProps {
  features: {
    invite: {
      single: boolean;
      bulk: boolean;
      domain: boolean;
    };
    roles: {
      custom: boolean;
      predefined: string[];
    };
    access: {
      features: string[];
      permissions: string[];
    };
  };
  limits: {
    maxUsers: number;
    currentUsers: number;
    invitesLeft: number;
  };
}
```

## Error Scenarios

### 1. Tenant Access Errors
```typescript
type TenantError =
  | 'TENANT_NOT_FOUND'
  | 'MAX_USERS_REACHED'
  | 'INVALID_TENANT_TYPE'
  | 'DOMAIN_NOT_VERIFIED'
  | 'ACCESS_DENIED'
  | 'TENANT_SUSPENDED'
  | 'PAYMENT_REQUIRED'
  | 'UPGRADE_REQUIRED'
  | 'INVITATION_EXPIRED'
  | 'DOMAIN_MISMATCH'
  | 'SSO_REQUIRED';

interface ErrorHandling {
  type: TenantError;
  message: string;
  action?: {
    type: 'redirect' | 'modal' | 'toast' | 'page';
    destination?: string;
    component?: React.ComponentType;
  };
  recovery?: {
    possible: boolean;
    steps?: string[];
  };
}
```

### 2. Feature Access Errors
```typescript
interface FeatureError {
  feature: string;
  currentPlan: string;
  requiredPlan: string;
  difference: {
    price?: string;
    features: string[];
  };
  action: {
    type: 'upgrade' | 'contact';
    url?: string;
  };
}
```

## State Management

### 1. Tenant Context
```typescript
interface TenantContext {
  current: Tenant;
  available: Tenant[];
  permissions: string[];
  features: string[];
  loading: boolean;
  error?: TenantError;
}
```

### 2. User Context
```typescript
interface UserContext {
  user: User;
  preferences: {
    defaultTenant?: string;
    theme: 'light' | 'dark' | 'system';
    notifications: NotificationPreferences;
  };
  sessions: {
    current: string;
    active: number;
  };
}
```

### 3. Feature Flags
```typescript
interface FeatureFlags {
  flags: Record<string, boolean>;
  override?: Record<string, boolean>;
  compute: (feature: string) => {
    enabled: boolean;
    source: 'tenant' | 'override' | 'system';
  };
}
```

## Security Implementation

### 1. Tenant Isolation
```typescript
interface TenantSecurity {
  headers: {
    'X-Tenant-ID': string;
    'X-Tenant-Type': string;
    Authorization: string;
  };
  storage: {
    namespace: string;
    encryption: boolean;
    isolation: boolean;
  };
  validation: {
    routes: string[];
    methods: string[];
    features: string[];
  };
}
```

### 2. Domain Verification
```typescript
interface DomainVerification {
  status: 'unverified' | 'pending' | 'verified' | 'failed';
  methods: {
    dns: {
      type: 'TXT' | 'CNAME';
      record: string;
      value: string;
    };
    file: {
      path: string;
      content: string;
    };
  };
  autoJoin: {
    enabled: boolean;
    domains: string[];
    defaultRole: string;
  };
}
```

### 3. Access Control Matrix
```typescript
interface AccessControl {
  roles: {
    owner: string[];
    admin: string[];
    member: string[];
    guest: string[];
  };
  features: {
    [feature: string]: {
      requiredRole: string;
      requiredPlan: string;
      customCheck?: (user: User, tenant: Tenant) => boolean;
    };
  };
  ui: {
    hideDisabled: boolean;
    showUpgrade: boolean;
    customMessage?: string;
  };
} 