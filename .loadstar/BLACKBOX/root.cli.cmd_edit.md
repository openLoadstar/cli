<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_edit
## [STATUS] S_STB
## [SYNCED_AT] 2026-03-13

### 1. DESCRIPTION
- SUMMARY: `avcs edit [ADDRESS]` 구현. 수정 전 Shadow History 스냅샷을 자동 생성하고, 환경변수 기반 에디터를 실행한다.
- LINKED_WP: W://root/cli/cmd_edit

### 2. CODE_MAP

**구현 후 (실측)**
- `cmd/element.go:78-128`
  - `editCmd.Run()` → ADDRESS 파싱 → fs.Exists() → Shadow History(CopyFile) → resolveEditor() → os/exec 실행 → mtime 비교
- `cmd/element.go:253-264`
  - `resolveEditor()` → LOADSTAR_EDITOR → EDITOR → notepad(Windows)/vi(나머지)

### 3. ISSUES
- 에디터 우선순위: `AVCS_EDITOR` → `EDITOR` → OS 기본값(Windows: notepad)
- Shadow History 파일명: `HISTORY/root.cli.cmd_edit_20260305T153000.md`

### 4. TODO
- [x] ADDRESS 파싱 및 파일 존재 확인 [WP_REF:1]
- [x] Shadow History 스냅샷 생성 (`storage.CopyFile`) [WP_REF:2]
- [x] 에디터 우선순위 환경변수 체인 구현 [WP_REF:3]
- [x] `os/exec`로 에디터 프로세스 실행 [WP_REF:4]
- [x] 에디터 종료 후 mtime 비교로 변경 감지 [WP_REF:5]
</BLACKBOX>
