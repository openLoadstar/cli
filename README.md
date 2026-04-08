# LOADSTAR CLI

Go 기반 프로젝트 메타데이터 관리 도구. WayPoint/Map 구조를 통해 AI 에이전트와 개발자가 프로젝트 작업을 추적하고 관리합니다.

## 명령어

| 명령 | 용도 |
|:---|:---|
| `loadstar init` | `.loadstar/` 디렉토리 구조 초기화 |
| `loadstar show [FILTER]` | 전체 WayPoint 목록 출력, 키워드 필터 |
| `loadstar validate` | 모든 요소의 참조 링크 검증, 깨진 참조 보고 |
| `loadstar log [ADDR] [KIND] "[MSG]"` | 메타 이벤트 로그 기록 (6종 KIND) |
| `loadstar findlog [OFFSET] [LIMIT]` | 로그 검색 (최신순, 주소/KIND 필터, 페이징) |
| `loadstar todo sync` | WayPoint STATUS 기반 TODO 목록 자동 동기화 |
| `loadstar todo list` | 현재 PENDING/ACTIVE/BLOCKED 작업 목록 출력 |
| `loadstar todo history [MAP_ADDR]` | WayPoint TECH_SPEC 완료 항목 수집 |

## 빠른 시작

```bash
# 빌드
go build -o bin/loadstar.exe .

# 새 프로젝트 초기화
cd /my/project
loadstar init

# 전체 WayPoint 상태 확인
loadstar show

# 특정 키워드로 필터
loadstar show frontend

# 깨진 참조 검증
loadstar validate

# TODO 동기화 및 확인
loadstar todo sync
loadstar todo list

# 완료 히스토리 조회
loadstar todo history
loadstar todo history M://root/cli

# 메타 이벤트 로그
loadstar log cmd_show MODIFIED "show 명령 개편"
loadstar findlog 0 10
loadstar findlog 0 5 --kind ISSUE
```

## 주소 체계

```
M://root/cli            →  .loadstar/MAP/root.cli.md
W://root/cli/cmd_show   →  .loadstar/WAYPOINT/root.cli.cmd_show.md
```

- **M** (Map): WayPoint를 묶는 구조적 그룹
- **W** (WayPoint): 작업의 최소 단위

## 디렉토리 구조

```
.loadstar/
├── MAP/          M:// 요소
├── WAYPOINT/     W:// 요소
├── COMMON/       설정 파일
└── .clionly/     ← CLI 전용 (직접 편집 금지)
    ├── LOG/       변경 이력 로그
    └── TODO/      작업 목록 (sync가 관리)
```

## AI 협업 워크플로우

1. 코드 수정 전: 대상 WayPoint TECH_SPEC에 `[ ]` 항목 등록
2. 코드 수정 완료 후: `[x] YYYY-MM-DD` 체크, 필요 시 STATUS 갱신
3. `loadstar todo sync` → TODO_LIST 자동 동기화
4. Claude Code 사용 시: PostToolUse Hook이 TECH_SPEC 등록/갱신 리마인더 출력

## 의존성

- [cobra](https://github.com/spf13/cobra) — CLI 프레임워크

## 관련 프로젝트

- [loadstar_SPEC](https://github.com/aeolusk/loadstar_SPEC) — LOADSTAR 방법론 명세
- [loadstar_ui](https://github.com/aeolusk/loadstar_ui) — 웹 기반 Explorer UI
