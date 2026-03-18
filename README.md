# avcs-cli

AVCS CLI — Go 기반 프로젝트 메타데이터 관리 도구

## 프로젝트 구조

```
avcs-cli/
├── main.go                  # 진입점
├── go.mod
├── cmd/                     # CLI 명령어 (cobra)
│   ├── root.go              # 루트 명령어 및 서브커맨드 등록
│   ├── element.go           # create / edit / delete
│   ├── checkpoint.go        # checkpoint / history / diff / rollback
│   ├── nav.go               # link / show
│   └── todo.go              # todo add / done / list
├── internal/
│   ├── core/
│   │   └── element.go       # 요소 생명주기 비즈니스 로직
│   ├── address/
│   │   └── address.go       # AVCS URI 주소 파싱 (W://root/dev/auth)
│   ├── storage/
│   │   └── fs.go            # .avcs/ 파일 시스템 I/O
│   └── git/
│       └── client.go        # go-git 래핑 (commit, hash)
└── build/                   # 크로스 컴파일 결과물
```

## 의존성

- [cobra](https://github.com/spf13/cobra) — CLI 명령어 프레임워크
- [go-git](https://github.com/go-git/go-git) — 순수 Go Git 라이브러리

## 빌드

```bash
# 현재 플랫폼
go build -o avcs .

# 크로스 컴파일
GOOS=windows GOARCH=amd64 go build -o build/avcs.exe .
GOOS=darwin  GOARCH=amd64 go build -o build/avcs-mac .
GOOS=linux   GOARCH=amd64 go build -o build/avcs-linux .
```

## 사용 예시

```bash
avcs create M root
avcs create W auth --parent M://root
avcs todo add W://root/auth M://root "인증 모듈 구현"
avcs checkpoint -m "auth 모듈 초기 구조 설계"
avcs show W://root/auth
```
