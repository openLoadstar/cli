<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_edit
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `avcs edit [ADDRESS]` 구현. 수정 전 Shadow History 스냅샷을 자동 생성하고, 환경변수 기반 에디터를 실행한다.
- LINKED_WP: W://root/cli/cmd_edit

### CODE_MAP

**구현 후 (실측)**
- `cmd/element.go:78-128`
  - `editCmd.Run()` → ADDRESS 파싱 → fs.Exists() → Shadow History(CopyFile) → resolveEditor() → os/exec 실행 → mtime 비교
- `cmd/element.go:253-264`
  - `resolveEditor()` → LOADSTAR_EDITOR → EDITOR → notepad(Windows)/vi(나머지)

### TODO
- [x] ADDRESS 파싱 및 파일 존재 확인 [WP_REF:1]
- [x] Shadow History 스냅샷 생성 (`storage.CopyFile`) [WP_REF:2]
- [x] 에디터 우선순위 환경변수 체인 구현 [WP_REF:3]
- [x] `os/exec`로 에디터 프로세스 실행 [WP_REF:4]
- [x] 에디터 종료 후 mtime 비교로 변경 감지 [WP_REF:5]

### ISSUE
- 에디터 우선순위: `AVCS_EDITOR` → `EDITOR` → OS 기본값(Windows: notepad)
- Shadow History 파일명: `HISTORY/root.cli.cmd_edit_20260305T153000.md`

### COMMENT
- [2026-04-02T17:36:21] [MODIFIED] Shadow History 제거, LOADSTAR_EDITOR 우선순위 갱신, avcs→loadstar 명칭 갱신
</BLACKBOX>
