<WAYPOINT>
## [ADDRESS] W://root/cli/test_nav
## [STATUS] S_IDL

### IDENTITY
- SUMMARY: link/show 명령어 통합 테스트. Link 파일 생성·양방향 등록 및 show의 depth 재귀 출력을 검증한다.
- METADATA: [Ver: 1.0, Created: 2026-03-10, Priority: MEDIUM]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/test_nav

### TODO
- [x] `TestLink_FileCreated`: link 실행 후 LINK/ 폴더에 md 파일 생성 확인
- [x] `TestLink_BidirectionalRegistration`: SOURCE와 TARGET 양쪽 CONNECTIONS.LINKS에 등록 확인
- [x] `TestLink_InvalidType`: L_REF/L_SEQ/L_TST 외 타입 입력 시 에러 확인
- [x] `TestLink_Duplicate`: 동일 조합 재등록 시 에러 확인
- [x] `TestShow_Depth0`: depth 0에서 단일 요소 내용 출력 확인
- [x] `TestShow_DepthN`: depth N에서 하위 요소 트리 재귀 출력 확인
- [x] `TestShow_CircularRef`: 순환 참조 시 무한 루프 없이 종료 확인

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
