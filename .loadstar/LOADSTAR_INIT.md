# LOADSTAR_INIT — 세션 진입점

> AI 에이전트는 세션 시작 시 이 파일을 읽어 현재 프로젝트 상태를 파악한다.
> 작업 완료 후 갱신한다.

## 마지막 갱신
- DATE: 2026-04-02
- SESSION: CLI 프로토타입 완료 + 테스트 프로젝트 검증 + log 명령 개선 (ID 단축, --list, 페이징)

---

## 현재 TODO (PENDING)

| Address | Summary | Depends_On |
|---|---|---|
| W://root/cli/cmd_create | appendToContains → appendToWaypoints/appendToChildren 리팩토링 후 multiline 파싱 버그 재검증 | - |

> 전체 목록: `loadstar todo list`

---

## 최근 작업 요약 (2026-04-02)

### 구조 대규모 리팩토링
- **요소 포맷 단순화**:
  - Map: IDENTITY + WAYPOINTS + COMMENT (인덱스만)
  - WayPoint: IDENTITY + CONNECTIONS(PARENT/CHILDREN/REFERENCE/BLACKBOX) + TODO + ISSUE + COMMENT
  - BlackBox: DESCRIPTION + CODE_MAP + TODO + ISSUE + COMMENT
- **삭제된 디렉토리**: LINK/, SAVEPOINT/, HISTORY/
- **삭제된 명령**: history, diff, rollback, link (git 직접 활용으로 대체)
- **이름 변경**: .clionly/CHANGE_LOG → .clionly/LOG
- **주소 체계**: M://, W://, B:// 3종류만 유지 (L://, S://, H:// 제거)
- **todo 단순화**: REQUESTER/EXECUTOR → ADDRESS 단일 필드

### checkpoint 단순화 (이전 세션에서 시작)
- GIT_INDEX 폐지, 커밋 메시지에 변경 요소 자동 나열
- `--auto` 플래그, `.clionly/MONITOR/checkpoint_needed.flag` 연동
- `loadstar_monitor.exe` 별도 프로세스 신규 구현

### Go 코드 변경 파일
- `cmd/element.go` — 템플릿 전면 교체, appendToWaypoints/appendToChildren/parseParent 신규
- `cmd/checkpoint.go` — history/diff/rollback 삭제, checkpoint만 유지
- `cmd/nav.go` — link 삭제, show만 유지 (extractChildren 리팩토링)
- `cmd/todo.go` — REQUESTER/EXECUTOR 제거, 2-arg add
- `cmd/log.go` — CHANGE_LOG→LOG, LOG→COMMENT
- `cmd/root.go` — 4개 명령 등록 제거
- `internal/address/address.go` — L, S, H 타입 제거
- `internal/storage/fs.go` — LINK/SAVEPOINT/HISTORY 디렉토리 제거

### 메타데이터 마이그레이션
- MAP 2개, WAYPOINT 15개, BLACKBOX 7개 전체 새 포맷 적용
- .loadstar/LINK/, SAVEPOINT/, HISTORY/ 디렉토리 삭제

---

## 다음 권장 작업

1. **UI 설계 및 구현** — WayPoint 흐름 시각화, 관리자 편집 기능, CLI 연동 정보 표시
   - 기술 스택: 미정 (JavaFX, Web 등 논의 필요)
   - 핵심: 프로젝트 방향성 확인 + 수정 + 진행 상황 모니터링
2. **테스트 프로젝트** — `C:\bono\MCP\GIT\test_calc\` (Java/JavaFX 계산기, CLI 전체 워크플로우 검증 완료)

---

## 프로젝트 구조 참고

```
loadstar_cli/
├── main.go                    # loadstar.exe 엔트리포인트
├── cmd/monitor/main.go        # loadstar_monitor.exe 엔트리포인트
├── bin/
│   ├── loadstar.exe
│   └── loadstar_monitor.exe
└── .loadstar/
    ├── MAP/
    ├── WAYPOINT/
    ├── BLACKBOX/
    ├── COMMON/
    └── .clionly/
        ├── LOG/
        ├── MONITOR/
        └── TODO/
```
