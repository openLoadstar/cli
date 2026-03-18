<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_checkpoint
## [STATUS] S_STB
## [SYNCED_AT] 2026-03-18

### 1. DESCRIPTION
- SUMMARY: `loadstar checkpoint -m "[Message]"` 구현. git commit → 커밋 해시 → ACTIVE SavePoint에 해시 기입을 원자적으로 처리.
- LINKED_WP: W://root/cli/cmd_checkpoint

### 2. CODE_MAP

**구현 후 (실측)**
- `cmd/checkpoint.go:17-49`
  - `checkpointCmd.Run()` → `-m` 플래그 파싱 → `gitClient.Commit(message)` 호출 → 실패 시 즉시 종료(Atomic) → SAVEPOINT 폴더 전체 스캔 → `S_ACT` 포함 파일에 `git: [hash]` append
- `cmd/checkpoint.go:51-92`
  - `historyCmd.Run()` → ADDRESS 파싱 → HISTORY/ 에서 `[dot-path]_*.md` 목록 수집 → 역순 정렬 → tabwriter 테이블 출력
- `cmd/checkpoint.go:94-141`
  - `diffCmd.Run()` → ADDRESS + H_ID 파싱 → 두 파일 읽기 → `diffmatchpatch`로 semantic diff → ANSI 컬러 출력 (green `+`, red `-`)
- `cmd/checkpoint.go:143-195`
  - `rollbackCmd.Run()` → ADDRESS + H_ID 파싱 → 확인 프롬프트(--force 우회) → pre_rollback 백업 → H_ID 스냅샷으로 현재 파일 덮어쓰기

### 3. ISSUES
- go-git의 `.loadstar/` 스테이징 범위가 프로젝트 전체 git repo와 충돌 가능 → `.loadstar/` 경로만 Add하는지 확인 완료 (client.go에서 `AddGlob(".loadstar/*")` 사용)

### 4. TODO
- [x] `gitClient.Commit(message)` 호출 및 해시 수신 [WP_REF:2]
- [x] Commit 실패 시 조기 반환 [WP_REF:5]
- [x] SAVEPOINT 폴더 스캔 → ACTIVE 파일 필터링 [WP_REF:4]
- [x] 해당 파일 `SAVEPOINTS` 섹션에 `git: [hash]` append [WP_REF:4]
- [x] historyCmd: HISTORY/ 스캔 후 역순 tabwriter 출력 [WP_REF:1]
- [x] diffCmd: sergi/go-diff로 unified diff 출력 [WP_REF:2]
- [x] rollbackCmd: pre_rollback 백업 후 복원 [WP_REF:3]

### 5. LOG
- [2026-03-18T00:00:00] [MODIFIED] WayPoint STATUS S_IDL → S_STB 갱신, CODE_MAP 구현 전 계획을 실측 라인번호(checkpoint:17-49, history:51-92, diff:94-141, rollback:143-195)로 교체, TODO 전체 [x] 처리, SYNCED_AT 갱신
</BLACKBOX>
