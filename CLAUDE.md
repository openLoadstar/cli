# LOADSTAR CLI — Claude Agent 운영 규칙

## 세션 시작 절차 (매 세션 필수)

1. 이 파일을 읽는다.
2. `.loadstar/LOADSTAR_INIT.md` 를 읽어 현재 프로젝트 상태를 파악한다.
3. 사용자에게 아래 질문을 한다:

> **LOADSTAR SPEC 파일을 로드할까요?**
> - 새로운 기능 구현, 설계 변경, SPEC 참조가 필요한 작업이면 **권장**
> - 버그 수정, 단순 코드 수정이면 **불필요**

4. **Yes** → `C:\bono\MCP\GIT\SPEC\` 에서 관련 파일 로드
5. **No** → `LOADSTAR_INIT.md` 내용만으로 진행

---

## 프로젝트 개요

- **언어/프레임워크**: Go + cobra CLI
- **저장소**: `C:\bono\MCP\GIT\loadstar_cli\`
- **바이너리**: `bin/loadstar.exe`, `bin/loadstar_monitor.exe`
- **빌드**: `go build -o bin/loadstar.exe .` / `go build -o bin/loadstar_monitor.exe ./cmd/monitor`
- **SPEC 문서**: `C:\bono\MCP\GIT\SPEC\` (01~14번 파일)
- **LOADSTAR 메타데이터**: `.loadstar\`

---

## LOADSTAR 핵심 규칙

### 절대 규칙 (위반 금지)
- **`.loadstar/.clionly/` 직접 읽기·쓰기 금지** — CLI 전담 영역
- MD 직접 편집 후 반드시 `loadstar log [ADDR] MODIFIED "[내용]"` 실행

### 주소 체계 (3종류만)
```
M://root/cli            →  .loadstar/MAP/root.cli.md
W://root/cli/cmd_create  →  .loadstar/WAYPOINT/root.cli.cmd_create.md
B://root/cli/cmd_create  →  .loadstar/BLACKBOX/root.cli.cmd_create.md
```

### 요소 포맷
- **Map**: IDENTITY, WAYPOINTS, COMMENT (인덱스 역할만)
- **WayPoint**: IDENTITY, CONNECTIONS (PARENT/CHILDREN/REFERENCE/BLACKBOX), TODO, ISSUE, COMMENT
- **BlackBox**: DESCRIPTION, CODE_MAP, TODO, ISSUE, COMMENT

### 작업 진입 순서
| 작업 유형 | 진입 순서 |
|---|---|
| 기능 구현 / 설계 변경 / 영향 범위 불명확 | WayPoint → BlackBox → 코드 |
| 명확한 버그 수정 / 단일 함수 수정 | grep → 코드 → WayPoint·BlackBox 사후 갱신 |

### 작업 후 처리
- 코드 수정 완료 → `loadstar checkpoint -m "..."` (git commit + push 자동)
- MD 직접 편집 → `loadstar log [ADDR] MODIFIED "[변경 내용]"`

---

## 구현 완료 명령어

`init` · `create` · `edit` · `delete` · `checkpoint` · `show` · `todo (add/list/update/done/delete/history)` · `git (set/status/unset)` · `log` · `findlog`

## 삭제된 명령어

`history` · `diff` · `rollback` · `link` — git 직접 활용으로 대체 (2026-04-02)

## 별도 프로세스

`loadstar_monitor.exe` — 5분 간격 .loadstar/ 변경 감시, checkpoint 필요 시 flag 파일 생성

---

## 디렉토리 구조

```
.loadstar/
├── MAP/          M:// 요소
├── WAYPOINT/     W:// 요소
├── BLACKBOX/     B:// 요소
├── COMMON/       git_config.json 등
└── .clionly/     ← AI 직접 접근 금지
    ├── LOG/       변경 이력 로그
    ├── MONITOR/   모니터 프로세스 상태
    └── TODO/      작업 목록 관리
```
