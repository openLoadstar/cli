<BLACKBOX>
## [ADDRESS] B://root/cli/cmd_todo
## [STATUS] S_STB

### DESCRIPTION
- SUMMARY: `avcs todo add/done/list/update/delete/history` 구현. .clionly/TODO/TODO_LIST.md 테이블 파싱·수정, done/update/delete 시 .clionly/TODO/TODO_HISTORY.md 통합 이력 기록, EXECUTION_HISTORY 동시 기록(done만), history 명령으로 이력 조회.
- LINKED_WP: W://root/cli/cmd_todo

### CODE_MAP

**구현 후 (실측)**
- `cmd/todo.go:14-15`
  - 경로 상수: `todoListFile = ".clionly/TODO/TODO_LIST.md"`, `todoHistoryFile = ".clionly/TODO/TODO_HISTORY.md"`
- `cmd/todo.go:25-69`
  - `todoAddCmd.Run()` → ensureTodoFile() → 헤더·구분선 이후 신규 행 insert
- `cmd/todo.go:71-127`
  - `todoDoneCmd.Run()` → 행 탐색/저장 → appendTodoHistory() → appendExecutionHistory()
- `cmd/todo.go:155-195`
  - `todoUpdateCmd.Run()` → 상태 변경 후 appendTodoHistory()로 UPDATED(OLD→NEW) 이력 기록
- `cmd/todo.go:198-236`
  - `todoDeleteCmd.Run()` → 행 삭제 후 appendTodoHistory()로 DELETED 이력 기록
- `cmd/todo.go:258-305`
  - `appendTodoHistory(histPath, row, action)` → Action/At 컬럼 포함 TODO_HISTORY.md append
- `cmd/todo.go:237-271`
  - `todoListCmd.Run()` → PENDING/ACTIVE 필터, Depends_On 미완료 시 [BLOCKED] 표시
- `cmd/todo.go:273-318`
  - `todoHistoryCmd.Run()` → TODO_HISTORY.md 읽기, executor 인자 있으면 완전 일치 필터, 없으면 전체 출력

### TODO
- [x] TODO_LIST 경로 확인 및 초기화 [WP_REF:1]
- [x] 마크다운 테이블 파서 구현 [WP_REF:2]
- [x] `todo add`: 신규 행 추가, --depends 플래그 처리 [WP_REF:3]
- [x] `todo done`: 행 삭제 → TODO_HISTORY append → EXECUTION_HISTORY 이관 [WP_REF:4]
- [x] `todo done` 부수 효과: Depends_On 재평가 [WP_REF:5]
- [x] `todo list`: 필터링 출력 및 BLOCKED 표시 [WP_REF:6]
- [x] `todo history`: TODO_HISTORY 전체 또는 executor 필터 출력 [WP_REF:7]

### ISSUE
- 해결됨: Executor 매칭 완전 일치 전용 (Q1 RESOLVED)

### COMMENT
(없음)
</BLACKBOX>
