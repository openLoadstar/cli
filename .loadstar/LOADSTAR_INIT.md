# LOADSTAR_INIT — 세션 진입점

> AI 에이전트는 세션 시작 시 이 파일을 읽어 현재 프로젝트 상태를 파악한다.
> 작업 완료 후 갱신한다.

## 마지막 갱신
- DATE: 2026-04-08
- SESSION: CLI 대규모 리팩토링 — 명령어 정리, BlackBox 제거, Hook 기반 메타 동기화

---

## 현재 TODO (PENDING)

> 전체 목록: `loadstar todo list`

---

## 최근 작업 요약 (2026-04-08)

### CLI 리팩토링 (2차)
- **제거된 명령어**: create, edit, delete, checkpoint, git (set/status/unset)
- **제거된 프로세스**: loadstar_monitor.exe
- **BlackBox 완전 제거**: B:// 주소 체계, .loadstar/BLACKBOX/ 디렉토리
- **show 개편**: `loadstar show [FILTER]` — 전체 WayPoint 목록 + 키워드 필터
- **validate 추가**: `loadstar validate` — 깨진 링크 검출
- **Hook 추가**: `.claude/hooks/loadstar-drift-check.sh` — 소스 수정 시 메타 갱신 리마인더
- **init 단순화**: .loadstar/ 구조만 생성 (git 초기화 제거)
- **주소 체계**: M://, W:// 2종류만 유지

### 현재 명령어 (6개)
`init` · `show` · `todo` · `log` · `findlog` · `validate`

---

## 다음 권장 작업

1. **UI 프로젝트 (loadstar_ui)** — 진행 중: search, dashboard, monitor_view
2. **CLI 테스트 보강** — test_nav (S_IDL)

---

## 프로젝트 구조 참고

```
loadstar_cli/
├── main.go              # loadstar.exe 엔트리포인트
├── bin/loadstar.exe
├── .claude/
│   ├── settings.json    # PostToolUse hook 설정
│   └── hooks/           # 메타 동기화 리마인더 스크립트
└── .loadstar/
    ├── MAP/
    ├── WAYPOINT/
    ├── COMMON/
    └── .clionly/
        ├── LOG/
        ├── MONITOR/
        └── TODO/
```
