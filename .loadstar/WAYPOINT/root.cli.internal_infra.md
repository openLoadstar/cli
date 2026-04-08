<WAYPOINT>
## [ADDRESS] W://root/cli/internal_infra
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `internal/` 패키지군 완성 작업. address 파싱, storage I/O, git 클라이언트, core ElementService 등 모든 명령어가 공유하는 인프라 레이어 구현.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: CRITICAL]

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
- [x] `cmd/root.go` PersistentPreRun에서 fs, svc, gitClient 공통 초기화

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
