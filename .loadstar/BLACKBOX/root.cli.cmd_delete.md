<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_delete
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `avcs delete [ADDRESS]` 구현. 삭제 전 H:// 백업(History First 원칙), 대상 파일 삭제, 부모 CONTAINS.ITEMS에서 해당 주소 제거.
- LINKED_WP: W://root/cli/cmd_delete

### CODE_MAP

**구현 후 (실측)**
- `cmd/element.go:132-170`
  - `deleteCmd.Run()` → --force 체크 → 확인 프롬프트 → History First(CopyFile, _deleted) → LINEAGE 파싱(parseLineageParent) → removeFromContains() → os.Remove
- `cmd/element.go:224-240`
  - `removeFromContains()` → CONTAINS.ITEMS 정규식 라인 스캔 후 특정 주소 제거
- `cmd/element.go:242-252`
  - `parseLineageParent()` → LINEAGE: [PARENT: ...] 정규식 추출

### TODO
- [x] ADDRESS 파싱 및 파일 존재 확인 [WP_REF:1]
- [x] `--force` 플래그 파싱 및 확인 프롬프트 구현 [WP_REF:6]
- [x] History First 백업 (`storage.CopyFile`, suffix: `_deleted`) [WP_REF:2]
- [x] 대상 파일 삭제 (`os.Remove`) [WP_REF:3]
- [x] LINEAGE.PARENT 정규식 파싱으로 부모 주소 추출 [WP_REF:4]
- [x] 부모 CONTAINS.ITEMS에서 해당 주소 제거 후 저장 [WP_REF:5]

### ISSUE
- LINEAGE.PARENT 파싱 실패 시: 경고 메시지 출력 후 계속 진행(파일 삭제는 수행).
- CONTAINS.ITEMS 파싱은 create와 동일한 라인 스캔 헬퍼 공유.

### COMMENT
- [2026-04-02T17:36:21] [MODIFIED] HISTORY 백업 제거, CONNECTIONS.PARENT 파싱으로 변경, avcs→loadstar 명칭 갱신
</BLACKBOX>
