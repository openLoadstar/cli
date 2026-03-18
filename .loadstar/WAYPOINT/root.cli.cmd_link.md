<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_link
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar link [SOURCE] [TARGET] --type [L_REF|L_SEQ|L_TST]` 구현. 두 요소 간 논리적 Link 파일을 생성하고 양쪽 요소의 CONNECTIONS.LINKS에 등록한다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/nav.go` (linkCmd.Run)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`, `internal/core/element.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/link_to_show | TYPE: L_REF,
    L://root/cli/nav_to_test_nav | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_link

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_link
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] SOURCE, TARGET 주소 파싱 및 파일 존재 확인
  - [x] --type 유효성 확인 (L_REF, L_SEQ, L_TST)
  - [x] Link ID 자동 생성: `[source_id]_to_[target_id]`
  - [x] `LINK/[dot-path].md` 파일 생성 (SOURCE/TARGET/TYPE 기록)
  - [x] SOURCE 요소의 CONNECTIONS.LINKS에 `L://[ID] | TYPE: [type]` 추가
  - [x] TARGET 요소의 CONNECTIONS.LINKS에 역방향 참조 추가
  - [x] 동일 SOURCE-TARGET-TYPE 조합 중복 시 오류 반환
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] L_SEQ 역방향 참조는 TARGET의 LINKS에 동일 TYPE(L_SEQ)으로 기재한다.
- EXECUTION_HISTORY: [
    * 2026-03-18: 전체 구현 완료 확인 (Link md 생성, 양방향 CONNECTIONS.LINKS 등록, 중복 검사), STATUS S_IDL → S_STB, 테스트 PASS 확인
  ]
</WAYPOINT>
