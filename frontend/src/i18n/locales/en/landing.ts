export default {
  batchImageGuide: {
    title: 'Batch Image Generation',
    description: 'Submit multiple prompts in one job and download the generated images when complete'
  },
  // Home Page
  home: {
    viewOnGithub: 'View on GitHub',
    viewDocs: 'View Documentation',
    docs: 'Docs',
    switchToLight: 'Switch to Light Mode',
    switchToDark: 'Switch to Dark Mode',
    dashboard: 'Dashboard',
    login: 'Login',
    getStarted: 'Start Relay Access',
    goToDashboard: 'Go to Dashboard',
    // User-focused value proposition
    heroEyebrow: 'Codex Compute Relay',
    heroTitle: 'Professional Codex Compute Relay',
    heroTitlePrefix: 'Professional Codex',
    heroTitleSuffix: 'Compute Relay',
    heroSubtitle: 'Stable relay access for Codex CLI, Codex Desktop, and development teams',
    heroDescription:
      'Use managed Codex compute relay immediately, or let us handle private deployment, domain setup, account pool configuration, usage controls, and ongoing operations.',
    metrics: {
      relay: {
        value: 'Relay',
        label: 'Codex access'
      },
      setup: {
        value: 'Setup',
        label: 'Managed delivery'
      },
      billing: {
        value: 'Control',
        label: 'Usage limits'
      }
    },
    download: {
      badge: 'Desktop builds ready',
      title: 'Download Sub2API Desktop',
      description:
        'Both macOS and Windows builds are available. Mac users can download the signed and Apple-notarized ZIP; Windows users can use the Microsoft Installer link. Sign in, choose your relay key, and write the local Codex config without copying Base URL or API Key by hand.',
      macButton: 'Download for Mac',
      windowsButton: 'Download for Windows',
      versionLabel: 'Available builds',
      version: 'macOS Apple Silicon / Windows',
      integrity: 'Mac: Notarized ZIP · Windows: Microsoft Installer',
      steps: {
        install: {
          title: '1. Download and install',
          description: 'Choose the build for your system; unzip and move the macOS app to Applications, or follow the Windows installer prompts.'
        },
        login: {
          title: '2. Sign in and choose a key',
          description: 'Open the dashboard and select an available Codex relay group and API key.'
        },
        write: {
          title: '3. Write Codex config',
          description: 'The desktop app backs up and updates local ~/.codex config, then Codex is ready in terminal.'
        }
      }
    },
    desktopGuide: {
      badge: 'Desktop guide',
      title: 'Finish local Codex setup',
      description: 'You are already in the desktop app, so this page should guide setup instead of offering another app download. Check whether Codex is installed, then install Codex or import the relay config.',
      checkButton: 'Check Codex',
      downloadCodex: 'Download Codex',
      importConfig: 'Import config',
      recheck: 'Check again',
      codexStatusLabel: 'Local Codex status',
      statusChecking: 'Checking...',
      statusUnknown: 'Not checked yet',
      statusInstalled: 'Installed',
      statusMissing: 'Codex not detected',
      statusError: 'Check failed',
      statusHint: 'Click check to read the local codex command status.',
      missingHint: 'The codex command was not found in PATH. Install Codex first.',
      checkFailed: 'Could not check the Codex install status.',
      downloadFailed: 'Could not open the Codex download link.',
      steps: {
        check: {
          title: '1. Check local Codex',
          description: 'Look for the codex command in the system PATH and show its version or install path.'
        },
        download: {
          title: '2. Install Codex',
          description: 'If Codex is missing, open the official install entry, then return here and check again.'
        },
        import: {
          title: '3. Import relay config',
          description: 'Sign in, choose a Codex relay key, and write the local ~/.codex config in one click.'
        }
      }
    },
    tags: {
      subscriptionToApi: 'Codex Compute Relay',
      stickySession: 'Managed Setup',
      realtimeBilling: 'Stable Session Routing'
    },
    // Pain points section
    painPoints: {
      title: 'Sound Familiar?',
      items: {
        expensive: {
          title: 'High Subscription Costs',
          desc: 'Paying for multiple AI subscriptions that add up every month'
        },
        complex: {
          title: 'Account Chaos',
          desc: 'Managing scattered accounts and API keys across different platforms'
        },
        unstable: {
          title: 'Service Interruptions',
          desc: 'Single accounts hitting rate limits and disrupting your workflow'
        },
        noControl: {
          title: 'No Usage Control',
          desc: "Can't track where your money goes or limit team member usage"
        }
      }
    },
    // Solutions section
    solutions: {
      title: 'We Solve These Problems',
      subtitle: 'Three simple steps to stress-free AI access'
    },
    features: {
      unifiedGateway: 'Professional Codex Relay',
      unifiedGatewayDesc: 'API access built for Codex CLI and Desktop, with one key and one entry point ready for development workflows.',
      multiAccount: 'Stable Pool Scheduling',
      multiAccountDesc: 'Health checks, load distribution, and session persistence reduce rate-limit interruptions and account switching.',
      balanceQuota: 'Managed Setup Delivery',
      balanceQuotaDesc: 'Private relay setup is available, including deployment, domain wiring, configuration, verification, and ongoing operations.'
    },
    // Comparison section
    comparison: {
      title: 'Why Choose Us?',
      headers: {
        feature: 'Comparison',
        official: 'Official Subscriptions',
        us: 'Our Platform'
      },
      items: {
        pricing: {
          feature: 'Pricing',
          official: 'Fixed monthly fee, pay even if unused',
          us: 'Pay only for what you use'
        },
        models: {
          feature: 'Model Selection',
          official: 'Single provider only',
          us: 'Switch between models freely'
        },
        management: {
          feature: 'Account Management',
          official: 'Manage each service separately',
          us: 'Unified key, one dashboard'
        },
        stability: {
          feature: 'Stability',
          official: 'Single account rate limits',
          us: 'Multi-account pool, auto-failover'
        },
        control: {
          feature: 'Usage Control',
          official: 'Not available',
          us: 'Quotas & detailed analytics'
        }
      }
    },
    providers: {
      title: 'Codex Workflow Access',
      description: 'Unified access for terminals, desktop clients, and private deployments',
      supported: 'Ready',
      soon: 'Soon',
      service: 'Service',
      codexCli: 'Codex CLI',
      codexDesktop: 'Codex Desktop',
      gptCodex: 'GPT Codex',
      claudeCode: 'Claude Code',
      privateDeploy: 'Private Deploy',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: 'More'
    },
    // CTA section
    cta: {
      title: 'Ready to Get Started?',
      description: 'Sign up now and get free trial credits to experience seamless AI access',
      button: 'Sign Up Free'
    },
    footer: {
      allRightsReserved: 'All rights reserved.'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key Usage',
    subtitle: 'Enter your API Key to view real-time spending and usage status',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: 'Query',
    querying: 'Querying...',
    privacyNote: 'Your Key is processed locally in the browser and will not be stored',
    dateRange: 'Date Range:',
    dateRangeToday: 'Today',
    dateRange7d: '7 Days',
    dateRange30d: '30 Days',
    dateRange90d: '90 Days',
    dateRangeCustom: 'Custom',
    apply: 'Apply',
    used: 'Used',
    detailInfo: 'Detail Information',
    tokenStats: 'Token Statistics',
    dailyDetail: 'Daily Detail',
    modelStats: 'Model Usage Statistics',
    // Table headers
    date: 'Date',
    model: 'Model',
    requests: 'Requests',
    inputTokens: 'Input Tokens',
    outputTokens: 'Output Tokens',
    cacheCreationTokens: 'Cache Creation',
    cacheReadTokens: 'Cache Read',
    cacheWriteTokens: 'Cache Write',
    totalTokens: 'Total Tokens',
    cost: 'Cost',
    // Status
    quotaMode: 'Key Quota Mode',
    walletBalance: 'Wallet Balance',
    // Ring card titles
    totalQuota: 'Total Quota',
    limit5h: '5-Hour Limit',
    limitDaily: 'Daily Limit',
    limit7d: '7-Day Limit',
    limitWeekly: 'Weekly Limit',
    limitMonthly: 'Monthly Limit',
    // Detail rows
    remainingQuota: 'Remaining Quota',
    expiresAt: 'Expires At',
    todayExpires: '(expires today)',
    daysLeft: '({days} days)',
    usedQuota: 'Used Quota',
    resetNow: 'Resetting soon',
    subscriptionType: 'Subscription Type',
    subscriptionExpires: 'Subscription Expires',
    // Usage stat cells
    todayRequests: 'Today Requests',
    todayInputTokens: 'Today Input',
    todayOutputTokens: 'Today Output',
    todayTokens: 'Today Tokens',
    todayCacheCreation: 'Today Cache Creation',
    todayCacheRead: 'Today Cache Read',
    todayCost: 'Today Cost',
    rpmTpm: 'RPM / TPM',
    totalRequests: 'Total Requests',
    totalInputTokens: 'Total Input',
    totalOutputTokens: 'Total Output',
    totalTokensLabel: 'Total Tokens',
    totalCacheCreation: 'Total Cache Creation',
    totalCacheRead: 'Total Cache Read',
    totalCost: 'Total Cost',
    avgDuration: 'Avg Duration',
    // Messages
    enterApiKey: 'Please enter an API Key',
    querySuccess: 'Query successful',
    queryFailed: 'Query failed',
    queryFailedRetry: 'Query failed, please try again later',
    noDailyUsage: 'No daily usage data',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API Setup',
    description: 'Configure your Sub2API instance',
    database: {
      title: 'Database Configuration',
      description: 'Connect to your PostgreSQL database',
      host: 'Host',
      port: 'Port',
      username: 'Username',
      password: 'Password',
      databaseName: 'Database Name',
      sslMode: 'SSL Mode',
      passwordPlaceholder: 'Password',
      ssl: {
        disable: 'Disable',
        require: 'Require',
        verifyCa: 'Verify CA',
        verifyFull: 'Verify Full'
      }
    },
    redis: {
      title: 'Redis Configuration',
      description: 'Connect to your Redis server',
      host: 'Host',
      port: 'Port',
      password: 'Password (optional)',
      database: 'Database',
      passwordPlaceholder: 'Password',
      enableTls: 'Enable TLS',
      enableTlsHint: 'Use TLS when connecting to Redis (public CA certs)'
    },
    admin: {
      title: 'Admin Account',
      description: 'Create your administrator account',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      passwordPlaceholder: 'Min 8 characters',
      confirmPasswordPlaceholder: 'Confirm password',
      passwordMismatch: 'Passwords do not match'
    },
    ready: {
      title: 'Ready to Install',
      description: 'Review your configuration and complete setup',
      database: 'Database',
      redis: 'Redis',
      adminEmail: 'Admin Email'
    },
    status: {
      testing: 'Testing...',
      success: 'Connection Successful',
      testConnection: 'Test Connection',
      installing: 'Installing...',
      completeInstallation: 'Complete Installation',
      completed: 'Installation completed!',
      redirecting: 'Redirecting to login page...',
      restarting: 'Service is restarting, please wait...',
      timeout: 'Service restart is taking longer than expected. Please refresh the page manually.'
    }
  },

  // Common
}
