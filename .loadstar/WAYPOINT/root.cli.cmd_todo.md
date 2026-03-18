<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_todo
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `loadstar todo add/done/list/update/delete/history` 서브커맨드 구현. .clionly/TODO/TODO_LIST.md 테이블을 파싱·수정하고, done 시 .clionly/TODO/TODO_HISTORY.md 이관 및 Executor의 EXECUTION_HISTORY에 기록한다. update/delete 이벤트도 TODO_HISTORY에 통합 기록하며, history 명령으로 이력을 조회할 수 있다.
- METADATA: [Ver: 1.5, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-03-16

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/todo.go` (todoAddCmd, todoDoneCmd, todoListCmd, todoUpdateCmd, todoDeleteCmd, todoHistoryCmd)
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`, `internal/core/element.go`
  - 저장 경로:
    - 현재 목록: `.loadstar/.clionly/TODO/TODO_LIST.md`
    - 완료 이력: `.loadstar/.clionly/TODO/TODO_HISTORY.md`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/infra_to_todo | TYPE: L_REF,
    L://root/cli/todo_to_test_todo | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []
- BLACKBOX_REF: B://root/cli/cmd_todo

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/cmd_todo
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
  - [x] TODO_LIST 경로 확인: `.loadstar/.clionly/TODO/TODO_LIST.md`
  - [x] 마크다운 테이블 파서 구현: 헤더·구분선 이후 데이터 행을 `|` 기준으로 파싱
  - [x] `todo add`: 신규 행 추가 (PENDING, 현재 시각 자동 기입, `--depends` 플래그 처리)
  - [x] `todo done`: 행 삭제 → TODO_HISTORY.md append → Executor EXECUTION_HISTORY 기록
  - [x] `todo done` 부수 효과: 완료 항목을 Depends_On으로 참조하는 다른 행의 상태 재평가
  - [x] `todo list`: PENDING/ACTIVE 행 필터링 출력, Depends_On 미완료 항목에 `[BLOCKED]` 표시
  - [x] 시각 포맷: `2006-01-02 15:04` (Go time 레이아웃)
  - [x] `todo update [EXECUTOR] [STATUS]`: 행의 상태 컬럼을 PENDING/ACTIVE/BLOCKED 중 하나로 직접 변경
  - [x] `todo update` 시 변경 이력을 TODO_HISTORY에 `UPDATED(OLD→NEW)` 액션으로 기록
  - [x] `todo delete [EXECUTOR]`: TODO_LIST에서 행 삭제 후 TODO_HISTORY에 `DELETED` 액션으로 기록
  - [x] TODO_HISTORY 포맷에 Action/At 컬럼 추가 (done/update/delete 모든 이벤트 통합 기록)
  - [x] `todo history`: TODO_HISTORY 전체 출력
  - [x] `todo history [EXECUTOR]`: 특정 executor 필터링 출력 (완전 일치)
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] `todo done`의 Executor 매칭은 완전 일치만 허용한다.
  - [Q2 RESOLVED] `todo update`의 허용 상태값은 PENDING, ACTIVE, BLOCKED 세 가지로 제한. COMPLETED/FAILED는 `todo done` 경로로만 처리한다.
- EXECUTION_HISTORY: [
    * 2026-03-12: add/done/list 구현 완료로 기존 체크박스 갱신, update 서브커맨드 추가 착수
    * 2026-03-13: TODO/ 폴더 구조로 변경 (TODO_LIST.md + TODO_HISTORY.md), done 시 이중 이력 기록 추가
    * 2026-03-13: TODO/ 위치를 .loadstar/.clionly/TODO/로 이동 (CLI 전용 데이터 분리 정책 적용)
    * 2026-03-16: TODO_HISTORY 포맷에 Action/At 컬럼 추가, todo delete 서브커맨드 구현, update 시 이력 기록 추가 — 전체 테스트 26개 PASS
    * 2026-03-16: todo history / todo history [EXECUTOR] 서브커맨드 구현 — 전체 테스트 30개 PASS
  ]
</WAYPOINT>
