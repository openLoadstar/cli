<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_rollback
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar rollback [ADDRESS] [H_ID]` 구현. 특정 History 스냅샷으로 요소 파일을 복원. 단일 요소 범위에만 적용된다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/checkpoint.go` (rollbackCmd.Run)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/history_to_rollback | TYPE: L_SEQ,
    L://root/cli/diff_to_rollback | TYPE: L_SEQ,
    L://root/cli/delete_to_rollback | TYPE: L_REF,
    L://root/cli/checkpoint_to_test_checkpoint | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_checkpoint

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_rollback
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 → 현재 파일 경로 확인
  - [x] H_ID로 `HISTORY/[H_ID].md` 존재 확인
  - [x] 롤백 전 현재 상태를 Shadow History에 백업 (`[dot-path]_[TIMESTAMP]_pre_rollback.md`)
  - [x] H_ID 파일 내용을 현재 파일 경로에 덮어쓰기
  - [x] 확인 프롬프트 구현 (--force 플래그로 우회)
  - [x] 연쇄 복원 없음 안내 메시지 출력
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (pre_rollback 백업, CopyFile 복원, --force, 안내 메시지), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
