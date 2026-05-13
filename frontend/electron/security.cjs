const DEFAULT_SERVER_URL = "https://codex.apiz.ai";
const DEV_SERVER_URL = "http://localhost:3000";

function normalizeServerUrl(value) {
  const raw = String(value || "").trim();
  if (!raw) {
    throw new Error("服务地址不能为空。");
  }

  let parsed;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error("服务地址必须是有效 URL。");
  }

  if (!["http:", "https:"].includes(parsed.protocol)) {
    throw new Error("服务地址只支持 http 或 https。");
  }

  parsed.username = "";
  parsed.password = "";
  parsed.search = "";
  parsed.hash = "";

  return parsed.toString().replace(/\/+$/, "");
}

function resolveInitialServerUrl({ envUrl, isPackaged } = {}) {
  if (envUrl) {
    return normalizeServerUrl(envUrl);
  }

  return isPackaged ? DEFAULT_SERVER_URL : DEV_SERVER_URL;
}

function originOf(value) {
  return new URL(normalizeServerUrl(value)).origin;
}

function isDesktopBridgePathAllowed(urlValue) {
  let parsed;
  try {
    parsed = new URL(String(urlValue || ""));
  } catch {
    return false;
  }

  const path = parsed.pathname.toLowerCase();
  return !(
    path === "/auth" ||
    path.startsWith("/auth/") ||
    path === "/payment" ||
    path.startsWith("/payment/") ||
    path === "/purchase" ||
    path.startsWith("/purchase/") ||
    path === "/orders" ||
    path.startsWith("/orders/")
  );
}

function isTrustedSenderUrl(senderUrl, serverUrl) {
  if (!senderUrl || !serverUrl || !isDesktopBridgePathAllowed(senderUrl)) {
    return false;
  }

  try {
    return new URL(senderUrl).origin === originOf(serverUrl);
  } catch {
    return false;
  }
}

module.exports = {
  DEFAULT_SERVER_URL,
  DEV_SERVER_URL,
  isDesktopBridgePathAllowed,
  isTrustedSenderUrl,
  normalizeServerUrl,
  originOf,
  resolveInitialServerUrl
};
