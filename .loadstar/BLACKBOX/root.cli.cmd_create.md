<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_create
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `avcs create [TYPE] [ID] --parent [PARENT_ID]` 구현. 새 요소 MD 파일을 생성하고, 부모 요소의 CONTAINS.ITEMS에 자동 등록한다.
- LINKED_WP: W://root/cli/cmd_create

### CODE_MAP

**구현 후 (실측)**
- `cmd/element.go:27-72`
  - `createCmd.Run()` → TYPE 검증 → parent 파싱/존재확인 → 신규 주소 구성 → 중복검사 → buildTemplate() → appendToContains()
- `cmd/element.go:170-195`
  - `buildTemplate()` → M/W/L/S 타입별 fmt.Sprintf 템플릿 반환
- `cmd/element.go:197-222`
  - `appendToContains()` → CONTAINS.ITEMS 정규식 라인 스캔 후 append

### TODO
- [x] --parent 플래그 파싱 및 부모 주소 유효성 검증 [WP_REF:1]
- [x] TYPE 허용 목록 검증 (M, W, L, S) [WP_REF:2]
- [x] ID 중복 검사 (`storage.Exists`) [WP_REF:3]
- [x] 타입별 MD 템플릿 생성 (fmt.Sprintf) [WP_REF:4]
- [x] 부모 CONTAINS.ITEMS 라인 스캔 후 append [WP_REF:5]

### ISSUE
- `address.go:ToFilePath()`가 `a.ID`만 파일명으로 사용 중. 경로 전체를 dot-separated로 변환해야 함.
- 부모 CONTAINS.ITEMS 파싱은 라인 스캔 방식으로 결정됨 (Q2 RESOLVED).

### COMMENT
(없음)
</BLACKBOX>
