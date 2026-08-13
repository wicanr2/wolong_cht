#!/usr/bin/env bash
set -euo pipefail

# 本腳本只在 wolong-go Docker 容器內執行；不要在主機直接啟動 Go、Python 或 Xvfb。
repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_dir"

xvfb_pid=""
cleanup() {
  if [[ -n "$xvfb_pid" ]]; then
    kill "$xvfb_pid" 2>/dev/null || true
    wait "$xvfb_pid" 2>/dev/null || true
  fi
}
trap cleanup EXIT

if [[ -z "${DISPLAY:-}" ]]; then
  Xvfb :99 -screen 0 640x400x24 -nolisten tcp >/tmp/wolong-parity-xvfb.log 2>&1 &
  xvfb_pid=$!
  export DISPLAY=:99
fi

go test -p=1 -vet=off ./internal/state -run \
  'Test(ApproximateEvent10ProducerUsesKnownRawContract|ApproximateEvent10ProducerIsBoundedAndDisableable|ApproximateEvent10ReentersIdleClockConsumer|IdleClockDispatchesQueuedEvent10OnHourlyCadence|Event10ProducerWritesRawTalkPayload|QueuedEvent10TalkNotice|Event9LongNaturalRoute|Sub124FFMatchesRawSignedByteContract|MovingDisasterSub1248AUsesOnlyLastHalfOfRawSlots|MovingDisasterSub1248ARawWrapAndDirectionByte)$' \
  -count=1

go test -p=1 -vet=off ./cmd/wlgame -run \
  'Test(IdleClockGateRequiresStablePointerAndNoCommand|Event2To5TalkBranchParityGate|Event2To5FullTalkPageSampling|M7CorrectedTalkLayoutGate|ProjectileParityGate|Event9ShortFixtureGate|Event9LongNotificationRoute)$' \
  -count=1

go test -p=1 -vet=off ./internal/rules/tactical -run \
  'Test(SpecialProjectile.*|NormalProjectile.*|Projectile.*)$' \
  -count=1

python3 tools/talkdat_selftest.py

echo 'parity gate: PASS'
