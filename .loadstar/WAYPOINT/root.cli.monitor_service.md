<WAYPOINT>
## [ADDRESS] W://root/cli/monitor_service
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar_monitor.exe` 별도 프로세스. 5분 간격으로 .loadstar/ 변경 감시, checkpoint 필요 시 flag 파일 생성.
- METADATA: [Ver: 1.0, Created: 2026-04-02]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: [W://root/cli/cmd_checkpoint]
- BLACKBOX: B://root/cli/monitor_service

### TODO
- [x] 5분 간격 루프 실행
- [x] git status로 .loadstar/ 내 uncommitted 변경 확인
- [x] 감시 대상: WAYPOINT/, BLACKBOX/, MAP/
- [x] 변경 감지 시 .clionly/MONITOR/checkpoint_needed.flag 생성
- [x] flag 파일 중복 생성 방지

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
