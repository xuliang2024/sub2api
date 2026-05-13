const fs = require("node:fs/promises");
const path = require("node:path");
const os = require("node:os");

const CONFIG_KEYS = new Set([
  "model_provider",
  "model",
  "review_model",
  "model_reasoning_effort",
  "disable_response_storage",
  "network_access",
  "windows_wsl_setup_acknowledged",
  "model_context_window",
  "model_auto_compact_token_limit"
]);

const FEATURE_KEYS = new Set(["responses_websockets_v2"]);
const VALID_MODES = new Set(["standard", "websocket"]);
const DEFAULT_CODEX_MODEL = "gpt-5.4";

function codexPaths(homeDir = os.homedir(), platform = process.platform) {
  const configDir = path.join(homeDir, ".codex");

  return {
    platform,
    homeDir,
    configDir,
    configPath: path.join(configDir, "config.toml"),
    authPath: path.join(configDir, "auth.json")
  };
}

function quoteTomlString(value) {
  return JSON.stringify(value);
}

function normalizeMode(mode) {
  return VALID_MODES.has(mode) ? mode : "standard";
}

function buildConfigBlock(baseUrl, mode = "standard") {
  const normalizedMode = normalizeMode(mode);
  const lines = [
    'model_provider = "OpenAI"',
    `model = ${quoteTomlString(DEFAULT_CODEX_MODEL)}`,
    `review_model = ${quoteTomlString(DEFAULT_CODEX_MODEL)}`,
    'model_reasoning_effort = "xhigh"',
    "disable_response_storage = true",
    'network_access = "enabled"',
    "windows_wsl_setup_acknowledged = true",
    "model_context_window = 1000000",
    "model_auto_compact_token_limit = 900000",
    "",
    "[model_providers.OpenAI]",
    'name = "OpenAI"',
    `base_url = ${quoteTomlString(baseUrl)}`,
    'wire_api = "responses"'
  ];

  if (normalizedMode === "websocket") {
    lines.push("supports_websockets = true");
  }

  lines.push("requires_openai_auth = true");

  if (normalizedMode === "websocket") {
    lines.push("", "[features]", "responses_websockets_v2 = true");
  }

  return lines.join("\n");
}

function stripManagedConfig(existingConfig) {
  const output = [];
  let currentSection = "";
  let skippingOpenAiProvider = false;

  for (const line of String(existingConfig || "").split(/\r?\n/)) {
    const sectionMatch = line.match(/^\s*\[([^\]]+)]\s*(?:#.*)?$/);

    if (sectionMatch) {
      const sectionName = sectionMatch[1].trim();
      currentSection = sectionName;

      if (sectionName === "model_providers.OpenAI") {
        skippingOpenAiProvider = true;
        continue;
      }

      skippingOpenAiProvider = false;
    }

    if (skippingOpenAiProvider) {
      continue;
    }

    if (!currentSection) {
      const keyMatch = line.match(/^\s*([A-Za-z0-9_-]+)\s*=/);
      if (keyMatch && CONFIG_KEYS.has(keyMatch[1])) {
        continue;
      }
    }

    if (currentSection === "features") {
      const keyMatch = line.match(/^\s*([A-Za-z0-9_-]+)\s*=/);
      if (keyMatch && FEATURE_KEYS.has(keyMatch[1])) {
        continue;
      }
    }

    output.push(line);
  }

  return output
    .join("\n")
    .replace(/\n{3,}/g, "\n\n")
    .trim();
}

async function pathExists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

function backupPath(filePath) {
  const stamp = new Date().toISOString().replace(/[-:T.Z]/g, "").slice(0, 14);
  return `${filePath}.bak-${stamp}`;
}

async function backupIfExists(filePath) {
  if (!(await pathExists(filePath))) {
    return null;
  }

  const destination = backupPath(filePath);
  await fs.copyFile(filePath, destination);
  return destination;
}

function normalizeAndValidateInput({ baseUrl, apiKey, mode }) {
  const normalizedBaseUrl = String(baseUrl || "").trim();
  const normalizedApiKey = String(apiKey || "").trim();
  const normalizedMode = normalizeMode(mode);

  if (!normalizedBaseUrl) {
    throw new Error("Base URL 不能为空。");
  }

  let parsed;
  try {
    parsed = new URL(normalizedBaseUrl);
  } catch {
    throw new Error("Base URL 必须是有效 URL。");
  }

  if (!["http:", "https:"].includes(parsed.protocol)) {
    throw new Error("Base URL 只支持 http 或 https。");
  }

  if (!normalizedApiKey) {
    throw new Error("API key 不能为空。");
  }

  return {
    baseUrl: normalizedBaseUrl.replace(/\/+$/, ""),
    apiKey: normalizedApiKey,
    mode: normalizedMode
  };
}

async function readJsonIfPossible(filePath) {
  if (!(await pathExists(filePath))) {
    return {};
  }

  try {
    const raw = await fs.readFile(filePath, "utf8");
    const parsed = JSON.parse(raw);
    return parsed && typeof parsed === "object" && !Array.isArray(parsed) ? parsed : {};
  } catch {
    return {};
  }
}

function parseTomlValue(rawValue) {
  const trimmed = String(rawValue || "").trim();

  if (!trimmed) {
    return "";
  }

  if (trimmed.startsWith('"')) {
    const match = trimmed.match(/^"(?:\\.|[^"\\])*"/);
    if (!match) {
      return "";
    }

    try {
      return JSON.parse(match[0]);
    } catch {
      return "";
    }
  }

  if (trimmed.startsWith("'")) {
    const endIndex = trimmed.indexOf("'", 1);
    return endIndex === -1 ? "" : trimmed.slice(1, endIndex);
  }

  return trimmed.split("#")[0].trim();
}

function parseTomlBoolean(rawValue) {
  return parseTomlValue(rawValue).toLowerCase() === "true";
}

function readOpenAiSettingsFromToml(configText) {
  let currentSection = "";
  let baseUrl = "";
  let supportsWebsockets = false;
  let websocketFeature = false;

  for (const line of String(configText || "").split(/\r?\n/)) {
    const sectionMatch = line.match(/^\s*\[([^\]]+)]\s*(?:#.*)?$/);

    if (sectionMatch) {
      currentSection = sectionMatch[1].trim();
      continue;
    }

    if (currentSection === "model_providers.OpenAI") {
      const baseUrlMatch = line.match(/^\s*base_url\s*=\s*(.+)$/);
      if (baseUrlMatch) {
        baseUrl = parseTomlValue(baseUrlMatch[1]);
      }

      const websocketMatch = line.match(/^\s*supports_websockets\s*=\s*(.+)$/);
      if (websocketMatch) {
        supportsWebsockets = parseTomlBoolean(websocketMatch[1]);
      }
    }

    if (currentSection === "features") {
      const featureMatch = line.match(/^\s*responses_websockets_v2\s*=\s*(.+)$/);
      if (featureMatch) {
        websocketFeature = parseTomlBoolean(featureMatch[1]);
      }
    }
  }

  return {
    baseUrl,
    mode: supportsWebsockets || websocketFeature ? "websocket" : "standard"
  };
}

function readOpenAiBaseUrlFromToml(configText) {
  return readOpenAiSettingsFromToml(configText).baseUrl;
}

async function readExistingSettings(paths) {
  const [configText, auth] = await Promise.all([
    pathExists(paths.configPath).then((exists) => exists ? fs.readFile(paths.configPath, "utf8") : ""),
    readJsonIfPossible(paths.authPath)
  ]);
  const openAiSettings = readOpenAiSettingsFromToml(configText);

  return {
    baseUrl: openAiSettings.baseUrl,
    mode: openAiSettings.mode,
    apiKey: typeof auth.OPENAI_API_KEY === "string" ? auth.OPENAI_API_KEY : "",
    hasChatGptTokens: Boolean(auth.tokens?.access_token || auth.tokens?.refresh_token)
  };
}

async function getEnvironment(options = {}) {
  const paths = codexPaths(options.homeDir, options.platform);
  const [hasConfig, hasAuth, existingSettings] = await Promise.all([
    pathExists(paths.configPath),
    pathExists(paths.authPath),
    readExistingSettings(paths)
  ]);

  return {
    ...paths,
    hasConfig,
    hasAuth,
    existingBaseUrl: existingSettings.baseUrl,
    existingMode: existingSettings.mode,
    existingApiKey: existingSettings.apiKey,
    hasExistingApiKey: Boolean(existingSettings.apiKey),
    hasChatGptTokens: existingSettings.hasChatGptTokens
  };
}

function redactKey(apiKey) {
  if (apiKey.length <= 10) {
    return "********";
  }

  return `${apiKey.slice(0, 6)}${"*".repeat(Math.min(18, apiKey.length - 10))}${apiKey.slice(-4)}`;
}

async function writeCodexConfig(payload, options = {}) {
  const { baseUrl, apiKey, mode } = normalizeAndValidateInput(payload);
  const paths = codexPaths(options.homeDir, options.platform);

  await fs.mkdir(paths.configDir, { recursive: true });

  const [configBackup, authBackup] = await Promise.all([
    backupIfExists(paths.configPath),
    backupIfExists(paths.authPath)
  ]);

  let existingConfig = "";
  if (await pathExists(paths.configPath)) {
    existingConfig = await fs.readFile(paths.configPath, "utf8");
  }

  const preservedConfig = stripManagedConfig(existingConfig);
  const nextConfig = `${buildConfigBlock(baseUrl, mode)}${preservedConfig ? `\n\n${preservedConfig}` : ""}\n`;

  const existingAuth = await readJsonIfPossible(paths.authPath);
  const nextAuth = {
    ...existingAuth,
    OPENAI_API_KEY: apiKey
  };

  await fs.writeFile(paths.configPath, nextConfig, { encoding: "utf8", mode: 0o600 });
  await fs.writeFile(paths.authPath, `${JSON.stringify(nextAuth, null, 2)}\n`, {
    encoding: "utf8",
    mode: 0o600
  });

  await Promise.all([
    fs.chmod(paths.configPath, 0o600).catch(() => {}),
    fs.chmod(paths.authPath, 0o600).catch(() => {})
  ]);

  return {
    ...paths,
    mode,
    configBackup,
    authBackup,
    configPreview: buildConfigBlock(baseUrl, mode),
    authPreview: JSON.stringify({ OPENAI_API_KEY: redactKey(apiKey) }, null, 2)
  };
}

function buildGatewayUsageUrl(baseUrl) {
  const normalizedBaseUrl = String(baseUrl || "").trim().replace(/\/+$/, "");
  if (!normalizedBaseUrl) {
    throw new Error("Base URL 不能为空。");
  }

  let parsed;
  try {
    parsed = new URL(normalizedBaseUrl);
  } catch {
    throw new Error("Base URL 必须是有效 URL。");
  }

  if (!["http:", "https:"].includes(parsed.protocol)) {
    throw new Error("Base URL 只支持 http 或 https。");
  }

  const base = parsed.toString().replace(/\/+$/, "");
  return base.endsWith("/v1") ? `${base}/usage` : `${base}/v1/usage`;
}

async function queryGatewayUsage(payload, options = {}) {
  const apiKey = String(payload?.apiKey || "").trim();
  if (!apiKey) {
    throw new Error("API key 不能为空。");
  }

  const url = new URL(buildGatewayUsageUrl(payload?.baseUrl));
  const fetchImpl = options.fetchImpl || globalThis.fetch;

  if (payload?.startDate) {
    url.searchParams.set("start_date", String(payload.startDate));
  }
  if (payload?.endDate) {
    url.searchParams.set("end_date", String(payload.endDate));
  }

  if (typeof fetchImpl !== "function") {
    throw new Error("当前运行环境不支持 fetch。");
  }

  const response = await fetchImpl(url.toString(), {
    method: "GET",
    headers: {
      accept: "application/json",
      authorization: `Bearer ${apiKey}`
    }
  });
  const raw = await response.text();
  let data;

  try {
    data = raw ? JSON.parse(raw) : {};
  } catch {
    throw new Error("用量接口返回了非 JSON 内容。");
  }

  if (!response.ok) {
    throw new Error(data.message || data.error?.message || data.error || `用量查询失败：HTTP ${response.status}`);
  }

  return data;
}

module.exports = {
  buildConfigBlock,
  buildGatewayUsageUrl,
  codexPaths,
  getEnvironment,
  normalizeAndValidateInput,
  parseTomlValue,
  queryGatewayUsage,
  readExistingSettings,
  readJsonIfPossible,
  readOpenAiBaseUrlFromToml,
  readOpenAiSettingsFromToml,
  redactKey,
  stripManagedConfig,
  writeCodexConfig
};
