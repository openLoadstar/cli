<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_delete
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `avcs delete [ADDRESS]` 명령어 구현. 삭제 전 H:// 백업 후 파일 삭제 및 부모 CONTAINS에서 제거.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-03-13

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/element.go` (deleteCmd.Run)
  - 의존 패키지: `internal/core/element.go`, `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/edit_to_delete | TYPE: L_REF,
    L://root/cli/delete_to_rollback | TYPE: L_REF,
    L://root/cli/delete_to_test_element | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_delete

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_delete
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 및 파일 존재 확인
  - [x] 삭제 전 `HISTORY/[dot-path]_[TIMESTAMP]_deleted.md`에 최종 상태 백업
  - [x] 대상 파일 삭제 (`os.Remove`)
  - [x] 파일 내 LINEAGE.PARENT 파싱하여 부모 주소 자동 추출
  - [x] 부모 요소 파일의 CONTAINS.ITEMS에서 해당 주소 제거 후 저장
  - [x] 확인 프롬프트 구현 (--force 플래그로 우회)
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] LINEAGE.PARENT 파싱 실패 시 경고 메시지 출력 후 계속 진행(파일 삭제는 수행).
- EXECUTION_HISTORY: [
    * 2026-03-13: 전체 구현 완료 (History First 백업, 삭제, LINEAGE 파싱, CONTAINS 제거, --force), STATUS S_IDL → S_STB
  ]
</WAYPOINT>
