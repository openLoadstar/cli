<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_show
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar show [ADDRESS] [--depth N]` 구현. 요소 메타데이터를 터미널에 출력하고, --depth N 지정 시 하위 요소를 트리 형태로 재귀 출력한다.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: MEDIUM]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/cmd_show

### TODO
- [x] ADDRESS 파싱 → 파일 읽기 → 내용 출력 (depth 0 기본)
- [x] CONTAINS.ITEMS에서 주소 추출 (정규식 사용)
- [x] --depth N: CONTAINS.ITEMS의 각 주소를 재귀적으로 읽어 트리 출력 (들여쓰기 2칸씩 증가)
- [x] 각 노드: 주소 + STATUS 코드 표시
- [x] 순환 참조 방지: visited map으로 이미 방문한 주소 추적
- [x] depth 한도 초과 노드는 주소만 표시

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
