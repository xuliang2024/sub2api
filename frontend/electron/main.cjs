const { app, BrowserWindow, Menu, ipcMain, shell } = require("electron");
const fs = require("node:fs/promises");
const path = require("node:path");
const {
  getEnvironment,
  queryGatewayUsage,
  writeCodexConfig
} = require("./configWriter.cjs");
const {
  DEFAULT_SERVER_URL,
  isTrustedSenderUrl,
  normalizeServerUrl,
  resolveInitialServerUrl
} = require("./security.cjs");

const FALLBACK_DESKTOP_MANIFEST = {
  codexDownloads: {
    mac: "https://cdn-video.51sux.com/downloads/codex/Codex-1-20260512-215253.dmg",
    windows: "https://codex.apiz.ai/download/windows",
    linux: "https://codex.apiz.ai/download"
  }
};

let mainWindow = null;
let serverUrl = DEFAULT_SERVER_URL;

function settingsPath() {
  return path.join(app.getPath("userData"), "desktop-settings.json");
}

async function readStoredServerUrl() {
  try {
    const raw = await fs.readFile(settingsPath(), "utf8");
    const parsed = JSON.parse(raw);
    if (typeof parsed?.serverUrl === "string") {
      return normalizeServerUrl(parsed.serverUrl);
    }
  } catch {
    // Ignore missing or malformed settings and fall back below.
  }

  return null;
}

async function writeStoredServerUrl(nextServerUrl) {
  await fs.mkdir(app.getPath("userData"), { recursive: true });
  await fs.writeFile(
    settingsPath(),
    `${JSON.stringify({ serverUrl: nextServerUrl }, null, 2)}\n`,
    "utf8"
  );
}

function assertTrustedSender(event) {
  const senderUrl = event.senderFrame?.url || event.sender.getURL();
  if (!isTrustedSenderUrl(senderUrl, serverUrl)) {
    throw new Error("当前页面不允许访问本机桌面配置能力。");
  }
}

function platformDownloadKey(platform = process.platform) {
  if (platform === "darwin" || platform === "mac") {
    return "mac";
  }
  if (platform === "win32" || platform === "windows") {
    return "windows";
  }
  return "linux";
}

async function fetchDesktopManifest() {
  const url = `${serverUrl.replace(/\/+$/, "")}/desktop-manifest.json`;
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 3000);

  try {
    const response = await fetch(url, {
      cache: "no-store",
      signal: controller.signal
    });
    if (!response.ok) {
      return FALLBACK_DESKTOP_MANIFEST;
    }
    const data = await response.json();
    return {
      ...FALLBACK_DESKTOP_MANIFEST,
      ...data,
      codexDownloads: {
        ...FALLBACK_DESKTOP_MANIFEST.codexDownloads,
        ...(data.codexDownloads || {})
      }
    };
  } catch {
    return FALLBACK_DESKTOP_MANIFEST;
  } finally {
    clearTimeout(timeout);
  }
}

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1180,
    height: 820,
    minWidth: 980,
    minHeight: 680,
    title: "Sub2API Desktop",
    backgroundColor: "#f8fafc",
    webPreferences: {
      preload: path.join(__dirname, "preload.cjs"),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: false
    }
  });

  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url);
    return { action: "deny" };
  });

  mainWindow.loadURL(serverUrl);

  mainWindow.on("closed", () => {
    mainWindow = null;
  });
}

function buildMenu() {
  return Menu.buildFromTemplate([
    ...(process.platform === "darwin"
      ? [{
          label: app.name,
          submenu: [
            { role: "about" },
            { type: "separator" },
            { role: "hide" },
            { role: "hideOthers" },
            { role: "unhide" },
            { type: "separator" },
            { role: "quit" }
          ]
        }]
      : []),
    {
      label: "File",
      submenu: [
        {
          label: "Open Codex Config Directory",
          click: async () => {
            const environment = await getEnvironment();
            await fs.mkdir(environment.configDir, { recursive: true });
            await shell.openPath(environment.configDir);
          }
        },
        {
          label: "Open Server in Browser",
          click: () => shell.openExternal(serverUrl)
        },
        { type: "separator" },
        process.platform === "darwin" ? { role: "close" } : { role: "quit" }
      ]
    },
    {
      label: "Edit",
      submenu: [
        { role: "undo" },
        { role: "redo" },
        { type: "separator" },
        { role: "cut" },
        { role: "copy" },
        { role: "paste" },
        { role: "selectAll" }
      ]
    },
    {
      label: "View",
      submenu: [
        { role: "reload" },
        { role: "toggleDevTools" },
        { type: "separator" },
        { role: "resetZoom" },
        { role: "zoomIn" },
        { role: "zoomOut" }
      ]
    }
  ]);
}

ipcMain.handle("desktop:get-environment", async (event) => {
  assertTrustedSender(event);
  return getEnvironment();
});

ipcMain.handle("desktop:write-codex-config", async (event, payload) => {
  assertTrustedSender(event);
  return writeCodexConfig(payload);
});

ipcMain.handle("desktop:query-gateway-usage", async (event, payload) => {
  assertTrustedSender(event);
  return queryGatewayUsage(payload);
});

ipcMain.handle("desktop:open-config-dir", async (event) => {
  assertTrustedSender(event);
  const environment = await getEnvironment();
  await fs.mkdir(environment.configDir, { recursive: true });
  const result = await shell.openPath(environment.configDir);
  if (result) {
    throw new Error(result);
  }
  return environment.configDir;
});

ipcMain.handle("desktop:open-codex-download", async (event, platform) => {
  assertTrustedSender(event);
  const manifest = await fetchDesktopManifest();
  const key = platformDownloadKey(platform || process.platform);
  const url = manifest.codexDownloads?.[key] || manifest.codexDownloads?.windows;
  if (!url) {
    throw new Error("未找到当前平台的 Codex 下载地址。");
  }
  await shell.openExternal(url);
  return url;
});

ipcMain.handle("desktop:get-server-url", async (event) => {
  assertTrustedSender(event);
  return serverUrl;
});

ipcMain.handle("desktop:set-server-url", async (event, value) => {
  assertTrustedSender(event);
  const nextServerUrl = normalizeServerUrl(value);
  serverUrl = nextServerUrl;
  await writeStoredServerUrl(nextServerUrl);
  if (mainWindow && !mainWindow.isDestroyed()) {
    await mainWindow.loadURL(nextServerUrl);
  }
  return serverUrl;
});

app.whenReady().then(async () => {
  serverUrl = process.env.SUB2API_DESKTOP_URL
    ? resolveInitialServerUrl({ envUrl: process.env.SUB2API_DESKTOP_URL, isPackaged: app.isPackaged })
    : ((await readStoredServerUrl()) || resolveInitialServerUrl({ isPackaged: app.isPackaged }));

  Menu.setApplicationMenu(buildMenu());
  createWindow();

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
