<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_checkpoint
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar checkpoint -m "[Message]"` 구현. git commit + 메타데이터 스냅샷 + SavePoint 커밋 해시 기록을 원자적으로 처리하는 핵심 명령어.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: CRITICAL]
- SYNCED_AT: 2026-03-18

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
  - [x] `.loadstar/` 전체 변경 파일 탐지
  - [x] `git.Client.Commit(message)` 호출 → 커밋 해시 반환
  - [x] 변경된 MD 파일들의 이전 버전을 `HISTORY/[ID]_[TIMESTAMP].md`에 스냅샷
  - [x] ACTIVE 상태 SavePoint 파일의 SAVEPOINTS 섹션에 커밋 해시 자동 기입
  - [x] git commit 실패 시 SavePoint 기록 중단 및 에러 반환 (Atomic 보장)
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] SavePoint 업데이트 대상 범위: ACTIVE 전체 SavePoint 파일에 커밋 해시 기입.
  - [Q2 RESOLVED] HISTORY 스냅샷은 영구 보존. `history purge` 명령은 추후 구현.
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (git commit, S_ACT SavePoint 해시 기입, Atomic 보장), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
