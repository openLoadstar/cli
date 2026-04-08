<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_todo
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar todo sync/list/history` 서브커맨드 구현. WayPoint STATUS 기반으로 TODO_LIST를 자동 산출하며, TECH_SPEC [x] 항목으로 히스토리를 조회한다.
- METADATA: [Ver: 2.0, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-04-08

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
#### v1.x (구버전 — 수동 관리 방식)
- [x] 2026-03-10 TODO_LIST 마크다운 테이블 파서 구현
- [x] 2026-03-10 todo add: 신규 행 추가 (PENDING, --depends)
- [x] 2026-03-10 todo done: 행 삭제 → TODO_HISTORY append
- [x] 2026-03-10 todo list: PENDING/ACTIVE 필터, [BLOCKED] 표시
- [x] 2026-03-10 todo update: 상태 변경 + HISTORY 기록
- [x] 2026-03-10 todo delete: 삭제 + HISTORY 기록
- [x] 2026-03-10 todo history: 이력 조회 + 주소 필터

#### v2.0 (현행 — sync 기반 자동 관리)
- [x] 2026-04-08 add/update/done/delete 서브커맨드 제거
- [x] 2026-04-08 todo sync: MAP 순회 → 전체 WP 주소 수집
- [x] 2026-04-08 todo sync: WP_SNAPSHOT.json 캐시 기반 변경 감지 (modTime/size)
- [x] 2026-04-08 todo sync: WP STATUS → TODO 상태 매핑 (S_IDL→PENDING, S_PRG→ACTIVE, S_STB→제거)
- [x] 2026-04-08 todo sync: 신규/변경/삭제 WP 자동 반영 + 결과 카운트 출력
- [x] 2026-04-08 todo sync [ADDRESS]: 특정 WP 개별 동기화
- [x] 2026-04-08 todo list: 새 포맷 (ADDRESS/STATUS/SUMMARY) tabwriter 출력
- [x] 2026-04-08 todo list: REFERENCE 대상 미완료 시 [BLOCKED] 자동 표시
- [x] 2026-04-08 todo list: ACTIVE → PENDING → BLOCKED 순 정렬
- [x] 2026-04-08 todo history: WP TECH_SPEC [x] YYYY-MM-DD 항목 직접 수집 (TODO_HISTORY.md 폐지)
- [x] 2026-04-08 todo history [MAP_ADDRESS]: 특정 Map 하위 WP만 필터

### ISSUE
(없음)

### COMMENT
v2.0: 수동 관리(add/done/update/delete) → sync 기반 자동 관리로 전환. WP STATUS가 Single Source of Truth. TODO_HISTORY.md 폐지, TECH_SPEC [x] 항목이 히스토리 역할.
</WAYPOINT>
