<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_checkpoint
## [STATUS] S_PRG

### 1. IDENTITY
- SUMMARY: `loadstar checkpoint -m "[Message]"` 구현. git commit + 메타데이터 스냅샷 + SavePoint 커밋 해시 기록을 원자적으로 처리하는 핵심 명령어.
- METADATA: [Ver: 1.2, Created: 2026-03-04, Priority: CRITICAL]
- SYNCED_AT: 2026-04-01

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/checkpoint.go` (checkpointCmd.Run)
  - 의존 패키지: `internal/git/client.go`, `internal/storage/fs.go`, `internal/address/address.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/infra_to_checkpoint | TYPE: L_REF,
    L://root/cli/checkpoint_to_test_checkpoint | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_checkpoint

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_checkpoint
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] `git.Client.Commit(message)` 호출 → 커밋 해시 반환
  - [x] git commit 실패 시 SavePoint 기록 중단 및 에러 반환 (Atomic 보장)
  - [x] ACTIVE 상태 SavePoint 파일의 SAVEPOINTS 섹션에 커밋 해시 자동 기입
  - [x] git remote 설정 시 자동 push
  - [ ] GIT_INDEX 파일 생성 (`GIT.[DATE].[SEQ].[UUID].md`) — SPEC §3 step 1
  - [ ] 해당 세션의 CHANGE_LOG 파일들에 `GIT_REF` 역기입 — SPEC §3 step 2
  - [ ] Commit Hash를 GIT_INDEX에 역기입 후 추가 커밋 — SPEC §3 step 4
  - [ ] HISTORY/ 내 임시 스냅샷 정리 — SPEC §3 step 6
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] SavePoint 업데이트 대상 범위: ACTIVE 전체 SavePoint 파일에 커밋 해시 기입.
  - [Q2 RESOLVED] HISTORY 스냅샷은 영구 보존. `history purge` 명령은 추후 구현.
  - [Q3] GIT_INDEX의 SEQ 채번 방식 — 같은 날짜 내 순번 관리 방법 미결.
- EXECUTION_HISTORY: [
    * 2026-03-18: 1차 구현 완료 (git commit, S_ACT SavePoint 해시 기입, Atomic 보장), STATUS S_IDL → S_STB, 테스트 PASS 확인
    * 2026-04-01: SPEC 대조 결과 GIT_INDEX/CHANGE_LOG GIT_REF/HISTORY 정리 미구현 확인, STATUS S_STB → S_PRG, TECH_SPEC 갱신
  ]
</WAYPOINT>
