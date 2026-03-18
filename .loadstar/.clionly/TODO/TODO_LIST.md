# TODO_LIST — avcs-cli

| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 상태 (Status) | 선행 조건 (Depends_On) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| W://root/cli/cmd_rollback | NONE | 2026-03-13 12:38 | rollback 명령어 구현 (History 스냅샷으로 요소 파일 복원) | PENDING | W://root/cli/cmd_diff |
| W://root/cli/cmd_diff | NONE | 2026-03-13 12:38 | diff 명령어 구현 (현재 파일 vs History 스냅샷 unified diff 출력) | PENDING | W://root/cli/cmd_history |
| W://root/cli/cmd_show | NONE | 2026-03-13 12:37 | show 명령어 구현 (요소 출력 + depth N 트리 재귀) | PENDING | - |
| W://root/cli/cmd_link | NONE | 2026-03-13 12:37 | link 명령어 구현 (Link md 생성 + 양방향 CONNECTIONS 등록) | PENDING | - |
| W://root/cli/cmd_history | NONE | 2026-03-13 12:37 | history 명령어 구현 (HISTORY/ 스냅샷 목록 역순 출력) | PENDING | - |
| W://root/cli/cmd_checkpoint | NONE | 2026-03-13 12:37 | checkpoint -m 명령어 구현 (git commit + SavePoint 해시 기입 원자적 처리) | PENDING | - |

---
### 운영 지침 (Operational Rules)
1. **최소 정보 원칙**: 상세 구현 사항은 각 `W://` WayPoint 파일의 TECH_SPEC 및 OPEN_QUESTIONS 참조.
2. **완료 시 처리**: 완료 항목은 본 리스트에서 즉시 삭제 후 `TODO_HISTORY.md`에 이관하고 Executor의 EXECUTION_HISTORY에도 기록.
3. **AI 우선 순위**: 세션 시작 시 이 파일을 읽어 PENDING 작업을 식별하고 해당 WayPoint로 진입.
4. **Depends_On 규칙**: `Depends_On` 항목이 완료되지 않은 작업은 `[BLOCKED]` 상태로 간주하며 착수하지 않는다. `-`는 선행 조건 없음을 의미한다.
5. **OPEN_QUESTIONS 우선**: WayPoint 진입 후 OPEN_QUESTIONS에 미해결 항목(`[QN]`)이 있으면 코드 작성 전 사람에게 먼저 확인한다.
