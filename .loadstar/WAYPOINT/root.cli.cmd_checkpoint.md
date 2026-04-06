<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_checkpoint
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar checkpoint -m "[Message]" [--auto]` 구현. git commit + 변경 요소 자동 나열 커밋 메시지 + SavePoint 해시 기록 + monitor flag 삭제.
- METADATA: [Ver: 2.0, Created: 2026-03-04, Priority: CRITICAL]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/cmd_checkpoint

### TODO
- [x] `git.Client.Commit(message)` 호출 → 커밋 해시 반환
- [x] git commit 실패 시 SavePoint 기록 중단 및 에러 반환 (Atomic 보장)
- [x] ACTIVE 상태 SavePoint 파일의 SAVEPOINTS 섹션에 커밋 해시 자동 기입
- [x] git remote 설정 시 자동 push
- [x] 변경된 .loadstar/ 요소 주소를 커밋 메시지에 자동 append
- [x] `--auto` 플래그 → `[AUTO-CHECKPOINT]` 접두사 부여
- [x] checkpoint 완료 시 `.clionly/MONITOR/checkpoint_needed.flag` 삭제

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
