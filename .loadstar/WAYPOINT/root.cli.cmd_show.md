<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_show
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar show [FILTER]` 구현. 전체 WayPoint 목록을 ADDRESS + STATUS 테이블로 출력하고, FILTER 지정 시 주소에 해당 키워드가 포함된 항목만 필터링한다.
- METADATA: [Ver: 2.0, Created: 2026-03-04, Priority: MEDIUM]
- SYNCED_AT: 2026-04-08

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
- [x] 2026-04-08 WAYPOINT 디렉토리 스캔 → 전체 목록 출력
- [x] 2026-04-08 FILTER 인자로 주소 키워드 필터링 (대소문자 무시)
- [x] 2026-04-08 tabwriter로 ADDRESS + STATUS 테이블 출력

### ISSUE
(없음)

### COMMENT
v2.0: 기존 단일 요소 트리 표시 → 전체 WP 목록 + 필터로 개편 (2026-04-08)
</WAYPOINT>
