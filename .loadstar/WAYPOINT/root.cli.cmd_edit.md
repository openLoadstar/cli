<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_edit
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar edit [ADDRESS]` 구현. 기본 에디터를 실행하여 요소 파일 편집. 이력은 git이 관리.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/cmd_edit

### TODO
- [x] ADDRESS 파싱 → 물리 파일 경로 확인
- [x] 에디터 실행: `os/exec`로 프로세스 실행, Stdin/Stdout/Stderr 연결
- [x] 에디터 종료 후 변경 여부 감지 (수정 시각 비교)
- [x] 에디터 우선순위: LOADSTAR_EDITOR → EDITOR → OS 기본값(Windows: notepad)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
