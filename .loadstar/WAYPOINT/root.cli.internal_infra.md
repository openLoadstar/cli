<WAYPOINT>
## [ADDRESS] W://root/cli/internal_infra
## [STATUS] S_STB

### 1. IDENTITY
- SUMMARY: `internal/` 패키지군 완성 작업. address 파싱, storage I/O, git 클라이언트, core ElementService 등 모든 명령어가 공유하는 인프라 레이어 구현.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: CRITICAL]
- SYNCED_AT: 2026-03-10

### 2. CONTAINS
- ITEMS: []
- PAYLOAD:
  - 대상 파일:
    - `internal/interfaces.go`
    - `internal/address/address.go`
    - `internal/storage/fs.go`
    - `internal/git/client.go`
    - `internal/core/element.go`
    - `cmd/root.go`

### 3. CONNECTIONS
- LINEAGE: [PARENT: M://root/cli, CHILDREN: []]
- LINKS: [
    L://root/cli/infra_to_checkpoint | TYPE: L_REF,
    L://root/cli/infra_to_todo | TYPE: L_REF
  ]

### 4. RESOURCES
- SAVEPOINTS: []

### 5. TODO
- REQUESTER: M://root/cli
- EXECUTOR: W://root/cli/internal_infra
- RESPONSE_STATUS: COMPLETED
- TECH_SPEC:
  - [x] Storage, GitClient를 interface로 정의 (`internal/interfaces.go`)
  - [x] `git.Client.Commit()`: go-git으로 `.avcs/` 스테이징 및 커밋 구현
  - [x] `git.Client.LatestHash()`: HEAD 해시 반환 구현
  - [x] `ElementService`에 Storage, AddressParser 인터페이스 주입
  - [x] 프로젝트 루트 탐색 로직: `storage.FindRoot()` — `.avcs/` 상위 탐색, 미존재 시 현재 디렉토리에 자동 초기화
  - [x] `go mod tidy`로 go-git 의존성 및 go.sum 정리
  - [x] `cmd/root.go` PersistentPreRun에서 fs, svc, gitClient 공통 초기화
- OPEN_QUESTIONS:
  - [Q1 RESOLVED] 별도 `avcs init` 없이, 모든 명령 실행 시 `.avcs/` 미존재 감지 → 현재 디렉토리에 자동 초기화.
  - [Q2 RESOLVED] 루트 탐색 실패 시 현재 디렉토리에 자동으로 `.avcs/` 구조 생성 후 계속 진행.
- EXECUTION_HISTORY: [
    * 2026-03-04: internal/interfaces.go 생성 — Storage, GitClient 인터페이스 정의
    * 2026-03-04: storage/fs.go 재작성 — FS 구조체, FindRoot, CopyFile, ListByPrefix 추가
    * 2026-03-04: git/client.go 구현 — go-git 기반 Commit, LatestHash 완성
    * 2026-03-04: core/element.go 재작성 — ElementService에 Storage 의존성 주입
    * 2026-03-04: cmd/root.go 재작성 — PersistentPreRun으로 fs/svc/gitClient 공통 초기화
    * 2026-03-04: go mod tidy 완료 — 빌드 성공 확인
  ]
</WAYPOINT>
