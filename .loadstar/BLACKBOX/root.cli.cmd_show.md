<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_show
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `loadstar show [ADDRESS] [--depth N]` 구현. 요소 파일 내용 출력 + depth N 트리 재귀 출력.
- LINKED_WP: W://root/cli/cmd_show

### CODE_MAP

**구현 후 (실측)**
- `cmd/nav.go:96-111`
  - `showCmd.Run()` → --depth 플래그 파싱 → ADDRESS 파싱 → visited map 초기화 → `showElement()` 호출
- `cmd/nav.go:113-150`
  - `showElement()` → 순환 참조 visited 체크 → 파일 존재 확인 → STATUS 추출 → `[address]  [STATUS]` 출력 → currentDepth < maxDepth이면 `extractContainsItems()`로 자식 추출 → 재귀 호출
- `cmd/nav.go:152-154`
  - `indent()` → depth * 2칸 공백 반환
- `cmd/nav.go:156-163`
  - `extractField()` → 정규식으로 헤더 필드 값 추출
- `cmd/nav.go:165-173`
  - `extractContainsItems()` → `ITEMS: [...]` 정규식으로 주소 목록 추출 → `,` 분리

### TODO
- [x] ADDRESS 파싱 → 파일 읽기 → 내용 출력 [WP_REF:1]
- [x] CONTAINS.ITEMS 정규식 추출 [WP_REF:2]
- [x] --depth N 재귀 트리 출력 [WP_REF:3]
- [x] 주소 + STATUS 코드 표시 [WP_REF:4]
- [x] visited map 순환 참조 방지 [WP_REF:5]

### ISSUE
- depth 재귀 시 파일 없는 주소는 `[NOT FOUND: addr]` 표시 후 스킵.

### COMMENT
- [2026-03-18T00:00:00] [MODIFIED] WayPoint STATUS S_IDL → S_STB 갱신, CODE_MAP 구현 전 계획을 실측 라인번호(nav.go:96-111, 113-150, 152-154, 156-163, 165-173)로 교체, TODO 전체 [x] 처리, SYNCED_AT 갱신
</BLACKBOX>
