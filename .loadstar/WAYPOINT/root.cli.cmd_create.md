<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_create
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `avcs create [TYPE] [ID] --parent [PARENT_ID]` 명령어 구현. 새 요소 파일 생성, 부모 CONTAINS 등록, ID 중복 검증을 담당한다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-03-13

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/element.go` (createCmd.Run)
  - 의존 패키지: `internal/core/element.go`, `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/create_to_edit | TYPE: L_SEQ,
    L://root/cli/create_to_test_element | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_create

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_create
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] --parent 플래그 파싱 및 부모 주소 유효성 검증
  - [x] TYPE 유효성 확인 (M, W, L, S만 허용, H/B는 거부)
  - [x] 동일 타입 폴더 내 ID 중복 검사 (`storage.Exists`)
  - [x] ELEMENT_FORMAT 규격의 MD 파일 생성 (`storage.WriteFile`, Go text/template 사용)
  - [x] 부모 요소 파일의 CONTAINS.ITEMS에 신규 주소 추가 후 저장
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] MD 템플릿은 단순 문자열 포맷(fmt.Sprintf)으로 처리한다. text/template은 과도한 추상화.
  - [Q2 RESOLVED] 부모 CONTAINS.ITEMS 파싱은 라인 스캔 방식을 표준으로 채택한다. edit/delete에도 동일 방식 적용.
- EXECUTION_HISTORY: [
    * 2026-03-06: STATUS S_IDL → S_PRG, TECH_SPEC 1번 완료
    * 2026-03-13: 전체 구현 완료 (TYPE 검증, ID 중복, 템플릿 생성, CONTAINS 등록), STATUS S_PRG → S_STB
  ]
</WAYPOINT>
