# TODO_HISTORY — avcs-cli

| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 액션 (Action) | 처리 시각 (At) | 선행 조건 (Depends_On) |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| W://root/cli/cmd_todo | NONE | 2026-03-13 13:09 | todo done/update/delete 시 TODO_HISTORY 이관 및 TODO_LIST 정리 강화: 상태 변경·완료·삭제 모든 경우에 이력 보존 | DONE | 2026-03-16 12:41 | - |
| W://root/cli/cmd_checkpoint | NONE | 2026-03-13 12:37 | checkpoint -m 명령어 구현 (git commit + SavePoint 해시 기입 원자적 처리) | UPDATED(PENDING→ACTIVE) | 2026-03-16 12:45 | - |
| W://root/cli/cmd_checkpoint | NONE | 2026-03-13 12:37 | checkpoint -m 명령어 구현 (git commit + SavePoint 해시 기입 원자적 처리) | UPDATED(ACTIVE→PENDING) | 2026-03-16 12:45 | - |
