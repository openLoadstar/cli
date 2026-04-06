<WAYPOINT>
## [ADDRESS] W://root/cli/test_todo
## [STATUS] S_STB

### IDENTITY
- SUMMARY: todo add/done/list/update 명령어 통합 테스트. GLOBAL_TODO_LIST.md 마크다운 테이블 파싱·추가·삭제·상태 변경·BLOCKED 표시를 검증한다.
- METADATA: [Ver: 1.1, Created: 2026-03-10, Priority: HIGH]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/test_todo

### TODO
- [x] `TestTodoAdd_NewRow`: todo add 실행 후 GLOBAL_TODO_LIST에 행 추가 확인
- [x] `TestTodoAdd_WithDepends`: --depends 플래그로 선행 조건 등록 확인
- [x] `TestTodoDone_RemovesRow`: todo done 실행 후 해당 행 삭제 확인
- [x] `TestTodoDone_ExecutionHistory`: done 후 Executor 요소의 EXECUTION_HISTORY에 기록 추가 확인
- [x] `TestTodoDone_UnblocksDependent`: done 후 해당 항목을 Depends_On으로 참조하던 행 상태 재평가 확인
- [x] `TestTodoList_BlockedDisplay`: Depends_On 미완료 항목에 [BLOCKED] 표시 확인
- [x] `TestTodoList_Empty`: 항목 없을 때 안내 메시지 확인
- [x] `TestTodoUpdate_StatusChange`: todo update로 PENDING → ACTIVE 상태 변경 확인
- [x] `TestTodoUpdate_PendingToBlocked`: todo update로 PENDING → BLOCKED 상태 변경 확인
- [x] `TestTodoUpdate_InvalidStatus`: COMPLETED/FAILED 등 허용 외 상태값 입력 시 에러 확인
- [x] `TestTodoUpdate_AllowedStatuses`: PENDING/ACTIVE/BLOCKED 허용 확인
- [x] `TestTodoUpdate_NotFound`: 존재하지 않는 executor 입력 시 에러 확인
- [x] `TestTodoUpdate_PreservesOtherCols`: 상태 변경 후 나머지 컬럼 보존 확인

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
