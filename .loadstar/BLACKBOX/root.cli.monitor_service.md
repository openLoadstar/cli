<BLACKBOX>
## [ADDRESS] B://root/cli/monitor_service
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: loadstar_monitor.exe 구현. go-git의 ChangedLoadstarFiles()를 사용하여 변경 감지, flag 파일 기반으로 checkpoint와 연동.
- LINKED_WP: W://root/cli/monitor_service

### CODE_MAP
- `cmd/monitor/main.go:27-48` — `main()`: FindRoot → monitorDir 생성 → 무한 루프 (5분 간격)
- `cmd/monitor/main.go:51-100` — `checkAndFlag()`: ChangedLoadstarFiles() 호출 → monitoredDirs 필터 → flag 존재 확인 → flag 파일 생성 (시각 + 변경 파일 목록)

### TODO
- [x] 5분 간격 루프
- [x] go-git 기반 변경 감지
- [x] flag 파일 생성/중복 방지

### ISSUE
(없음)

### COMMENT
(없음)
</BLACKBOX>
