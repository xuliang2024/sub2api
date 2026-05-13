# codex.apiz.ai 源码修改与上线手册

本文档记录 `codex.apiz.ai` 当前采用的 Apiz 定制部署流程。目标是：本地改代码、GitHub 记录变更、云端从指定分支构建镜像、用 green 容器验证并切流，保留旧服务用于快速回滚。

## 当前环境

| 项目 | 值 |
| --- | --- |
| 本地仓库 | `/Users/m007/codes/sub2api-apiz` |
| 云端源码 | `/opt/sub2api-src` |
| 云端部署 | `/opt/sub2api-green/docker-compose.yml` |
| 线上域名 | `https://codex.apiz.ai` |
| 生产服务器 | `jp` / `8.216.36.146` |
| 当前工作分支 | `apiz/custom-home` |
| 你的仓库 | `https://github.com/xuliang2024/sub2api.git` |
| 上游仓库 | `https://github.com/Wei-Shaw/sub2api.git` |
| 原始基线 | `v0.1.123` / `df722c9a6e97312491232c11bf305d5f93b45e04` |

线上 Nginx 目前指向 green 容器：

```text
codex.apiz.ai -> nginx -> 127.0.0.1:8081 -> sub2api-green
```

旧容器仍保留在 `127.0.0.1:8080`：

```text
127.0.0.1:8080 -> sub2api
```

## 本地开发流程

进入本地仓库并确认分支：

```bash
cd /Users/m007/codes/sub2api-apiz
git status --short --branch
```

如果不在 `apiz/custom-home`，切换回来：

```bash
git switch apiz/custom-home
git pull --ff-only origin apiz/custom-home
```

首页主要入口：

```text
frontend/src/views/HomeView.vue
```

改完后提交并推送到你的仓库：

```bash
git add frontend/src/views/HomeView.vue
git commit -m "customize apiz home page"
git push origin apiz/custom-home
```

如改了依赖，必须同步 `frontend/pnpm-lock.yaml`：

```bash
cd frontend
pnpm install
cd ..
git add frontend/package.json frontend/pnpm-lock.yaml
git commit -m "chore: update frontend dependencies"
git push origin apiz/custom-home
```

## 云端上线流程

先确认线上服务正常：

```bash
curl -sS -w '\n%{http_code} %{time_total}\n' https://codex.apiz.ai/health
ssh jp 'cd /opt/sub2api-green && docker compose ps; cd /opt/sub2api && docker compose ps'
```

登录服务器并拉取你的分支：

```bash
ssh jp
cd /opt/sub2api-src
git status --short --branch
git fetch origin
git pull origin apiz/custom-home
```

云端源码可能保留过构建专用改动，正常情况下只应看到这些本地修改：

```text
M Dockerfile
M frontend/vite.config.ts
```

它们用于服务器构建稳定性：固定 `pnpm 9.15.9`、移除 Vite checker、限制 Go 构建并发。若 `git pull` 提示这些文件冲突，先停止上线，不要强行 reset。

生成新镜像 tag：

```bash
STAMP=$(date +%Y%m%d%H%M%S)
IMAGE_TAG="sub2api:apiz-custom-home-${STAMP}"
COMMIT=$(git rev-parse HEAD)
```

构建镜像：

```bash
docker build \
  --cpu-period=100000 --cpu-quota=100000 \
  --memory=4g --memory-swap=5g \
  --build-arg VERSION=0.1.123 \
  --build-arg COMMIT="$COMMIT" \
  --label org.opencontainers.image.version=0.1.123 \
  --label org.opencontainers.image.revision="$COMMIT" \
  --label org.opencontainers.image.source=https://github.com/xuliang2024/sub2api \
  -t "$IMAGE_TAG" \
  /opt/sub2api-src
```

更新 green Compose 的镜像 tag：

```bash
cd /opt/sub2api-green
cp docker-compose.yml "docker-compose.yml.bak.$(date +%Y%m%d%H%M%S)"
perl -0pi -e "s#image: sub2api:.*#image: ${IMAGE_TAG}#" docker-compose.yml
docker compose config >/tmp/sub2api-green-compose.check
```

只重建 green app 容器，不动 Postgres/Redis：

```bash
docker compose up -d --no-deps sub2api-green
```

验证 green：

```bash
curl -sS -H 'Host: codex.apiz.ai' http://127.0.0.1:8081/health
curl -sS -H 'Host: codex.apiz.ai' -o /tmp/green-home.html -w '%{http_code} %{time_total} %{size_download}\n' http://127.0.0.1:8081/home
docker logs --since 5m sub2api-green | grep -Ei 'panic|fatal|oom|killed|status_code":5|ERROR' || true
docker stats --no-stream sub2api-green
```

只在 green 验证通过后切流。如果 Nginx 已经指向 `8081`，通常不需要 reload：

```bash
grep -n 'proxy_pass' /etc/nginx/conf.d/codex-apiz.conf
```

如果需要从旧服务切到 green，执行：

```bash
cp /etc/nginx/conf.d/codex-apiz.conf "/etc/nginx/conf.d/codex-apiz.conf.bak.$(date +%Y%m%d%H%M%S)"
perl -0pi -e 's#proxy_pass http://127\.0\.0\.1:8080;#proxy_pass http://127.0.0.1:8081;#' /etc/nginx/conf.d/codex-apiz.conf
nginx -t
systemctl reload nginx
```

切流后外部验证：

```bash
curl -sS -w '\nhealth=%{http_code} %{time_total}\n' https://codex.apiz.ai/health
curl -sS -o /tmp/codex-home.html -w 'home=%{http_code} %{time_total} %{size_download}\n' https://codex.apiz.ai/home
curl -sS -o /tmp/codex-sub.yaml -w 'sub=%{http_code} %{time_total} %{size_download}\n' http://8.216.36.146/sub/f4b9931680e40f879b5e5765b023d8ac.yaml
```

观察 10 分钟：

```bash
for i in {1..20}; do
  date
  curl -sS -o /tmp/health.out -w 'health=%{http_code} time=%{time_total}\n' https://codex.apiz.ai/health
  curl -sS -o /tmp/home.out -w 'home=%{http_code} time=%{time_total} size=%{size_download}\n' https://codex.apiz.ai/home
  docker stats --no-stream sub2api-green
  sleep 30
done
```

## 回滚

秒级回滚到旧容器：

```bash
ssh jp 'perl -0pi -e "s#proxy_pass http://127\\.0\\.0\\.1:8081;#proxy_pass http://127.0.0.1:8080;#" /etc/nginx/conf.d/codex-apiz.conf && nginx -t && systemctl reload nginx'
```

验证旧服务：

```bash
curl -sS -w '\n%{http_code} %{time_total}\n' https://codex.apiz.ai/health
ssh jp 'curl -sS -H "Host: codex.apiz.ai" http://127.0.0.1/health'
```

如果要重新切回 green：

```bash
ssh jp 'perl -0pi -e "s#proxy_pass http://127\\.0\\.0\\.1:8080;#proxy_pass http://127.0.0.1:8081;#" /etc/nginx/conf.d/codex-apiz.conf && nginx -t && systemctl reload nginx'
```

不要默认恢复数据库。只有确认新版本执行了破坏性 schema/data 变更时，才使用上线前备份做数据库恢复。

## 备份与清理

上线前建议创建备份：

```bash
ssh jp 'set -e
STAMP=$(date +%Y%m%d%H%M%S)
BACKUP_DIR=/opt/sub2api-migration-backups/$STAMP
mkdir -p "$BACKUP_DIR"
cp /etc/nginx/conf.d/codex-apiz.conf "$BACKUP_DIR/codex-apiz.conf.before"
tar -C /opt/sub2api -czf "$BACKUP_DIR/sub2api-config-data-redis.tgz" .env docker-compose.yml data redis_data
docker exec sub2api-postgres pg_dump -U sub2api -d sub2api -Fc -f /tmp/sub2api-${STAMP}.dump
docker cp sub2api-postgres:/tmp/sub2api-${STAMP}.dump "$BACKUP_DIR/sub2api-${STAMP}.dump"
docker exec sub2api-postgres rm -f /tmp/sub2api-${STAMP}.dump
sha256sum "$BACKUP_DIR"/* > "$BACKUP_DIR/SHA256SUMS"
echo "$BACKUP_DIR"
'
```

稳定运行 24 小时后再考虑清理旧镜像。旧容器 `sub2api` 建议至少保留 24 小时，方便秒级回滚。

## 注意事项

- 不要在服务器上直接修改业务代码后不提交；业务改动应从本地提交到 `xuliang2024/sub2api`。
- 不要打印或提交 `/opt/sub2api/.env`、数据库 dump、API Key、OAuth token。
- 不要删除 `/opt/sub2api/data`、`/opt/sub2api/postgres_data`、`/opt/sub2api/redis_data`。
- 不要同时让多个不同版本长期写同一套数据库；green 验证和切流窗口应尽量短。
- 如果改动涉及数据库 migrations，先单独制定数据库迁移和恢复计划，不走本文档的秒级回滚简化流程。
