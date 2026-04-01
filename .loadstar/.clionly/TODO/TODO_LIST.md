# TODO_LIST — avcs-cli

| 실행 요소 (Executor) | 요청 요소 (Requester) | 발생 시간 (Time) | 작업 요약 (Summary) | 상태 (Status) | 선행 조건 (Depends_On) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| W://root/cli/cmd_rollback | NONE | 2026-04-01 17:45 | 코드 검토 및 SPEC 대조 — 현재 HISTORY 스냅샷 복원 방식, SPEC은 GIT_REF → git checkout 방식 | PENDING | - |
| W://root/cli/cmd_diff | NONE | 2026-04-01 17:45 | 코드 검토 및 SPEC 대조 — 현재 HISTORY 스냅샷 직접 비교 방식, SPEC은 GIT_REF → git diff 방식 | PENDING | - |
| W://root/cli/cmd_history | NONE | 2026-04-01 17:45 | 코드 검토 및 SPEC 대조 — 현재 HISTORY 파일 스캔 방식, SPEC은 CHANGE_LOG/GIT_REF 기반 | PENDING | - |
| W://root/cli/meta_sync | NONE | 2026-04-01 17:44 | 구현 완료 명령어 WayPoint/BlackBox 신규 생성 — cmd_git, cmd_log, cmd_init | PENDING | - |
| W://root/cli/cmd_create | NONE | 2026-04-01 17:44 | appendToContains multiline ITEMS 포맷 파싱 버그 수정 | PENDING | - |
| W://root/cli/cmd_checkpoint | NONE | 2026-04-01 17:36 | checkpoint SPEC 미구현 항목 완료 — GIT_INDEX 생성, CHANGE_LOG GIT_REF 역기입, GIT_INDEX 커밋 해시 역기입, HISTORY 임시 스냅샷 정리 | PENDING | - |

---
### 운영 지침 (Operational Rules)
1. **최소 정보 원칙**: 상세 구현 사항은 각 `W://` WayPoint 파일의 TECH_SPEC 및 OPEN_QUESTIONS 참조.
2. **완료 시 처리**: 완료 항목은 본 리스트에서 즉시 삭제 후 `TODO_HISTORY.md`에 이관하고 Executor의 EXECUTION_HISTORY에도 기록.
3. **AI 우선 순위**: 세션 시작 시 이 파일을 읽어 PENDING 작업을 식별하고 해당 WayPoint로 진입.
4. **Depends_On 규칙**: `Depends_On` 항목이 완료되지 않은 작업은 `[BLOCKED]` 상태로 간주하며 착수하지 않는다. `-`는 선행 조건 없음을 의미한다.
5. **OPEN_QUESTIONS 우선**: WayPoint 진입 후 OPEN_QUESTIONS에 미해결 항목(`[QN]`)이 있으면 코드 작성 전 사람에게 먼저 확인한다.
