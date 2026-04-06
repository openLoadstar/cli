<WAYPOINT>
## [ADDRESS] W://root/cli/test_element
## [STATUS] S_IDL

### IDENTITY
- SUMMARY: create/edit/delete 명령어 통합 테스트. 요소 생성·수정·삭제의 정상 흐름 및 엣지케이스(ID 중복, 잘못된 TYPE, 부모 없음 등)를 검증한다.
- METADATA: [Ver: 1.0, Created: 2026-03-10, Priority: HIGH]

### CONNECTIONS
- PARENT: M://root/cli
- CHILDREN: []
- REFERENCE: []
- BLACKBOX: B://root/cli/test_element

### TODO
- [x] `TestCreate_ValidTypes`: M/W/L/S 타입으로 요소 생성 후 파일 존재 및 내용 확인
- [x] `TestCreate_InvalidType`: H/B 타입 입력 시 에러 반환 확인
- [x] `TestCreate_DuplicateID`: 동일 주소 재생성 시 에러 반환 확인
- [x] `TestCreate_ParentContains`: 생성 후 부모 파일 CONTAINS.ITEMS에 주소 등록 확인
- [x] `TestCreate_NoParent`: --parent 미지정 시 에러 반환 확인
- [x] `TestEdit_ShadowHistory`: edit 실행 후 HISTORY/ 에 스냅샷 생성 확인
- [x] `TestDelete_HistoryBackup`: delete 실행 후 HISTORY/ 에 _deleted 스냅샷 생성 확인
- [x] `TestDelete_ParentContainsRemoved`: delete 후 부모 CONTAINS.ITEMS에서 주소 제거 확인
- [x] `TestAddress_ToFilePath`: dot-separated 파일명 변환 정확성 확인

### ISSUE
(없음)

### COMMENT
(없음)
</WAYPOINT>
