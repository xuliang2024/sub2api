<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-hidden bg-[#f7faf8] text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <!-- Background Decorations -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(15,23,42,0.045)_1px,transparent_1px),linear-gradient(90deg,rgba(15,23,42,0.045)_1px,transparent_1px)] bg-[size:72px_72px] dark:bg-[linear-gradient(rgba(148,163,184,0.08)_1px,transparent_1px),linear-gradient(90deg,rgba(148,163,184,0.08)_1px,transparent_1px)]"
      ></div>
      <div
        class="absolute inset-x-0 top-0 h-80 bg-[linear-gradient(180deg,rgba(20,184,166,0.16),transparent)] dark:bg-[linear-gradient(180deg,rgba(20,184,166,0.12),transparent)]"
      ></div>
      <div class="absolute inset-0 bg-[linear-gradient(135deg,transparent_0%,rgba(250,204,21,0.08)_42%,transparent_66%)] dark:bg-[linear-gradient(135deg,transparent_0%,rgba(250,204,21,0.04)_42%,transparent_66%)]"></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <div class="h-10 w-10 overflow-hidden rounded-lg shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-14">
      <div class="mx-auto w-full max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-stretch justify-between gap-12 lg:flex-row lg:items-center lg:gap-16">
          <!-- Left: Text Content -->
          <div class="w-full flex-1 text-center lg:flex-[1.12] lg:text-left">
            <div class="mb-4 inline-flex items-center gap-2 border border-primary-200 bg-white/80 px-3 py-1.5 text-xs font-semibold uppercase text-primary-700 shadow-sm backdrop-blur-sm dark:border-primary-800/70 dark:bg-dark-900/80 dark:text-primary-300">
              <Icon name="cpu" size="sm" />
              <span>{{ t('home.heroEyebrow') }}</span>
            </div>
            <h1
              class="mx-auto mb-4 max-w-3xl text-3xl font-bold leading-tight text-gray-950 dark:text-white sm:text-4xl md:text-5xl lg:mx-0 lg:text-5xl 2xl:text-6xl"
            >
              <span class="block sm:inline">{{ t('home.heroTitlePrefix') }}</span>
              <span class="block sm:inline sm:ml-3">{{ t('home.heroTitleSuffix') }}</span>
            </h1>
            <p class="mx-auto mb-4 max-w-[22rem] text-lg font-medium text-gray-700 dark:text-dark-200 sm:max-w-2xl md:text-xl lg:mx-0">
              {{ t('home.heroSubtitle') }}
            </p>
            <p class="mx-auto mb-8 max-w-[22rem] text-base leading-8 text-gray-600 dark:text-dark-300 sm:max-w-2xl lg:mx-0">
              {{ t('home.heroDescription') }}
            </p>

            <!-- CTA Button -->
            <div class="flex flex-col items-center gap-3 sm:flex-row lg:justify-start">
              <template v-if="isDesktopRuntime">
                <button
                  type="button"
                  class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-primary-500/30"
                  :disabled="checkingCodex"
                  @click="checkCodexStatus"
                >
                  <Icon :name="checkingCodex ? 'refresh' : 'terminal'" size="md" class="mr-2" :stroke-width="2" />
                  {{ t('home.desktopGuide.checkButton') }}
                </button>
                <button
                  type="button"
                  class="inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
                  @click="openCodexDownload"
                >
                  <Icon name="download" size="sm" />
                  {{ t('home.desktopGuide.downloadCodex') }}
                </button>
                <router-link
                  :to="isAuthenticated ? '/codex-setup' : '/login'"
                  class="inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
                >
                  {{ t('home.desktopGuide.importConfig') }}
                  <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
                </router-link>
              </template>
              <template v-else>
              <a
                :href="macDownloadUrl"
                target="_blank"
                rel="noopener noreferrer"
                download
                class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-primary-500/30"
              >
                <Icon name="download" size="md" class="mr-2" :stroke-width="2" />
                {{ t('home.download.macButton') }}
              </a>
              <a
                :href="windowsDownloadUrl"
                target="_blank"
                rel="noopener noreferrer"
                download
                class="inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
              >
                <Icon name="download" size="sm" />
                {{ t('home.download.windowsButton') }}
              </a>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
              >
                <Icon name="book" size="sm" />
                {{ t('home.viewDocs') }}
              </a>
              </template>
            </div>

            <div class="mx-auto mt-8 grid w-full max-w-xl grid-cols-1 gap-3 text-left sm:grid-cols-3 lg:mx-0">
              <div
                v-for="metric in heroMetrics"
                :key="metric.label"
                class="border border-gray-200 bg-white/65 px-4 py-3 backdrop-blur-sm dark:border-dark-700/70 dark:bg-dark-900/60"
              >
                <div class="text-lg font-bold text-gray-950 dark:text-white">{{ metric.value }}</div>
                <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ metric.label }}</div>
              </div>
            </div>
          </div>

          <!-- Right: Terminal Animation -->
          <div class="flex w-full flex-1 justify-center lg:flex-[0.88] lg:justify-end">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">codex relay</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">codex</span>
                    <span class="code-flag">--model</span>
                    <span class="code-url">gpt-5.3-codex</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># Relay pool scheduling...</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">200 OK</span>
                    <span class="code-response">{ "service": "ready" }</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Desktop Runtime Guide -->
        <section
          v-if="isDesktopRuntime"
          class="mb-12 border border-primary-200/80 bg-white/80 p-5 shadow-sm backdrop-blur-sm dark:border-primary-800/60 dark:bg-dark-900/80 md:p-6"
        >
          <div class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div class="flex items-start gap-4">
              <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-cyan-500 text-white shadow-lg shadow-primary-500/20"
              >
                <Icon name="terminal" size="lg" :stroke-width="2" />
              </div>
              <div>
                <div class="mb-2 inline-flex border border-primary-200 bg-primary-50 px-2.5 py-1 text-xs font-semibold text-primary-700 dark:border-primary-800/60 dark:bg-primary-900/30 dark:text-primary-300">
                  {{ t('home.desktopGuide.badge') }}
                </div>
                <h2 class="text-2xl font-bold text-gray-950 dark:text-white">
                  {{ t('home.desktopGuide.title') }}
                </h2>
                <p class="mt-2 max-w-2xl text-sm leading-6 text-gray-600 dark:text-dark-300">
                  {{ t('home.desktopGuide.description') }}
                </p>
              </div>
            </div>

            <div class="w-full shrink-0 border border-gray-200 bg-white px-4 py-3 text-left dark:border-dark-700 dark:bg-dark-950 lg:w-80">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <div class="text-xs uppercase text-gray-500 dark:text-dark-400">
                    {{ t('home.desktopGuide.codexStatusLabel') }}
                  </div>
                  <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">
                    {{ codexStatusText }}
                  </div>
                </div>
                <span
                  class="inline-flex h-9 w-9 items-center justify-center rounded-full"
                  :class="codexInstalled ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'"
                >
                  <Icon :name="codexInstalled ? 'check' : 'download'" size="sm" :stroke-width="2" />
                </span>
              </div>
              <div v-if="codexStatusDetail" class="mt-2 break-all text-xs leading-5 text-gray-500 dark:text-dark-400">
                {{ codexStatusDetail }}
              </div>
            </div>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-3">
            <div
              v-for="step in desktopGuideSteps"
              :key="step.title"
              class="border border-gray-200/80 bg-white/70 p-4 dark:border-dark-700/70 dark:bg-dark-950/60"
            >
              <div class="mb-3 flex h-9 w-9 items-center justify-center rounded-lg bg-gray-950 text-white dark:bg-white dark:text-gray-950">
                <Icon :name="step.icon" size="sm" :stroke-width="2" />
              </div>
              <div class="text-sm font-semibold text-gray-900 dark:text-white">{{ step.title }}</div>
              <div class="mt-1 text-xs leading-5 text-gray-500 dark:text-dark-400">{{ step.description }}</div>
              <button
                v-if="step.action === 'check'"
                type="button"
                class="mt-4 inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-3 py-2 text-xs font-semibold text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200"
                :disabled="checkingCodex"
                @click="checkCodexStatus"
              >
                <Icon name="refresh" size="xs" />
                {{ t('home.desktopGuide.recheck') }}
              </button>
              <button
                v-else-if="step.action === 'download'"
                type="button"
                class="mt-4 inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-3 py-2 text-xs font-semibold text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200"
                @click="openCodexDownload"
              >
                <Icon name="download" size="xs" />
                {{ t('home.desktopGuide.downloadCodex') }}
              </button>
              <router-link
                v-else
                :to="isAuthenticated ? '/codex-setup' : '/login'"
                class="mt-4 inline-flex items-center justify-center gap-2 border border-gray-300 bg-white px-3 py-2 text-xs font-semibold text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200"
              >
                <Icon name="arrowRight" size="xs" />
                {{ t('home.desktopGuide.importConfig') }}
              </router-link>
            </div>
          </div>
        </section>

        <!-- Desktop Download -->
        <section
          v-else
          class="mb-12 border border-primary-200/80 bg-white/80 p-5 shadow-sm backdrop-blur-sm dark:border-primary-800/60 dark:bg-dark-900/80 md:p-6"
        >
          <div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
            <div class="flex items-start gap-4">
              <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-cyan-500 text-white shadow-lg shadow-primary-500/20"
              >
                <Icon name="download" size="lg" :stroke-width="2" />
              </div>
              <div>
                <div class="mb-2 inline-flex border border-amber-200 bg-amber-50 px-2.5 py-1 text-xs font-semibold text-amber-700 dark:border-amber-800/60 dark:bg-amber-900/30 dark:text-amber-300">
                  {{ t('home.download.badge') }}
                </div>
                <h2 class="text-2xl font-bold text-gray-950 dark:text-white">
                  {{ t('home.download.title') }}
                </h2>
                <p class="mt-2 max-w-2xl text-sm leading-6 text-gray-600 dark:text-dark-300">
                  {{ t('home.download.description') }}
                </p>
              </div>
            </div>

            <div class="flex shrink-0 flex-col items-stretch gap-3 sm:flex-row sm:items-center">
              <div class="border border-gray-200 bg-white px-4 py-3 text-left dark:border-dark-700 dark:bg-dark-950">
                <div class="text-xs uppercase text-gray-500 dark:text-dark-400">
                  {{ t('home.download.versionLabel') }}
                </div>
                <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">
                  {{ t('home.download.version') }}
                </div>
                <div class="mt-1 max-w-64 break-all text-xs text-gray-500 dark:text-dark-400">
                  {{ t('home.download.integrity') }}
                </div>
              </div>
              <a
                :href="macDownloadUrl"
                target="_blank"
                rel="noopener noreferrer"
                download
                class="inline-flex items-center justify-center gap-2 rounded-lg bg-gray-950 px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-gray-950/15 transition-colors hover:bg-gray-800 dark:bg-white dark:text-gray-950 dark:hover:bg-dark-100"
              >
                <Icon name="download" size="sm" :stroke-width="2" />
                {{ t('home.download.macButton') }}
              </a>
              <a
                :href="windowsDownloadUrl"
                target="_blank"
                rel="noopener noreferrer"
                download
                class="inline-flex items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white px-5 py-3 text-sm font-semibold text-gray-800 shadow-sm transition-colors hover:border-primary-300 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-primary-700 dark:hover:text-primary-300"
              >
                <Icon name="download" size="sm" :stroke-width="2" />
                {{ t('home.download.windowsButton') }}
              </a>
            </div>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-3">
            <div
              v-for="step in downloadSteps"
              :key="step.title"
              class="border border-gray-200/80 bg-white/70 p-4 dark:border-dark-700/70 dark:bg-dark-950/60"
            >
              <div class="text-sm font-semibold text-gray-900 dark:text-white">{{ step.title }}</div>
              <div class="mt-1 text-xs leading-5 text-gray-500 dark:text-dark-400">{{ step.description }}</div>
            </div>
          </div>
        </section>

        <!-- Feature Tags - Centered -->
        <div class="mb-12 flex flex-wrap items-center justify-center gap-4 md:gap-6">
          <div
            class="inline-flex items-center gap-2.5 border border-gray-200/70 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-900/80"
          >
            <Icon name="cpu" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.subscriptionToApi')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 border border-gray-200/70 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-900/80"
          >
            <Icon name="terminal" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.stickySession')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 border border-gray-200/70 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-900/80"
          >
            <Icon name="shield" size="sm" class="text-primary-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.realtimeBilling')
            }}</span>
          </div>
        </div>

        <!-- Features Grid -->
        <div class="mb-12 grid gap-6 md:grid-cols-3">
          <!-- Feature 1: Unified Gateway -->
          <div
            class="group border border-gray-200/70 bg-white/70 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-900/70"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 shadow-lg shadow-blue-500/25 transition-transform group-hover:scale-105"
            >
              <Icon name="cpu" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.unifiedGateway') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.unifiedGatewayDesc') }}
            </p>
          </div>

          <!-- Feature 2: Account Pool -->
          <div
            class="group border border-gray-200/70 bg-white/70 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-900/70"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-emerald-600 shadow-lg shadow-primary-500/25 transition-transform group-hover:scale-105"
            >
              <Icon name="shield" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.multiAccount') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.multiAccountDesc') }}
            </p>
          </div>

          <!-- Feature 3: Billing & Quota -->
          <div
            class="group border border-gray-200/70 bg-white/70 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-900/70"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-amber-500 to-orange-500 shadow-lg shadow-amber-500/25 transition-transform group-hover:scale-105"
            >
              <Icon name="terminal" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.balanceQuota') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.balanceQuotaDesc') }}
            </p>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="mb-16 flex flex-wrap items-center justify-center gap-4">
          <!-- Claude - Supported -->
          <div
            class="flex items-center gap-2 rounded-lg border border-primary-200 bg-white/70 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-900/70"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-cyan-500"
            >
              <span class="text-xs font-bold text-white">C</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.codexCli') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- GPT - Supported -->
          <div
            class="flex items-center gap-2 rounded-lg border border-primary-200 bg-white/70 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-900/70"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-500"
            >
              <span class="text-xs font-bold text-white">D</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.codexDesktop') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Gemini - Supported -->
          <div
            class="flex items-center gap-2 rounded-lg border border-primary-200 bg-white/70 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-900/70"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-blue-600"
            >
              <span class="text-xs font-bold text-white">G</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.gptCodex') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Antigravity - Supported -->
          <div
            class="flex items-center gap-2 rounded-lg border border-primary-200 bg-white/70 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-900/70"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-rose-500 to-pink-600"
            >
              <span class="text-xs font-bold text-white">CC</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.claudeCode') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- More - Coming Soon -->
          <div
            class="flex items-center gap-2 rounded-lg border border-gray-200/70 bg-white/50 px-5 py-3 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-900/50"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-amber-500 to-orange-500"
            >
              <span class="text-xs font-bold text-white">S</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.privateDeploy') }}</span>
            <span
              class="rounded bg-amber-100 px-1.5 py-0.5 text-[10px] font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
              >{{ t('home.providers.service') }}</span
            >
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import type { CodexInstallStatus } from '@/types/global'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const heroMetrics = computed(() => [
  { value: t('home.metrics.relay.value'), label: t('home.metrics.relay.label') },
  { value: t('home.metrics.setup.value'), label: t('home.metrics.setup.label') },
  { value: t('home.metrics.billing.value'), label: t('home.metrics.billing.label') }
])

const downloadSteps = computed(() => [
  { title: t('home.download.steps.install.title'), description: t('home.download.steps.install.description') },
  { title: t('home.download.steps.login.title'), description: t('home.download.steps.login.description') },
  { title: t('home.download.steps.write.title'), description: t('home.download.steps.write.description') }
])

const desktopGuideSteps = computed(() => [
  {
    icon: 'terminal' as const,
    action: 'check',
    title: t('home.desktopGuide.steps.check.title'),
    description: t('home.desktopGuide.steps.check.description')
  },
  {
    icon: 'download' as const,
    action: 'download',
    title: t('home.desktopGuide.steps.download.title'),
    description: t('home.desktopGuide.steps.download.description')
  },
  {
    icon: 'cog' as const,
    action: 'import',
    title: t('home.desktopGuide.steps.import.title'),
    description: t('home.desktopGuide.steps.import.description')
  }
])

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'
const macDownloadUrl = 'https://pub-30b53688750140ac8432d5cea73f95e4.r2.dev/downloads/sub2api-desktop/Sub2API-Desktop-1.0.0-mac-arm64.zip'
const windowsDownloadUrl = 'https://get.microsoft.com/installer/download/9PLM9XGG6VKS?cid=website_cta_psi'
const isDesktopRuntime = computed(() => typeof window !== 'undefined' && !!window.sub2apiDesktop)
const checkingCodex = ref(false)
const codexStatus = ref<CodexInstallStatus | null>(null)
const codexStatusError = ref('')
const codexInstalled = computed(() => !!codexStatus.value?.installed)
const codexStatusText = computed(() => {
  if (checkingCodex.value) return t('home.desktopGuide.statusChecking')
  if (codexStatusError.value) return t('home.desktopGuide.statusError')
  if (!codexStatus.value) return t('home.desktopGuide.statusUnknown')
  return codexStatus.value.installed ? t('home.desktopGuide.statusInstalled') : t('home.desktopGuide.statusMissing')
})
const codexStatusDetail = computed(() => {
  if (codexStatusError.value) return codexStatusError.value
  if (!codexStatus.value) return t('home.desktopGuide.statusHint')
  if (!codexStatus.value.installed) return t('home.desktopGuide.missingHint')
  return codexStatus.value.version || codexStatus.value.executablePath
})

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

async function checkCodexStatus() {
  if (!window.sub2apiDesktop) return
  checkingCodex.value = true
  codexStatusError.value = ''
  try {
    codexStatus.value = await window.sub2apiDesktop.getCodexStatus()
  } catch (error: any) {
    codexStatusError.value = error?.message || t('home.desktopGuide.checkFailed')
  } finally {
    checkingCodex.value = false
  }
}

async function openCodexDownload() {
  if (!window.sub2apiDesktop) return
  try {
    await window.sub2apiDesktop.openCodexDownload()
  } catch (error: any) {
    codexStatusError.value = error?.message || t('home.desktopGuide.downloadFailed')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  if (isDesktopRuntime.value) {
    checkCodexStatus()
  }
})
</script>

<style scoped>
/* Terminal Container */
.terminal-container {
  position: relative;
  display: inline-block;
  max-width: 100%;
}

/* Terminal Window */
.terminal-window {
  width: calc(100vw - 48px);
  max-width: 460px;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 8px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.5s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #14b8a6;
}
.code-comment {
  color: #64748b;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fbbf24;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:deep(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(20, 184, 166, 0.2),
    0 0 40px rgba(20, 184, 166, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

@media (max-width: 640px) {
  .terminal-body {
    padding: 16px;
    font-size: 12px;
  }
}
</style>
