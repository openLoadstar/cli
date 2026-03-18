<MAP>
## [ADDRESS] M://root/cli
## [STATUS] S_PRG

### 1. IDENTITY
- SUMMARY: AVCS CLI 명령어 구현 전체를 담는 중간 지도. 각 명령어 구현을 WayPoint로 관리한다.
- METADATA: [Ver: 1.0, Created: 2026-03-04, Priority: HIGH]
- SYNCED_AT: 2026-03-18

### 2. CONTAINS
- ITEMS: [
    W://root/cli/internal_infra,
    W://root/cli/cmd_create,
    W://root/cli/cmd_edit,
    W://root/cli/cmd_delete,
    W://root/cli/cmd_checkpoint,
    W://root/cli/cmd_history,
    W://root/cli/cmd_diff,
    W://root/cli/cmd_rollback,
    W://root/cli/cmd_link,
    W://root/cli/cmd_show,
    W://root/cli/cmd_todo,
    W://root/cli/test_element,
    W://root/cli/test_checkpoint,
    W://root/cli/test_nav,
    W://root/cli/test_todo
  ]
- PAYLOAD: cobra 기반 명령어 구현 WayPoint + 각 명령어 대응 테스트 WayPoint 집합.

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root, CHILDREN: (위 ITEMS 참조)]
- LINKS: []

### 4. RESOURCES
- SAVEPOINTS: []

### 5. TODO
- REQUESTER: NONE
- RESPONSE_STATUS: PENDING
- TECH_SPEC:
- EXECUTION_HISTORY: [
    * 2026-03-18: WayPoint 메타데이터 일괄 동기화 완료 — cmd_checkpoint, cmd_history, cmd_diff, cmd_rollback, cmd_link, cmd_show 6개 WayPoint STATUS S_IDL → S_STB 갱신, BlackBox CODE_MAP 실측 라인번호로 교체 (cmd_checkpoint, cmd_link, cmd_show), SPEC/12.GLOBAL_TODO_LIST.md 더미 데이터 제거
  ]
</MAP>
