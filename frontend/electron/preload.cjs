const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("sub2apiDesktop", {
  isDesktop: true,
  appName: "codex 助手",
  getEnvironment: () => ipcRenderer.invoke("desktop:get-environment"),
  writeCodexConfig: (payload) => ipcRenderer.invoke("desktop:write-codex-config", payload),
  queryGatewayUsage: (payload) => ipcRenderer.invoke("desktop:query-gateway-usage", payload),
  openCodexDownload: (platform) => ipcRenderer.invoke("desktop:open-codex-download", platform),
  getCodexStatus: () => ipcRenderer.invoke("desktop:get-codex-status"),
  openConfigDir: () => ipcRenderer.invoke("desktop:open-config-dir"),
  getServerUrl: () => ipcRenderer.invoke("desktop:get-server-url"),
  setServerUrl: (url) => ipcRenderer.invoke("desktop:set-server-url", url)
});
