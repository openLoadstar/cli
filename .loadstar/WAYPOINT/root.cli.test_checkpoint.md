<WAYPOINT>
## [ADDRESS] W://root/cli/test_checkpoint
## [STATUS] S_IDL

### 1. IDENTITY
- SUMMARY: checkpoint/history/diff/rollback 명령어 통합 테스트. git 연동 없이 Mock GitClient를 사용하여 스냅샷 생성·조회·비교·복원 흐름을 검증한다.
- METADATA: [Ver: 1.0, Created: 2026-03-10, Priority: HIGH]
- SYNCED_AT: 2026-03-10

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/checkpoint_test.go`
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`, `internal/git/client.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/checkpoint_to_test_checkpoint | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/test_checkpoint
- RESPONSE_STATUS: COMPLETED
- TECH_SPEC:
  - [x] `MockGitClient` 구현: `internal.GitClient` 인터페이스 충족, Commit 호출 시 고정 해시 반환
  - [x] `TestCheckpoint_CommitAndSavePoint`: checkpoint 실행 후 ACTIVE SavePoint에 해시 기입 확인
  - [x] `TestCheckpoint_GitFailure`: Commit 실패 시 SavePoint 미수정 확인 (Atomic 보장)
  - [x] `TestHistory_List`: edit으로 스냅샷 생성 후 history 목록에 표시 확인
  - [x] `TestHistory_Empty`: 스냅샷 없을 때 안내 메시지 확인
  - [x] `TestDiff_Output`: 두 스냅샷 간 diff 출력 형식 확인 (+/- 라인 포함)
  - [x] `TestRollback_RestoresFile`: rollback 후 파일 내용이 스냅샷과 동일한지 확인
  - [x] `TestRollback_PreBackup`: rollback 전 현재 상태가 _pre_rollback 스냅샷으로 백업되는지 확인
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-11: [COMPLETED] 전체 테스트 55개 PASS (go test ./cmd/... -v)
  ]
</WAYPOINT>
