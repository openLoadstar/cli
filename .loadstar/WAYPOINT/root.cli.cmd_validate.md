<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_validate
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar validate` 구현. 모든 WayPoint와 Map의 링크를 추적하여 사라진 요소를 참조하는 깨진 링크를 검출한다.
- METADATA: [Ver: 1.0, Created: 2026-04-08, Priority: MEDIUM]
- SYNCED_AT: 2026-04-08

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
- [x] 2026-04-08 WAYPOINT 파일 스캔 → PARENT, CHILDREN, REFERENCE 주소 추출
- [x] 2026-04-08 MAP 파일 스캔 → WAYPOINTS 목록 추출
- [x] 2026-04-08 각 참조 주소의 파일 존재 여부 확인
- [x] 2026-04-08 깨진 참조 테이블 출력 (SOURCE, FIELD, BROKEN_ADDRESS)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
