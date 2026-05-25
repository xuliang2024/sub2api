# codex.apiz.ai 图像生成 API 接入指南

本文档说明如何通过 `https://codex.apiz.ai` 调用 OpenAI Images 兼容接口，覆盖文生图、图生图、参数、响应、流式返回和错误格式。

## 基础信息

### Base URL

```text
https://codex.apiz.ai
```

### 鉴权

所有请求都需要在 Header 中携带 API Key：

```http
Authorization: Bearer sk-你的_API_Key
Content-Type: application/json
```

示例中使用环境变量保存 key，避免把密钥写进代码：

```bash
export APIZ_API_KEY="sk-你的_API_Key"
```

### 接口列表

| 功能 | Method | Path |
|------|--------|------|
| 文生图 | `POST` | `/v1/images/generations` |
| 图生图 / 改图 | `POST` | `/v1/images/edits` |

这两个接口是同步等待型接口：默认情况下，HTTP 请求会等待图片生成完成后直接返回图片结果，不是“提交任务 ID 后再轮询”的异步接口。

如果设置 `stream: true`，接口会使用 SSE 流式返回图片事件，但仍然是同一个 HTTP 连接，不是任务轮询。

## 模型

推荐使用：

```text
gpt-image-2
```

`model` 可省略，服务端默认使用 `gpt-image-2`。当前图片接口要求模型名以 `gpt-image-` 开头。

## 文生图

### 请求

```http
POST /v1/images/generations
```

### JSON 示例

```bash
curl -X POST "https://codex.apiz.ai/v1/images/generations" \
  -H "Authorization: Bearer ${APIZ_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "A clean product poster of a matte black wireless speaker on a white background, studio lighting, crisp shadow",
    "size": "1024x1024",
    "quality": "low",
    "output_format": "png",
    "response_format": "b64_json",
    "n": 1
  }'
```

### 参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `model` | string | 否 | `gpt-image-2` | 图片模型，需为 `gpt-image-*` |
| `prompt` | string | 是 | - | 图片描述 |
| `size` | string | 否 | 上游默认 | 图片尺寸，如 `1024x1024` |
| `n` | number | 否 | `1` | 返回图片数量，必须大于 0 |
| `quality` | string | 否 | 上游默认 | 可传 `low`、`medium`、`high`、`auto` 等，具体以模型支持为准 |
| `background` | string | 否 | 上游默认 | 可传 `transparent`、`opaque`、`auto` 等，具体以模型支持为准 |
| `output_format` | string | 否 | 上游默认 | `png`、`jpeg`、`webp` |
| `response_format` | string | 否 | `b64_json` | `b64_json` 或 `url` |
| `stream` | boolean | 否 | `false` | 是否使用 SSE 流式返回 |
| `moderation` | string | 否 | 上游默认 | 上游原生图像安全参数，按模型支持传入 |
| `style` | string | 否 | 上游默认 | 上游原生风格参数，按模型支持传入 |
| `output_compression` | number | 否 | 上游默认 | 输出压缩参数，常用于 `jpeg` / `webp` |
| `partial_images` | number | 否 | 上游默认 | 配合 `stream: true` 返回中间图，按模型支持传入 |

### 尺寸说明

服务端会透传 `size` 给上游。已在当前线上接口实测成功：

| size | 实测结果 |
|------|----------|
| `1024x1024` | 返回 PNG，实际 `1024 x 1024` |
| `1536x1024` | 返回 PNG，实际 `1536 x 1024` |

项目后端会识别以下尺寸用于计费分档，但是否能生成成功仍以当前上游模型实际支持为准：

| 示例尺寸 | 计费分档 |
|----------|----------|
| `1024x1024`、`1024x768` | `1K` |
| `1536x1024`、`1024x1536`、`2048x2048`、`2048x1152`、`1280x768` | `2K` |
| `2560x1440`、`2560x1600`、`3840x2160`、`2160x3840` | `4K` |
| `auto` 或未识别尺寸 | 默认按 `2K` 分档 |

## 图生图 / 改图

图生图接口支持两种输入方式：

- JSON：传图片 URL 或 data URL，适合已经有公网图片地址、Base64 图片的场景。
- multipart/form-data：直接上传本地图片文件，适合服务端或命令行直接传文件。

### JSON 请求

```http
POST /v1/images/edits
```

```bash
curl -X POST "https://codex.apiz.ai/v1/images/edits" \
  -H "Authorization: Bearer ${APIZ_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-image-2",
    "prompt": "Keep the product shape, change the background to a warm walnut desk with soft daylight",
    "images": [
      {
        "image_url": "https://example.com/source.png"
      }
    ],
    "size": "1024x1024",
    "input_fidelity": "high",
    "quality": "low",
    "output_format": "png",
    "response_format": "b64_json"
  }'
```

也可以传 data URL：

```json
{
  "images": [
    {
      "image_url": "data:image/png;base64,iVBORw0KGgo..."
    }
  ]
}
```

### 带 Mask 的 JSON 请求

```json
{
  "model": "gpt-image-2",
  "prompt": "Replace only the transparent area with a blue sky",
  "images": [
    {
      "image_url": "https://example.com/source.png"
    }
  ],
  "mask": {
    "image_url": "https://example.com/mask.png"
  },
  "size": "1024x1024",
  "response_format": "b64_json"
}
```

### multipart/form-data 请求

```bash
curl -X POST "https://codex.apiz.ai/v1/images/edits" \
  -H "Authorization: Bearer ${APIZ_API_KEY}" \
  -F "model=gpt-image-2" \
  -F "prompt=Keep the product, replace the background with a warm studio scene" \
  -F "image=@./source.png" \
  -F "size=1024x1024" \
  -F "quality=low" \
  -F "output_format=png" \
  -F "response_format=b64_json"
```

带 mask：

```bash
curl -X POST "https://codex.apiz.ai/v1/images/edits" \
  -H "Authorization: Bearer ${APIZ_API_KEY}" \
  -F "model=gpt-image-2" \
  -F "prompt=Only redraw the transparent area" \
  -F "image=@./source.png" \
  -F "mask=@./mask.png" \
  -F "size=1024x1024" \
  -F "response_format=b64_json"
```

多张输入图可用多个 `image` 字段，或使用 `image[0]`、`image[1]` 这类字段名。

### 图生图参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `model` | string | 否 | `gpt-image-2` | 图片模型 |
| `prompt` | string | 是 | - | 改图描述 |
| `images` | array | JSON 必填 | - | JSON 模式输入图片数组 |
| `images[].image_url` | string | JSON 必填 | - | 图片 URL 或 data URL |
| `image` | file | multipart 必填 | - | multipart 模式输入图片文件 |
| `mask` | object/file | 否 | - | JSON 模式用 `{ "image_url": "..." }`；multipart 模式用文件字段 |
| `size` | string | 否 | 上游默认 | 输出尺寸 |
| `n` | number | 否 | `1` | 返回图片数量 |
| `input_fidelity` | string | 否 | 上游默认 | 输入保持强度，如 `high`、`low`，以模型支持为准 |
| `quality` | string | 否 | 上游默认 | 输出质量 |
| `background` | string | 否 | 上游默认 | 背景策略 |
| `output_format` | string | 否 | 上游默认 | `png`、`jpeg`、`webp` |
| `response_format` | string | 否 | `b64_json` | `b64_json` 或 `url` |
| `stream` | boolean | 否 | `false` | 是否使用 SSE 流式返回 |
| `moderation` | string | 否 | 上游默认 | 上游原生安全参数 |
| `style` | string | 否 | 上游默认 | 上游原生风格参数 |
| `output_compression` | number | 否 | 上游默认 | 输出压缩参数 |
| `partial_images` | number | 否 | 上游默认 | 流式中间图数量 |

JSON 模式暂不支持 `images[].file_id` 和 `mask.file_id`，请使用 `image_url`。

## 非流式响应

默认 `stream` 为 `false`，接口返回 JSON。

### `response_format: "b64_json"`

```json
{
  "created": 1710000000,
  "data": [
    {
      "b64_json": "iVBORw0KGgoAAAANSUhEUg...",
      "revised_prompt": "A clean product poster featuring..."
    }
  ],
  "model": "gpt-image-2",
  "size": "1024x1024",
  "quality": "low",
  "output_format": "png",
  "usage": {
    "input_tokens": 66,
    "output_tokens": 781,
    "total_tokens": 847,
    "input_tokens_details": {
      "image_tokens": 0,
      "text_tokens": 66
    },
    "output_tokens_details": {
      "image_tokens": 781,
      "text_tokens": 0
    }
  }
}
```

字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| `created` | number | 创建时间戳 |
| `data` | array | 图片结果列表 |
| `data[].b64_json` | string | Base64 图片内容，不含 `data:image/...;base64,` 前缀 |
| `data[].url` | string | 当 `response_format` 为 `url` 时返回。当前兼容实现可能返回 data URL |
| `data[].revised_prompt` | string | 上游改写后的 prompt，可能为空 |
| `model` | string | 实际模型，可能为空或由上游返回 |
| `size` | string | 输出尺寸，可能为空或由上游返回 |
| `quality` | string | 输出质量，可能为空或由上游返回 |
| `output_format` | string | 输出格式，可能为空或由上游返回 |
| `usage` | object | token 用量，可能为空 |

### 保存 Base64 图片

```bash
jq -r '.data[0].b64_json' response.json | base64 --decode > output.png
```

## 流式响应

请求中传入：

```json
{
  "stream": true,
  "partial_images": 2
}
```

响应类型为：

```http
Content-Type: text/event-stream
```

### 文生图事件

可能返回：

```text
event: image_generation.partial_image
data: {"type":"image_generation.partial_image","created_at":1710000001,"partial_image_index":0,"b64_json":"...","output_format":"png","size":"1024x1024"}

event: image_generation.completed
data: {"type":"image_generation.completed","created_at":1710000001,"b64_json":"...","output_format":"png","size":"1024x1024","usage":{"input_tokens":5,"output_tokens":9}}
```

### 图生图事件

可能返回：

```text
event: image_edit.partial_image
data: {"type":"image_edit.partial_image","created_at":1710000001,"partial_image_index":0,"b64_json":"...","output_format":"png","size":"1024x1024"}

event: image_edit.completed
data: {"type":"image_edit.completed","created_at":1710000001,"b64_json":"...","revised_prompt":"replace background with aurora","output_format":"png","size":"1024x1024"}
```

### 流式错误事件

```text
event: error
data: {"type":"error","error":{"type":"upstream_error","message":"upstream request failed"}}
```

## 错误响应

常见错误格式：

```json
{
  "error": {
    "type": "invalid_request_error",
    "message": "images[].image_url is required"
  }
}
```

也可能返回简化格式：

```json
{
  "code": "INVALID_API_KEY",
  "message": "Invalid API key"
}
```

常见 HTTP 状态：

| 状态码 | 含义 | 常见原因 |
|--------|------|----------|
| `400` | 请求参数错误 | 缺少 `prompt`、图生图缺少 `images[].image_url`、字段类型错误 |
| `401` | 鉴权失败 | API Key 不存在或无效 |
| `403` | 权限不足 | 当前 key 所属分组未开启图片生成 |
| `404` | 接口不支持 | 非 OpenAI 平台分组调用 Images API，或路径错误 |
| `413` | 请求体过大 | 上传图片或 data URL 过大 |
| `429` | 限流 | 并发或上游限流 |
| `502` / `503` | 上游不可用 | 无可用账号、上游失败、调度失败 |

## JavaScript 示例

### 文生图

```js
const response = await fetch("https://codex.apiz.ai/v1/images/generations", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${process.env.APIZ_API_KEY}`,
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    model: "gpt-image-2",
    prompt: "A clean product poster of a matte black wireless speaker",
    size: "1024x1024",
    response_format: "b64_json",
    output_format: "png",
  }),
});

if (!response.ok) {
  throw new Error(await response.text());
}

const data = await response.json();
const base64 = data.data[0].b64_json;
```

### 图生图

```js
const response = await fetch("https://codex.apiz.ai/v1/images/edits", {
  method: "POST",
  headers: {
    Authorization: `Bearer ${process.env.APIZ_API_KEY}`,
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    model: "gpt-image-2",
    prompt: "Change the background to a warm studio scene",
    images: [
      {
        image_url: "https://example.com/source.png",
      },
    ],
    size: "1024x1024",
    input_fidelity: "high",
    response_format: "b64_json",
  }),
});

if (!response.ok) {
  throw new Error(await response.text());
}

const data = await response.json();
```

## 接入建议

- 服务端调用优先使用 `response_format: "b64_json"`，拿到后自行上传到对象存储或 CDN。
- 浏览器前端直连可能遇到 CORS 限制，建议通过自己的后端转发，或在本地开发时使用 dev proxy。
- 图生图使用 data URL 时请求体会变大，线上 Nginx 当前支持较大请求体，但仍建议图片先压缩到合理大小。
- 如果只需要最终图片，不要传 `stream: true`。
- 如果要展示生成进度，可使用 `stream: true` 配合 `partial_images`。
