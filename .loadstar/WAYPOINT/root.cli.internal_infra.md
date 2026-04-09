<WAYPOINT>
## [ADDRESS] W://root/cli/internal_infra
## [STATUS] S_PRG

### IDENTITY
- SUMMARY: `internal/` 패키지군. address 파싱, storage I/O, core ElementService. Init 시 .loadstar/ + .claude/ hooks 자동 생성.
- METADATA: [Ver: 1.2, Created: 2026-03-04, Priority: CRITICAL]
- SYNCED_AT: 2026-04-08

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
- [x] Storage, GitClient를 interface로 정의 (`internal/interfaces.go`)
- [x] `git.Client.Commit()`: go-git으로 `.avcs/` 스테이징 및 커밋 구현
- [x] `git.Client.LatestHash()`: HEAD 해시 반환 구현
- [x] `ElementService`에 Storage, AddressParser 인터페이스 주입
- [x] 프로젝트 루트 탐색 로직: `storage.FindRoot()` — `.avcs/` 상위 탐색, 미존재 시 현재 디렉토리에 자동 초기화
- [x] `go mod tidy`로 go-git 의존성 및 go.sum 정리
- [x] `cmd/root.go` PersistentPreRun에서 fs, svc 공통 초기화
- [x] 2026-04-08 git 패키지 제거 (GitClient 인터페이스 폐지, git 직접 사용)
- [x] 2026-04-08 Init()에 .claude/hooks 자동 생성 (settings.json + loadstar-drift-check.sh)
- [ ] GitHub Actions CI 워크플로우 설정 (.github/workflows/ — go test + 리포트)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
