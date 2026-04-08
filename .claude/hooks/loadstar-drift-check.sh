#!/bin/bash
# loadstar-drift-check.sh
# PostToolUse hook: 소스코드 수정 시 LOADSTAR 메타데이터 갱신 리마인더
#
# Claude Code의 Edit/Write 도구 사용 후 실행됨.
# .loadstar/ 내부 파일 수정은 무시하고, 소스코드 수정 시에만 리마인더를 출력한다.

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# 파일 경로가 없으면 종료
if [[ -z "$FILE_PATH" ]]; then
  exit 0
fi

# .loadstar/ 메타데이터 파일은 무시
if [[ "$FILE_PATH" == *".loadstar"* ]]; then
  exit 0
fi

# .claude/ 설정 파일은 무시
if [[ "$FILE_PATH" == *".claude"* ]]; then
  exit 0
fi

# 테스트 파일, 설정 파일 등 메타 갱신 불필요한 파일 무시
BASENAME=$(basename "$FILE_PATH")
case "$BASENAME" in
  go.mod|go.sum|*.json|*.yaml|*.yml|*.toml|*.md|*.txt|LICENSE|.gitignore)
    exit 0
    ;;
esac

# 소스코드 수정 감지 → 리마인더 출력
echo "[LOADSTAR] 소스 파일 수정됨: $FILE_PATH"
echo "[LOADSTAR] 관련 WayPoint의 TODO 체크박스 갱신이 필요할 수 있습니다."

exit 0
