<WAYPOINT>
## [ADDRESS] W://root/maintenance/add_question_cmd
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar question` 명령 신규 구현 — 전체 WayPoint를 스캔하여 OPEN/DEFERRED 상태의 OPEN_QUESTIONS 목록을 표시
- METADATA: [Priority: P2, Created: 2026-04-24]
- SYNCED_AT: 2026-04-24

### CONNECTIONS
- PARENT: M://root/maintenance
- CHILDREN: []
- REFERENCE: []

### CODE_MAP
- scope:
  - cmd/

### TODO
- ADDRESS: W://root/maintenance/add_question_cmd
- SUMMARY: loadstar question 명령 구현 — 질의/응답 인터페이스 CLI 진입점
- TECH_SPEC:
  # TASK
  - [x] 2026-04-28 cmd/question.go 신규 작성
  - [x] 2026-04-28 WAYPOINT/*.md 파일 전체 스캔 로직
  - [x] 2026-04-28 OPEN_QUESTIONS 라인 regex 추출 (`[Q(\d+)(?:\s+(DEFERRED|RESOLVED\s+[\w.-]+))?\]`)
  - [x] 2026-04-28 기본 출력: ADDRESS, QID, STATE, QUESTION 테이블
  - [x] 2026-04-28 FILTER 인자 지원 (주소/키워드 부분일치)
  - [x] 2026-04-28 `loadstar question stats` 서브옵션 — OPEN/DEFERRED/RESOLVED 집계
  - [x] 2026-04-28 cmd/root.go에 AddCommand 등록
  - [x] 2026-04-28 단위 테스트 작성
  # TASK — v2 (2026-04-28)
  - [x] 2026-04-28 qRe regex에 DONE 상태 추가 — [Q1 DONE file.md] 파싱
  - [x] 2026-04-28 STATE_COLOR/표시 테이블에 DONE 추가
  - [x] 2026-04-28 `loadstar question` 기본 출력: OPEN + DEFERRED (RESOLVED/DONE 제외)
  - [x] 2026-04-28 `--with-resolved` 플래그: RESOLVED + DONE도 포함 출력
  - [x] 2026-04-28 `loadstar question done <wpAddr> <qid>` — RESOLVED → DONE 전환
  - [x] 2026-04-28 `loadstar question close <wpAddr> <qid> [사유]` — 결정 파일 없이 직접 DONE 처리
  - [x] 2026-04-28 단위 테스트: DONE 상태 파싱, done/close 명령 검증 (8개 PASS)

  # RECURRING
  - (R) 변경 후 `go build -o bin/loadstar.exe .` 검증
  - (R) 변경 후 `go test ./...` 실행

### ISSUE
- OPEN_QUESTIONS:
  - [Q1] `loadstar question --with-resolved` 플래그로 RESOLVED 질문도 표시할 수 있게 할지? (Git log 대체 가치)
  - [Q2] Decision 파일 참조 무결성 검사(DECISIONS/<id>.md 존재 여부)를 `validate`에 추가할지, 이 명령에 추가할지?

### COMMENT
- Decision 파일 생성 스캐폴드 명령(`loadstar decision add`)은 **의도적으로 보류**. 피드백 기반으로 v2에서 결정. 지금은 UI 또는 수동 작성.
- 포맷 규약: 02.SCHEMA_DEF §4, 05.ELEMENT_FORMAT Decision 섹션 참조.
</WAYPOINT>
