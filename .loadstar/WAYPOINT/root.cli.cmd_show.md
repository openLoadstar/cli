<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_show
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar show [ADDRESS] [--depth N]` 구현. 요소 메타데이터를 터미널에 출력하고, --depth N 지정 시 하위 요소를 트리 형태로 재귀 출력한다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/nav.go` (showCmd.Run)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/link_to_show | TYPE: L_REF,
    L://root/cli/nav_to_test_nav | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_show

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_show
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 → 파일 읽기 → 내용 출력 (depth 0 기본)
  - [x] CONTAINS.ITEMS에서 주소 추출 (정규식 사용)
  - [x] --depth N: CONTAINS.ITEMS의 각 주소를 재귀적으로 읽어 트리 출력 (들여쓰기 2칸씩 증가)
  - [x] 각 노드: 주소 + STATUS 코드 표시
  - [x] 순환 참조 방지: visited map으로 이미 방문한 주소 추적
  - [x] depth 한도 초과 노드는 주소만 표시
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (showElement 재귀, extractContainsItems, visited map 순환 방지), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
