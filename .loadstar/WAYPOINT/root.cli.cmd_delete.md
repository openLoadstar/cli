<WAYPOINT>
## [ADDRESS] W://root/cli/cmd_delete
## [STATUS] S_STB

### IDENTITY
- SUMMARY: `loadstar delete [ADDRESS]` 구현. 파일 삭제 및 부모 WAYPOINTS/CHILDREN에서 제거. 이력은 git이 관리.
- METADATA: [Ver: 1.1, Created: 2026-03-04, Priority: HIGH]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/cmd_delete

### TODO
- [x] ADDRESS 파싱 및 파일 존재 확인
- [x] 대상 파일 삭제 (`os.Remove`)
- [x] 파일 내 CONNECTIONS.PARENT 파싱하여 부모 주소 자동 추출
- [x] 부모 요소의 WAYPOINTS(Map) 또는 CHILDREN(WayPoint)에서 해당 주소 제거
- [x] 확인 프롬프트 구현 (--force 플래그로 우회)

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
