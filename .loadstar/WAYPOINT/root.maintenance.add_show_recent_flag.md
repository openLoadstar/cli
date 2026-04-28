<WAYPOINT>
## [ADDRESS] W://root/maintenance/add_show_recent_flag
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar show` 출력에 LAST_MODIFIED 컬럼 추가(파일 mtime 기반, 분 단위) + `--recent` 플래그로 최신 변경 순 정렬. 워크스페이스 CLAUDE.md의 폐기된 `loadstar check` 안내도 함께 정리.
- METADATA: [Priority: P3, Created: 2026-04-28]
- SYNCED_AT: 2026-04-28

### CONNECTIONS
- PARENT: M://root/maintenance
- CHILDREN: []
- REFERENCE: [W://root/cli/cmd_show, W://root/maintenance/remove_check_cmd]

### CODE_MAP
- scope:
  - cmd/

### TODO
- ADDRESS: W://root/maintenance/add_show_recent_flag
- SUMMARY: show 출력에 mtime 컬럼·--recent 정렬 옵션 추가, CLAUDE.md 잔여 check 안내 제거
- TECH_SPEC:
  # TASK
  - [x] 2026-04-28 cmd/nav.go: wpInfo 구조에 LastModified time.Time 필드 추가, 파일 mtime 캡처
  - [x] 2026-04-28 cmd/nav.go: tabwriter 출력에 LAST_MODIFIED 컬럼 추가 (YYYY-MM-DD HH:MM)
  - [x] 2026-04-28 cmd/nav.go: --recent 플래그 추가, true일 때 mtime 내림차순 정렬
  - [x] 2026-04-28 cmd/nav.go: --recent와 FILTER 조합 동작 보장
  - [x] 2026-04-28 cmd/nav.go: showCmd Long 도움말 갱신 (예시에 --recent 추가)
  - [x] 2026-04-28 cmd/nav_test.go: --recent 정렬 동작 + formatMTime 테스트 추가
  - [x] 2026-04-28 워크스페이스 CLAUDE.md(`C:\bono\MCP\GIT\CLAUDE.md`)의 "유용한 CLI 명령" 표에서 `loadstar check` 줄 제거, `loadstar show --recent` 안내 추가

  # RECURRING
  - (R) 변경 후 `go build -o bin/loadstar.exe .` 검증
  - (R) 변경 후 `go test ./...` 실행

### ISSUE
(없음)

### COMMENT
- 날짜 출처는 **파일 mtime** 사용 (사용자 결정). git checkout 후 mtime 리셋되는 한계가 있으나, "WP 변경"은 보통 작업 직후 의미가 살아있는 신호이므로 수용. LOG 타임스탬프 기반 대체는 추후 필요 시 별도 WP로 분리.
- 빌드/테스트 결과: `go build` OK, `go test ./...` PASS (cmd 패키지 모든 테스트 통과).
</WAYPOINT>
