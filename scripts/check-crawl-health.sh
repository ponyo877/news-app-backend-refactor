#!/usr/bin/env bash
# クロール停止検知(docs/GROWTH-PLAN.md Phase 0)
#
# /v1/site の last_updated_at が閾値(既定48時間)より古いサイトを列挙する。
# 「ワラノートが2年4ヶ月停止していたのに誰も気づかなかった」の再発防止。
# DB接続不要(公開APIのみ使用)なので、VMのcronでも手元のMacでも動く。
#
# 使い方:
#   ./check-crawl-health.sh                                          # 標準出力に表示のみ
#   WEBHOOK_URL=https://discord.com/api/webhooks/... ./check-crawl-health.sh   # 通知付き
#   THRESHOLD_HOURS=72 ./check-crawl-health.sh                       # 閾値変更
#
# cron設定例(VM上・毎朝9時JST):
#   0 9 * * * WEBHOOK_URL=... /home/keisuke877jp/check-crawl-health.sh >> /var/log/crawl-health.log 2>&1
#
# 終了コード: 停止サイトなし=0 / あり=1(CIや他スクリプトからの利用を想定)
set -euo pipefail

API_BASE="${API_BASE:-https://matome.folks-chat.com}"
THRESHOLD_HOURS="${THRESHOLD_HOURS:-48}"
WEBHOOK_URL="${WEBHOOK_URL:-}"

report=$(curl -fsS --max-time 30 "${API_BASE}/v1/site" | python3 -c "
import json, re, sys
from datetime import datetime, timedelta, timezone

# VM(Ubuntu 18.04)のpython3.6にはfromisoformatが無いため手でパースしてUTCに正規化する
def parse_utc(raw):
    m = re.match(r'(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.\d+)?(Z|[+-]\d{2}:?\d{2})?', raw)
    if not m:
        raise ValueError(raw)
    y, mo, d, h, mi, s = (int(v) for v in m.groups()[:6])
    tz = m.group(7)
    if tz in (None, 'Z'):
        offset = timedelta(0)
    else:
        sign = 1 if tz[0] == '+' else -1
        digits = tz[1:].replace(':', '')
        offset = sign * timedelta(hours=int(digits[:2]), minutes=int(digits[2:]))
    return datetime(y, mo, d, h, mi, s, tzinfo=timezone.utc) - offset

threshold = timedelta(hours=int('${THRESHOLD_HOURS}'))
now = datetime.now(timezone.utc)
sites = json.load(sys.stdin).get('data') or []

stale = []
for s in sites:
    raw = s.get('last_updated_at') or ''
    try:
        last = parse_utc(raw)
    except ValueError:
        stale.append((timedelta.max, s.get('titles', '?'), 'last_updated_at不明'))
        continue
    age = now - last
    if age > threshold:
        days = age.days
        label = f'{days}日停止' if days >= 1 else f'{int(age.total_seconds() // 3600)}時間停止'
        stale.append((age, s.get('titles', '?'), f'{label}(最終 {last.date()})'))

if not stale:
    print('OK')
else:
    print(f'{len(stale)}/{len(sites)} サイトが${THRESHOLD_HOURS}時間以上更新なし:')
    for _, title, detail in sorted(stale, key=lambda x: x[0], reverse=True):
        print(f'- {title}: {detail}')
")

echo "$report"

if [ "$report" = "OK" ]; then
  exit 0
fi

if [ -n "$WEBHOOK_URL" ]; then
  # Discordは content、Slack系は text を見る。両方入れてどちらでも通るようにする
  payload=$(printf '%s' "$report" | python3 -c "
import json, sys
body = '🕸️ まとめくんクロール停止検知\n' + sys.stdin.read()
print(json.dumps({'content': body, 'text': body}))
")
  curl -fsS --max-time 30 -H 'Content-Type: application/json' -d "$payload" "$WEBHOOK_URL" > /dev/null \
    || echo "警告: Webhook通知に失敗しました" >&2
fi

exit 1
