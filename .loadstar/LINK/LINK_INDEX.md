# LINK INDEX

링크 관계 색인 파일. 검색 및 전체 그래프 파악용.

| ADDRESS | SOURCE | TARGET | TYPE | SUMMARY |
| :--- | :--- | :--- | :--- | :--- |
| L://root/cli/create_to_edit | W://root/cli/cmd_create | W://root/cli/cmd_edit | L_SEQ | create 완료 후 edit을 구현한다 |
| L://root/cli/edit_to_delete | W://root/cli/cmd_edit | W://root/cli/cmd_delete | L_REF | Shadow History 스냅샷 패턴을 delete에서 참조 |
| L://root/cli/edit_to_history | W://root/cli/cmd_edit | W://root/cli/cmd_history | L_REF | edit이 생성하는 Shadow History를 history가 조회 |
| L://root/cli/history_to_diff | W://root/cli/cmd_history | W://root/cli/cmd_diff | L_SEQ | history 목록 조회 후 H_ID를 diff에 전달 |
| L://root/cli/history_to_rollback | W://root/cli/cmd_history | W://root/cli/cmd_rollback | L_SEQ | history 목록 조회 후 H_ID를 rollback에 전달 |
| L://root/cli/diff_to_rollback | W://root/cli/cmd_diff | W://root/cli/cmd_rollback | L_SEQ | diff로 차이를 확인한 뒤 rollback을 결정 |
| L://root/cli/delete_to_rollback | W://root/cli/cmd_delete | W://root/cli/cmd_rollback | L_REF | 삭제된 요소 복원 경로로 rollback을 참조 |
| L://root/cli/link_to_show | W://root/cli/cmd_link | W://root/cli/cmd_show | L_REF | link로 생성된 CONNECTIONS.LINKS를 show가 출력 |
| L://root/cli/infra_to_checkpoint | W://root/cli/internal_infra | W://root/cli/cmd_checkpoint | L_REF | git.Client 인프라를 checkpoint가 직접 사용 |
| L://root/cli/infra_to_todo | W://root/cli/internal_infra | W://root/cli/cmd_todo | L_REF | storage/fs 인프라를 todo가 GLOBAL_TODO_LIST 파싱에 사용 |
| L://root/cli/create_to_test_element | W://root/cli/cmd_create | W://root/cli/test_element | L_TST | create 구현을 test_element가 검증 |
| L://root/cli/edit_to_test_element | W://root/cli/cmd_edit | W://root/cli/test_element | L_TST | edit 구현을 test_element가 검증 |
| L://root/cli/delete_to_test_element | W://root/cli/cmd_delete | W://root/cli/test_element | L_TST | delete 구현을 test_element가 검증 |
| L://root/cli/checkpoint_to_test_checkpoint | W://root/cli/cmd_checkpoint | W://root/cli/test_checkpoint | L_TST | checkpoint/history/diff/rollback을 test_checkpoint가 검증 |
| L://root/cli/nav_to_test_nav | W://root/cli/cmd_link | W://root/cli/test_nav | L_TST | link/show 구현을 test_nav가 검증 |
| L://root/cli/todo_to_test_todo | W://root/cli/cmd_todo | W://root/cli/test_todo | L_TST | todo 구현을 test_todo가 검증 |
