# LOADSTAR_INIT — 세션 진입점

> AI 에이전트는 세션 시작 시 이 파일을 읽어 현재 프로젝트 상태를 파악한다.
> 작업 완료 후 갱신한다.

## 마지막 갱신
- DATE: 2026-04-01
- SESSION: loadstar log/findlog 구현, help 텍스트 전체 보강, GIT 디렉토리 통합

---

## 현재 TODO (PENDING)

| Executor | Summary | Depends_On |
|---|---|---|
| W://root/cli/cmd_checkpoint | checkpoint SPEC 미구현 항목 완료 — GIT_INDEX 생성, CHANGE_LOG GIT_REF 역기입, GIT_INDEX 커밋 해시 역기입, HISTORY 임시 스냅샷 정리 | - |
| W://root/cli/cmd_create | appendToContains multiline ITEMS 포맷 파싱 버그 수정 | - |
| W://root/cli/meta_sync | 구현 완료 명령어 WayPoint/BlackBox 신규 생성 — cmd_git, cmd_log, cmd_init | - |
| W://root/cli/cmd_history | 코드 검토 및 SPEC 대조 — 현재 HISTORY 파일 스캔 방식, SPEC은 CHANGE_LOG/GIT_REF 기반 | - |
| W://root/cli/cmd_diff | 코드 검토 및 SPEC 대조 — 현재 HISTORY 스냅샷 직접 비교 방식, SPEC은 GIT_REF → git diff 방식 | - |
| W://root/cli/cmd_rollback | 코드 검토 및 SPEC 대조 — 현재 HISTORY 스냅샷 복원 방식, SPEC은 GIT_REF → git checkout 방식 | - |

> 전체 목록: `loadstar todo list`

---

## 최근 작업 요약 (2026-04-01)

- `cmd/log.go` 신규 구현 — `loadstar log`, `loadstar findlog` 명령 완료
- 전체 명령어 `--help` 예시 텍스트 보강
- `C:\bono\MCP\GIT\loadstar_cli\` 단일 디렉토리로 통합 (AVCS 디렉토리 삭제)
- `bin/loadstar.exe` 빌드 경로 변경 (구 `build/loadstar.exe`, `buildloadstar.exe`)
- `W://root/cli/cmd_checkpoint` WayPoint/BlackBox SPEC 대조 갱신 — STATUS S_STB → S_PRG

---

## 다음 권장 작업

1. **`checkpoint` 미구현 항목** (W://root/cli/cmd_checkpoint)
   - GIT_INDEX 파일 생성 (`GIT.[DATE].[SEQ].[UUID].md`)
   - CHANGE_LOG 파일들에 `GIT_REF` 역기입
   - Commit Hash → GIT_INDEX 역기입 후 추가 커밋
   - HISTORY/ 임시 스냅샷 정리

2. **`create` multiline ITEMS 파싱 버그**
   - `appendToContains()` — multiline 포맷의 ITEMS 파싱 실패
   - 관련 이슈: `loadstar findlog 0 5 --kind ISSUE`

---

## 미구현 SPEC 항목 (저우선순위)

| 명령어 | SPEC 섹션 | 비고 |
|---|---|---|
| `loadstar sync` | §8 | git log 분석 → 누락 GIT_INDEX/CHANGE_LOG 소급 생성 |
| `loadstar monitor snapshot/diff/status/clean` | §14 | 파일 변경 추적 모니터링 |

---

## 프로젝트 구조 참고

```
.loadstar/
├── MAP/          M:// 요소
├── WAYPOINT/     W:// 요소
├── BLACKBOX/     B:// 요소
├── HISTORY/      자동 스냅샷
├── LINK/         L:// 요소
├── SAVEPOINT/    S:// 요소
├── COMMON/       git_config.json 등
└── .clionly/     ← AI 직접 접근 금지
    ├── CHANGE_LOG/
    ├── GIT_INDEX/
    └── TODO/
```
