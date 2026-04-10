<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_log
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar log [TIME_RANGE] [FILTER]` / `loadstar log add` 구현. .clionly/LOG에 날짜별 파일로 메타 이벤트 로그 기록 및 검색.
- METADATA: [Ver: 2.0, Created: 2026-04-01]
- SYNCED_AT: 2026-04-08

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
#### v1.x (구버전)
- [x] 2026-04-01 loadstar log — BlackBox에 로그 append
- [x] 2026-04-01 loadstar findlog — BlackBox 스캔, 최신순 정렬, 필터링

#### v2.0 (현행 — log add/find, 날짜별 파일)
- [x] 2026-04-08 findlog 폐지, log add/find 서브커맨드 구조로 통합
- [x] 2026-04-08 log add [ADDR] [KIND] "[MSG]" — 날짜별 .log 파일에 append
- [x] 2026-04-08 log find [FILTER] [TIME] — 키워드/KIND 필터 + Nd/Nh 시간 범위
- [x] 2026-04-08 날짜별 파일 구조 (YYYY-MM-DD.log, 파이프 구분)
- [x] 2026-04-08 레거시 .log.md 파일 하위 호환 읽기
- [x] 2026-04-08 최대 출력 1000라인 제한
- [x] 2026-04-08 인자 없이 log 실행 시 help 표시
- [x] 2026-04-08 레거시 로그 파일 날짜별 파일로 마이그레이션

#### v3.0 (log 직접 조회 — find 서브커맨드 제거)
- [x] 2026-04-10 log [TIME_RANGE] [FILTER] 직접 조회로 변경 (find 서브커맨드 제거)

### ISSUE
(없음)

### COMMENT
v2.0: log+findlog 통합 → log add/find. BlackBox 기록 제거, .clionly/LOG/ 날짜별 파일 전용. 레거시 호환 유지.
</WAYPOINT>
