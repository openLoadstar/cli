<WAYPOINT>
## [ADDRESS] W://root/cli/test_nav
## [STATUS] S_STB

### IDENTITY
- SUMMARY: show 명령어 테스트. WayPoint 목록 출력 및 키워드 필터링을 검증한다.
- METADATA: [Ver: 1.0, Created: 2026-03-10, Priority: MEDIUM]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []

### TODO
#### v1.x (legacy — 폐지된 명령 테스트)
- [x] `TestLink_FileCreated`: link 실행 후 LINK/ 폴더에 md 파일 생성 확인
- [x] `TestLink_BidirectionalRegistration`: SOURCE와 TARGET 양쪽 CONNECTIONS.LINKS에 등록 확인
- [x] `TestLink_InvalidType`: L_REF/L_SEQ/L_TST 외 타입 입력 시 에러 확인
- [x] `TestLink_Duplicate`: 동일 조합 재등록 시 에러 확인
- [x] `TestShow_Depth0`: depth 0에서 단일 요소 내용 출력 확인
- [x] `TestShow_DepthN`: depth N에서 하위 요소 트리 재귀 출력 확인
- [x] `TestShow_CircularRef`: 순환 참조 시 무한 루프 없이 종료 확인
#### v2.0 (현행 show 명령어)
- [x] 2026-04-10 `TestShow_ListAll`: 인자 없이 실행 시 전체 WP 목록 출력 확인 (TestListWaypoints_All로 구현됨)
- [x] 2026-04-10 `TestShow_FilterKeyword`: FILTER 인자로 주소 키워드 필터링 확인 (TestListWaypoints_Filter로 구현됨)
- [x] 2026-04-10 `TestShow_EmptyDir`: WAYPOINT 디렉토리가 비어있을 때 빈 테이블 출력 확인 (TestListWaypoints_Empty로 구현됨)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
