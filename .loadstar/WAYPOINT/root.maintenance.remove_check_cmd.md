<WAYPOINT>
## [ADDRESS] W://root/maintenance/remove_check_cmd
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar check` 명령 완전 제거 — 원래 의도(코드↔WP 정합성)는 폐기됐고, 현재 기능은 `git status` + `loadstar show`로 대체 가능하여 중복
- METADATA: [Priority: P2, Created: 2026-04-24]
- SYNCED_AT: 2026-04-24

### CONNECTIONS
- PARENT: M://root/maintenance
- CHILDREN: []
- REFERENCE: []

### CODE_MAP
- scope:
  - cmd/

### TODO
- ADDRESS: W://root/maintenance/remove_check_cmd
- SUMMARY: check 명령 소스/WP/문서 전면 제거
- TECH_SPEC:
  - [x] 2026-04-24 cmd/check.go 삭제
  - [x] 2026-04-24 cmd/root.go에서 checkCmd AddCommand 제거
  - [x] 2026-04-24 W://root/cli/cmd_check WP 파일 삭제
  - [x] 2026-04-24 root.cli MAP의 WAYPOINTS 리스트에서 cmd_check 제거
  - [x] 2026-04-24 loadstar_cli/CLAUDE.md 세션 시작 절차에서 `loadstar check` 단계 제거
  - [x] 2026-04-24 loadstar_ui/CLAUDE.md 세션 시작 절차에서 `loadstar check` 단계 제거
  - [x] 2026-04-24 loadstar_SPEC/README.md에서 check 언급 제거
  - [x] 2026-04-24 go build 검증
  - (R) 변경 후 `go build -o bin/loadstar.exe .` 실행하여 컴파일 검증

### ISSUE
(없음)

### COMMENT
- "느슨한 관리(Tolerable Consistency)" 방침은 유지. 정합성 검사를 도구로 강제하지 않고 프롬프트/CLAUDE.md 기반 리마인더만 남긴다.
</WAYPOINT>
