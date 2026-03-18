<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_diff
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar diff [ADDRESS] [H_ID]` 구현. 현재 요소 파일과 특정 History 스냅샷의 차이점을 unified diff 형식으로 출력.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/checkpoint.go` (diffCmd.Run)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/history_to_diff | TYPE: L_SEQ,
    L://root/cli/diff_to_rollback | TYPE: L_SEQ,
    L://root/cli/checkpoint_to_test_checkpoint | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_checkpoint

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_diff
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 → 현재 파일 경로 확인
  - [x] H_ID로 `HISTORY/[H_ID].md` 경로 확인
  - [x] 두 파일 내용을 라인 단위로 비교하여 unified diff 출력
  - [x] 컬러 출력: 추가(green `+`), 삭제(red `-`) — ANSI escape 직접 사용
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] diff 라이브러리: go.mod에 이미 포함된 `github.com/sergi/go-diff` 사용.
  - [Q2 RESOLVED] 컬러 출력은 ANSI escape 코드 직접 사용(외부 라이브러리 미추가).
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (go-diff semantic diff, ANSI 컬러 출력), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
