<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_edit
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `avcs edit [ADDRESS]` 명령어 구현. 수정 전 Shadow History 스냅샷 생성 후 기본 에디터 실행.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-03-13

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/element.go` (editCmd.Run)
  - 의존 패키지: `internal/core/element.go`, `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/create_to_edit | TYPE: L_SEQ,
    L://root/cli/edit_to_history | TYPE: L_REF,
    L://root/cli/edit_to_test_element | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_edit

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_edit
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] ADDRESS 파싱 → 물리 파일 경로 확인 (`storage.Exists`)
  - [x] Shadow History 스냅샷 생성: `HISTORY/[dot-path]_[YYYYMMDDTHHMMSS].md` 복사
  - [x] 에디터 실행: `os/exec`로 프로세스 실행, Stdin/Stdout/Stderr 연결
  - [x] 에디터 종료 후 변경 여부 감지 (수정 시각 비교)
  - [x] 에디터 우선순위: AVCS_EDITOR → EDITOR → OS 기본값(Windows: notepad)
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] Windows 기본 에디터는 notepad로 고정. VS Code 사용 시 AVCS_EDITOR=code --wait 환경변수로 설정.
- EXECUTION_HISTORY: [
    * 2026-03-13: 전체 구현 완료 (Shadow History, 에디터 실행, mtime 감지), STATUS S_IDL → S_STB
  ]
</WAYPOINT>
