const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("sub2apiDesktop", {
  isDesktop: true,
  getEnvironment: () => ipcRenderer.invoke("desktop:get-environment"),
  writeCodexConfig: (payload) => ipcRenderer.invoke("desktop:write-codex-config", payload),
  queryGatewayUsage: (payload) => ipcRenderer.invoke("desktop:query-gateway-usage", payload),
  openCodexDownload: (platform) => ipcRenderer.invoke("desktop:open-codex-download", platform),
  openConfigDir: () => ipcRenderer.invoke("desktop:open-config-dir"),
  getServerUrl: () => ipcRenderer.invoke("desktop:get-server-url"),
  setServerUrl: (url) => ipcRenderer.invoke("desktop:set-server-url", url)
});
