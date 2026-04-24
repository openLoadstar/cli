# LOADSTAR CLI — Claude Agent 운영 규칙

## 세션 시작 절차 (매 세션 필수)

1. 이 파일을 읽는다.
2. `.loadstar/LOADSTAR_INIT.md` 를 읽어 현재 프로젝트 상태를 파악한다.
3. `loadstar check` 를 실행하여 WP 동기화 상태를 확인한다 (git commit 대비 WP 수정시간 drift 검출).
4. 사용자에게 아래 질문을 한다:

> **LOADSTAR SPEC 파일을 로드할까요?**
> - 새로운 기능 구현, 설계 변경, SPEC 참조가 필요한 작업이면 **권장**
> - 버그 수정, 단순 코드 수정이면 **불필요**

5. **Yes** → `C:\bono\MCP\GIT\loadstar_SPEC\` 에서 관련 파일 로드
6. **No** → `LOADSTAR_INIT.md` 내용만으로 진행

---

## 프로젝트 개요

- **언어/프레임워크**: Go + cobra CLI
- **저장소**: `C:\bono\MCP\GIT\loadstar_cli\`
- **바이너리**: `bin/loadstar.exe`
- **빌드**: `go build -o bin/loadstar.exe .`
- **SPEC 문서**: `C:\bono\MCP\GIT\loadstar_SPEC\`
- **LOADSTAR 메타데이터**: `.loadstar\`

---

## LOADSTAR 핵심 규칙

### 절대 규칙 (위반 금지)
- **`.loadstar/.clionly/` 직접 읽기·쓰기 금지** — CLI 전담 영역
- MD 직접 편집 후 반드시 `loadstar log [ADDR] MODIFIED "[내용]"` 실행

### 주소 체계 (2종류)
```
M://root/cli            →  .loadstar/MAP/root.cli.md
W://root/cli/cmd_show   →  .loadstar/WAYPOINT/root.cli.cmd_show.md
```

### 요소 포맷
- **Map**: IDENTITY, WAYPOINTS, COMMENT (인덱스 역할만)
- **WayPoint**: IDENTITY, CONNECTIONS (PARENT/CHILDREN/REFERENCE), TODO, ISSUE, COMMENT

### 작업 착수 규칙 (필수)
1. **코드 수정 전**: 대상 WayPoint의 TECH_SPEC에 작업 항목을 `[ ]`로 추가
2. **코드 수정 완료 후**: `[x] YYYY-MM-DD 항목명`으로 체크
3. **WP 전체 완료 시**: STATUS를 `S_PRG → S_STB`로 변경
4. 항목 추가 없이 코드 수정 착수 금지 — Hook이 리마인드함

### 작업 진입 순서
| 작업 유형 | 진입 순서 |
|---|---|
| 기능 구현 / 설계 변경 / 영향 범위 불명확 | MAP → WayPoint TECH_SPEC 항목 등록 → 코드 |
| 명확한 버그 수정 / 단일 함수 수정 | grep → 코드 → WayPoint TECH_SPEC 사후 등록 + 체크 |

### 메타 동기화
- **Hook**: `.claude/hooks/loadstar-drift-check.sh`가 소스코드 편집 시 TECH_SPEC 등록/갱신 리마인더 출력
- **todo sync**: `loadstar todo sync`로 WP STATUS 기반 TODO_LIST 자동 동기화

---

## 구현 완료 명령어

`init` · `show` · `todo (sync/list/history)` · `log` · `validate` · `check`

- **`loadstar check`**: git 마지막 커밋 시점 대비 WP 파일 수정시간을 비교하여 OUTDATED WP를 표시합니다.
- **`loadstar log [TIME] [FILTER]`**: 기간·키워드 필터 로그 조회 (이전 `findlog`는 이 명령으로 통합됨).

## 삭제된 명령어

`todo add/update/done/delete` — sync가 WP STATUS 기반 자동 관리 (2026-04-08)
`create` · `edit` · `delete` — AI 직접 편집 + UI로 대체 (2026-04-08)
`checkpoint` · `git (set/status/unset)` — git 직접 사용으로 대체 (2026-04-08)
`history` · `diff` · `rollback` · `link` — git 직접 활용으로 대체 (2026-04-02)
`findlog` — `log [TIME_RANGE] [FILTER]`에 통합 (2026-04-10)

---

## 디렉토리 구조

```
.loadstar/
├── MAP/          M:// 요소
├── WAYPOINT/     W:// 요소
├── COMMON/       설정 파일
└── .clionly/     ← AI 직접 접근 금지
    ├── LOG/       변경 이력 로그
    ├── MONITOR/   (예약)
    └── TODO/      작업 목록 관리
```
