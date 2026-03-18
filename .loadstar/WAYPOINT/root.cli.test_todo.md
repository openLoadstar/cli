<WAYPOINT>
## [ADDRESS] W://root/cli/test_todo
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: todo add/done/list/update 명령어 통합 테스트. GLOBAL_TODO_LIST.md 마크다운 테이블 파싱·추가·삭제·상태 변경·BLOCKED 표시를 검증한다.
- METADATA: [Ver: 1.1, Created: 2026-03-10, Priority: HIGH]
- SYNCED_AT: 2026-03-13

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일: `cmd/todo_test.go`
  - 의존 패키지: `internal/address/address.go`, `internal/storage/fs.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/todo_to_test_todo | TYPE: L_TST
  ]

### 4. RESOURCES
- SAVEPOINTS: []

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/test_todo
- RESPONSE_STATUS: COMPLETED
- TECH_SPEC:
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
- OPEN_QUESTIONS: []
- EXECUTION_HISTORY: [
    * 2026-03-11: [COMPLETED] 전체 테스트 55개 PASS (go test ./cmd/... -v)
    * 2026-03-12: update 서브커맨드 테스트 6개 추가, 전체 17개 PASS (go test ./cmd/... -run Todo)
    * 2026-03-13: 모든 체크항목 확인 완료, STATUS S_PRG → S_STB
  ]
</WAYPOINT>
