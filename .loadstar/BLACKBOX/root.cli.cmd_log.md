<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_log
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: log/findlog 명령 구현. BlackBox COMMENT 섹션에 로그 append, .clionly/LOG에 변경 이력 기록, 전체 BlackBox 스캔 기반 로그 조회.
- LINKED_WP: W://root/cli/cmd_log

### CODE_MAP
- `cmd/log.go:29-98` — `logCmd.Run()`: ADDRESS/KIND/CONTENT 파싱 → BlackBox 자동 생성 → COMMENT 섹션에 entry append → .clionly/LOG 기록
- `cmd/log.go:101-158` — `findlogCmd.Run()`: BLACKBOX/ 전체 스캔 → logLineRe 매칭 → --address/--kind 필터 → 최신순 정렬 → offset/limit 출력
- `cmd/log.go:167-169` — `bbPathFromLogicalPath()`: 논리 경로 → BlackBox 파일 경로 변환
- `cmd/log.go:201-224` — `appendLogToBlackBox()`: COMMENT 섹션 탐색 → entry 삽입
- `cmd/log.go:227-238` — `writeLogChangeLog()`: .clionly/LOG에 CL 파일 생성
- `cmd/log.go:249-327` — `collectLogEntries()`: BLACKBOX/ 디렉토리 스캔, 로그 라인 파싱

### TODO
- [x] logCmd 구현
- [x] findlogCmd 구현
- [x] BlackBox 자동 생성 (buildBlackBoxTemplate)
- [x] COMMENT 섹션 append (backward compat: "### 5. LOG" 지원)
- [x] .clionly/LOG 기록

### ISSUE
(없음)

### COMMENT
(없음)
</BLACKBOX>
