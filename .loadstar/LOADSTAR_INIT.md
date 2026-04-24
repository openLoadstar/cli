# LOADSTAR_INIT — loadstar_cli

> AI 세션 진입 시 이 파일을 읽어 프로젝트 컨텍스트를 복원합니다.

## 프로젝트 개요

- **언어**: Go + cobra CLI
- **바이너리**: `bin/loadstar.exe`
- **명령어**: `init` · `show` · `todo (sync/list/history)` · `log` · `validate` · `check`

## AI 참고사항

- CLI 핵심 기능 구현 완료 (2026-04-24 기준 모든 기능 WP가 `S_STB`)
- 유지보수 작업은 `M://root/maintenance` 산하 WP로 관리
- loadstar_ui 프로젝트에서 이 CLI를 백엔드 브릿지로 호출

## 최근 변경

- 2026-04-24 `internal_infra` CI 워크플로우 항목 제거 → `S_STB` 전환
- 2026-04-24 `M://root/maintenance` 맵 신설 + `fix_claude_md` WP 완료
- 2026-04-24 SPEC에서 CodeBrief 개념 제거 (SPEC v3 관련)
