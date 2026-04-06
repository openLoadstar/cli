<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_checkpoint
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `loadstar checkpoint -m "[Message]" [--auto]` 구현. git commit + 변경 요소 자동 나열 커밋 메시지 + SavePoint 해시 기록 + monitor flag 삭제. GIT_INDEX 폐지됨 (2026-04-02).
- LINKED_WP: W://root/cli/cmd_checkpoint

### CODE_MAP

**구현 후 (실측) — 2026-04-02 checkpoint 단순화 + --auto 플래그 추가**
- `cmd/checkpoint.go:18-95`
  - `checkpointCmd.Run()` → `-m` + `--auto` 플래그 파싱 → `ChangedLoadstarFiles()` 수집 → `buildCheckpointMessage()` 구성 → `gitClient.Commit(commitMsg)` 호출 → 실패 시 즉시 종료(Atomic) → SAVEPOINT 폴더 전체 스캔 → `S_ACT` 포함 파일에 `git: [hash]` append → `.clionly/MONITOR/checkpoint_needed.flag` 삭제 → 원격 push
- `cmd/checkpoint.go` `buildCheckpointMessage()`
  - 사용자 메시지 + `[AUTO-CHECKPOINT]` 접두사(auto 시) + 변경 요소 주소 목록 append
- `cmd/checkpoint.go` `historyCmd`, `diffCmd`, `rollbackCmd`
  - 변경 없음 (기존 HISTORY 스냅샷 기반 유지)

### TODO
- [x] `gitClient.Commit(message)` 호출 및 해시 수신
- [x] Commit 실패 시 조기 반환
- [x] SAVEPOINT 폴더 스캔 → ACTIVE 파일 필터링
- [x] 해당 파일 `SAVEPOINTS` 섹션에 `git: [hash]` append
- [x] git remote 설정 시 자동 push
- [x] 변경된 .loadstar/ 요소 주소를 커밋 메시지에 자동 append
- [x] `--auto` 플래그 → `[AUTO-CHECKPOINT]` 접두사
- [x] checkpoint 완료 시 `.clionly/MONITOR/checkpoint_needed.flag` 삭제
- [x] historyCmd: HISTORY/ 스캔 후 역순 tabwriter 출력
- [x] diffCmd: sergi/go-diff로 unified diff 출력
- [x] rollbackCmd: pre_rollback 백업 후 복원

### ISSUE
- go-git의 `.loadstar/` 스테이징 범위 확인 완료 (client.go에서 `AddGlob(".loadstar/*")` 사용)
- GIT_INDEX 폐지 결정 (2026-04-02): 커밋 메시지로 대체, 2단계 커밋 불필요

### COMMENT
- [2026-04-02T00:00:00] [MODIFIED] GIT_INDEX 폐지, checkpoint 단순화 — 커밋 메시지 자동 구성(변경 요소 나열) + --auto 플래그 + monitor flag 연동 추가, STATUS S_PRG→S_STB
- [2026-04-01T17:36:07] [MODIFIED] SPEC 대조 후 WayPoint/BlackBox 갱신 — STATUS S_STB→S_PRG, TECH_SPEC 미구현 항목 4개 추가(GIT_INDEX 생성, CHANGE_LOG GIT_REF 역기입, GIT_INDEX 커밋 해시 역기입, HISTORY 정리), CODE_MAP 라인번호 갱신, SYNCED_AT 2026-04-01
- [2026-03-18T00:00:00] [MODIFIED] WayPoint STATUS S_IDL → S_STB 갱신, CODE_MAP 구현 전 계획을 실측 라인번호로 교체, TODO 전체 [x] 처리, SYNCED_AT 갱신
</BLACKBOX>
