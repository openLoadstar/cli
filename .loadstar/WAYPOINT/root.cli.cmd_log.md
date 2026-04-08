<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_log
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar log` / `loadstar findlog` 구현. BlackBox에 구조화된 로그를 누적하고, .clionly/LOG에 변경 이력을 기록.
- METADATA: [Ver: 1.0, Created: 2026-04-01]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
- [x] loadstar log [ADDRESS] [KIND] "[내용]" — BlackBox에 로그 append
- [x] BlackBox 미존재 시 자동 생성
- [x] .clionly/LOG에 변경 이력 기록
- [x] loadstar findlog [OFFSET] [LIMIT] — 전체 BlackBox 스캔, 최신순 정렬
- [x] --address, --kind 필터링
- [x] 6종 KIND 지원 (NOTE, DECISION, ISSUE, RESOLVED, PROGRESS, MODIFIED)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
