<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_check
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar check` 구현. WP 파일 수정 시간과 git 최신 커밋 시간을 비교하여 동기화 필요 여부를 표시한다. (30분 유예, 최대 10개 표시, OUTDATED/SYNCED + GAP)
- METADATA: [Ver: 1.0, Created: 2026-04-10, Priority: HIGH]
- SYNCED_AT: 2026-04-10

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
- [x] 2026-04-10 WP SYNCED_AT와 git 최신 커밋 시간 비교 → 동기화 필요 여부 출력
- [x] 2026-04-10 CLAUDE.md 프롬프트에 세션 시작/종료 시 check 실행 규칙 추가
- [x] 2026-04-10 rootCmd에 checkCmd 등록 (cmd/root.go)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
