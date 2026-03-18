<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_history
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar history [ADDRESS]` 구현. 특정 요소의 HISTORY/ 폴더 내 스냅샷 목록을 시간순으로 출력.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/checkpoint.go` (historyCmd.Run)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/history_to_diff | TYPE: L_SEQ,
    L://root/cli/history_to_rollback | TYPE: L_SEQ,
    L://root/cli/edit_to_history | TYPE: L_REF,
    L://root/cli/checkpoint_to_test_checkpoint | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_checkpoint

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_history
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 → dot-path 추출
  - [x] `.loadstar/HISTORY/`에서 `[dot-path]_*.md` 패턴 파일 목록 수집
  - [x] 타임스탬프 기준 역순 정렬
  - [x] `tabwriter`로 정렬된 테이블 출력 (컬럼: H_ID, 파일 크기)
  - [x] 조회 결과 없을 경우 안내 메시지 출력
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (dot-path 추출, ListByPrefix, 역순 정렬, tabwriter 출력), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
