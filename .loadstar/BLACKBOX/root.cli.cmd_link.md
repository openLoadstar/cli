<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_link
## [STATUS] S_STB
## [SYNCED_AT] 2026-03-18

### 1. DESCRIPTION
- SUMMARY: `loadstar link [SOURCE] [TARGET] --type [L_REF|L_SEQ|L_TST]` 구현. Link md 파일 생성 + 양방향 CONNECTIONS.LINKS 등록.
- LINKED_WP: W://root/cli/cmd_link

### 2. CODE_MAP

**구현 후 (실측)**
- `cmd/nav.go:15-17`
  - `allowedLinkTypes` → L_REF, L_SEQ, L_TST 허용 맵
- `cmd/nav.go:19-94`
  - `linkCmd.Run()` → --type 유효성 확인 → SOURCE/TARGET 주소 파싱 및 파일 존재 확인 → Link ID 생성(`[src_id]_to_[dst_id]`) → 중복 파일 존재 시 오류 → LINK/ md 파일 생성 → `appendToLinks()`로 SOURCE, TARGET 양방향 등록
- `cmd/nav.go:176-198`
  - `appendToLinks()` → 파일 읽기 → `CONNECTIONS.LINKS` 라인 스캔 정규식 → 기존 항목 있으면 `,\n    `로 append → 파일 저장

### 3. ISSUES
- CONNECTIONS.LINKS 파싱은 라인 스캔 방식 사용. 멀티라인 LINKS 섹션은 첫 번째 매칭 라인만 처리됨.

### 4. TODO
- [x] SOURCE, TARGET 주소 파싱 및 파일 존재 확인 [WP_REF:1]
- [x] --type 유효성 확인 [WP_REF:2]
- [x] Link ID 생성 및 LINK/ md 파일 생성 [WP_REF:3]
- [x] SOURCE CONNECTIONS.LINKS 등록 [WP_REF:5]
- [x] TARGET CONNECTIONS.LINKS 역방향 등록 [WP_REF:6]
- [x] 중복 조합 오류 반환 [WP_REF:7]

### 5. LOG
- [2026-03-18T00:00:00] [MODIFIED] WayPoint STATUS S_IDL → S_STB 갱신, CODE_MAP 구현 전 계획을 실측 라인번호(nav.go:15-17, 19-94, 176-198)로 교체, TODO 전체 [x] 처리, SYNCED_AT 갱신
</BLACKBOX>
