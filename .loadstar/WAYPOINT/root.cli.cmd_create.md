<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_create
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `avcs create [TYPE] [ID] --parent [PARENT_ID]` 명령어 구현. 새 요소 파일 생성, 부모 CONTAINS 등록, ID 중복 검증을 담당한다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/cmd_create

### TODO
- [x] --parent 플래그 파싱 및 부모 주소 유효성 검증
- [x] TYPE 유효성 확인 (M, W, L, S만 허용, H/B는 거부)
- [x] 동일 타입 폴더 내 ID 중복 검사 (`storage.Exists`)
- [x] ELEMENT_FORMAT 규격의 MD 파일 생성 (`storage.WriteFile`, Go text/template 사용)
- [x] 부모 요소 파일의 CONTAINS.ITEMS에 신규 주소 추가 후 저장

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
