<WAYPOINT>
## [ADDRESS] W://root/cli/test_nav
## [STATUS] S_IDL

### 1. IDENTITY
- SUMMARY: link/show 명령어 통합 테스트. Link 파일 생성·양방향 등록 및 show의 depth 재귀 출력을 검증한다.
- METADATA: [Ver: 1.0, Created: 2026-03-10, Priority: MEDIUM]
- SYNCED_AT: 2026-03-10

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/nav_test.go`
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/nav_to_test_nav | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/test_nav
- RESPONSE_STATUS: COMPLETED
- TECH_SPEC:
  - [x] `TestLink_FileCreated`: link 실행 후 LINK/ 폴더에 md 파일 생성 확인
  - [x] `TestLink_BidirectionalRegistration`: SOURCE와 TARGET 양쪽 CONNECTIONS.LINKS에 등록 확인
  - [x] `TestLink_InvalidType`: L_REF/L_SEQ/L_TST 외 타입 입력 시 에러 확인
  - [x] `TestLink_Duplicate`: 동일 조합 재등록 시 에러 확인
  - [x] `TestShow_Depth0`: depth 0에서 단일 요소 내용 출력 확인
  - [x] `TestShow_DepthN`: depth N에서 하위 요소 트리 재귀 출력 확인
  - [x] `TestShow_CircularRef`: 순환 참조 시 무한 루프 없이 종료 확인
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-11: [COMPLETED] 전체 테스트 55개 PASS (go test ./cmd/... -v)
  ]
</WAYPOINT>
