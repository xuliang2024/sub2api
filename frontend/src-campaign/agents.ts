/**
 * 常用编码 Agent 的接入配置。
 *
 * 为什么不放进 campaign.code_samples：SDK 示例是「一次调用」，后台按活动配置
 * 即可；agent 是「写哪个文件 / 起哪个进程」，属于网关协议面的固定事实，所有活动
 * 都一样。内置在代码里还能和 SDK 示例共用 {{BASE_URL}}/{{API_KEY}}/{{MODEL}}
 * 占位替换——用户领到 Key 后，复制过去就是能直接跑的命令。
 *
 * 每条 code 都是一段可直接粘进终端的脚本，第一行注释写明落盘路径。
 */
import type { CampaignCodeSample } from './types'

export interface AgentSample extends CampaignCodeSample {
  /** 一句话说明（协议 + 配置位置）的文案 key。 */
  noteKey: string
}

export const AGENT_SAMPLES: AgentSample[] = [
  {
    id: 'codex',
    label: 'Codex CLI',
    lang: 'bash',
    noteKey: 'agents.note.codex',
    code: `# ~/.codex/config.toml
mkdir -p ~/.codex
cat > ~/.codex/config.toml <<'EOF'
model = "{{MODEL}}"
model_provider = "apiz"

[model_providers.apiz]
name = "apiz"
base_url = "{{BASE_URL}}/v1"
env_key = "APIZ_API_KEY"
wire_api = "chat"
EOF

export APIZ_API_KEY="{{API_KEY}}"
codex`,
  },
  {
    id: 'claude-code',
    label: 'Claude Code',
    lang: 'bash',
    noteKey: 'agents.note.claude',
    code: `# append to ~/.zshrc to keep it across sessions
export ANTHROPIC_BASE_URL="{{BASE_URL}}"
export ANTHROPIC_AUTH_TOKEN="{{API_KEY}}"
export ANTHROPIC_MODEL="{{MODEL}}"
# optional: send small/fast calls to the same model
# export ANTHROPIC_SMALL_FAST_MODEL="{{MODEL}}"

claude`,
  },
  {
    id: 'aider',
    label: 'Aider',
    lang: 'bash',
    noteKey: 'agents.note.aider',
    code: `export OPENAI_API_BASE="{{BASE_URL}}/v1"
export OPENAI_API_KEY="{{API_KEY}}"

aider --model openai/{{MODEL}}`,
  },
  {
    id: 'cline',
    label: 'Cline',
    lang: 'text',
    noteKey: 'agents.note.cline',
    code: `Provider   OpenAI Compatible
Base URL   {{BASE_URL}}/v1
API Key    {{API_KEY}}
Model ID   {{MODEL}}`,
  },
]
