<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_checkpoint
## [STATUS] S_PRG
## [SYNCED_AT] 2026-04-01

### 1. DESCRIPTION
- SUMMARY: `loadstar checkpoint -m "[Message]"` 구현. git commit → 커밋 해시 → ACTIVE SavePoint에 해시 기입을 원자적으로 처리. GIT_INDEX / CHANGE_LOG GIT_REF 연동은 미구현.
- LINKED_WP: W://root/cli/cmd_checkpoint

### 2. CODE_MAP

**구현 후 (실측) — 2026-04-01 help text 추가로 라인 이동**
- `cmd/checkpoint.go:18-72`
  - `checkpointCmd.Run()` → `-m` 플래그 파싱 → `gitClient.Commit(message)` 호출 → 실패 시 즉시 종료(Atomic) → SAVEPOINT 폴더 전체 스캔 → `S_ACT` 포함 파일에 `git: [hash]` append → 원격 push
- `cmd/checkpoint.go:74-127`
  - `historyCmd.Run()` → ADDRESS 파싱 → HISTORY/ 에서 `[dot-path]_*.md` 목록 수집 → 역순 정렬 → tabwriter 테이블 출력
- `cmd/checkpoint.go:129-185`
  - `diffCmd.Run()` → ADDRESS + H_ID 파싱 → 두 파일 읽기 → `diffmatchpatch`로 semantic diff → ANSI 컬러 출력 (green `+`, red `-`)
- `cmd/checkpoint.go:187-254`
  - `rollbackCmd.Run()` → ADDRESS + H_ID 파싱 → 확인 프롬프트(--force 우회) → pre_rollback 백업 → H_ID 스냅샷으로 현재 파일 덮어쓰기

### 3. ISSUES
- go-git의 `.loadstar/` 스테이징 범위가 프로젝트 전체 git repo와 충돌 가능 → `.loadstar/` 경로만 Add하는지 확인 완료 (client.go에서 `AddGlob(".loadstar/*")` 사용)

### 4. TODO
- [x] `gitClient.Commit(message)` 호출 및 해시 수신 [WP_REF:1]
- [x] Commit 실패 시 조기 반환 [WP_REF:2]
- [x] SAVEPOINT 폴더 스캔 → ACTIVE 파일 필터링 [WP_REF:3]
- [x] 해당 파일 `SAVEPOINTS` 섹션에 `git: [hash]` append [WP_REF:3]
- [x] git remote 설정 시 자동 push [WP_REF:4]
- [x] historyCmd: HISTORY/ 스캔 후 역순 tabwriter 출력
- [x] diffCmd: sergi/go-diff로 unified diff 출력
- [x] rollbackCmd: pre_rollback 백업 후 복원
- [ ] GIT_INDEX 파일 생성 (`GIT.[DATE].[SEQ].[UUID].md`) [WP_REF:5]
- [ ] CHANGE_LOG 파일들에 `GIT_REF` 역기입 [WP_REF:6]
- [ ] Commit Hash를 GIT_INDEX에 역기입 후 추가 커밋 [WP_REF:7]
- [ ] HISTORY/ 내 임시 스냅샷 정리 [WP_REF:8]

### 5. LOG
- [2026-04-01T17:36:07] [MODIFIED] SPEC 대조 후 WayPoint/BlackBox 갱신 — STATUS S_STB→S_PRG, TECH_SPEC 미구현 항목 4개 추가(GIT_INDEX 생성, CHANGE_LOG GIT_REF 역기입, GIT_INDEX 커밋 해시 역기입, HISTORY 정리), CODE_MAP 라인번호 갱신, SYNCED_AT 2026-04-01
- [2026-03-18T00:00:00] [MODIFIED] WayPoint STATUS S_IDL → S_STB 갱신, CODE_MAP 구현 전 계획을 실측 라인번호로 교체, TODO 전체 [x] 처리, SYNCED_AT 갱신
</BLACKBOX>
