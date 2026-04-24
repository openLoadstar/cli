<WAYPOINT>
## [ADDRESS] W://root/maintenance/fix_claude_md
## [STATUS] S_STB

### IDENTITY
- SUMMARY: CLAUDE.md 문서 드리프트 정리 — 제거된 명령 참조 삭제, 현행 CLI 명령 목록 반영
- METADATA: [Priority: P2, Created: 2026-04-24]
- SYNCED_AT: 2026-04-24

### CONNECTIONS
- PARENT: M://root/maintenance
- CHILDREN: []
- REFERENCE: []

### TODO
- ADDRESS: W://root/maintenance/fix_claude_md
- SUMMARY: CLAUDE.md 명령어 목록 현행화 + 세션 진입 절차 정비
- TECH_SPEC:
  - [x] 2026-04-24 `findlog` 언급 삭제 (2026-04-10에 `log` 명령으로 통합됨)
  - [x] 2026-04-24 현행 명령어 목록 반영: `init`, `show`, `todo (sync/list/history)`, `log`, `validate`, `check`
  - [x] 2026-04-24 `loadstar check` 용도/출력 설명 추가 (git commit 대비 WP 수정시간 drift 검출)
  - [x] 2026-04-24 삭제된 명령 섹션에 `findlog` 이력 추가 (log 명령으로 통합)
  - [x] 2026-04-24 세션 시작 절차 번호 중복 교정 (4→5, 5→6)

### ISSUE
(없음)

### COMMENT
- 2026-04-24 SPEC 검토 세션 중 `/loadstar-enter` 테스트에서 드리프트 발견됨.
</WAYPOINT>
