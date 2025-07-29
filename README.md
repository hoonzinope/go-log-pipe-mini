[![en](https://img.shields.io/badge/lang-en-red.svg)](README.en.md)
[![ko](https://img.shields.io/badge/lang-ko-blue.svg)](README.md)

# go-log-pipe-mini

## 개요

`go-log-pipe-mini`는 로그 파일을 실시간으로 감시하고, 지정된 조건에 따라 로그를 필터링하여 다양한 대상으로 출력하는 고성능 로그 처리 파이프라인입니다. 유연한 설정 파일을 통해 다중 입력, 다중 출력, JSON 파싱, 로그 로테이션 등 고급 기능을 쉽게 구성할 수 있습니다.

## 주요 기능

-   **다중 입력 및 출력**: 여러 소스(`INPUTS`)로부터 로그를 수집하여 여러 대상(`OUTPUTS`)으로 동시에 전송할 수 있습니다.
-   **실시간 파일 감시**: 지정된 폴더의 모든 로그 파일을 실시간으로 감시하여 새로운 로그를 즉시 처리합니다.
-   **고급 로그 필터링**:
    -   일반 텍스트(`grep`) 및 JSON 형식(`json_grep`) 로그 필터링을 지원합니다.
    -   여러 필터 조건을 `AND` 또는 `OR` 논리로 결합할 수 있습니다.
    -   대소문자 구분 없는 패턴 매칭을 지원합니다.
-   **다양한 출력 대상**:
    -   **콘솔 (`stdout`)**: 표준 출력으로 로그를 표시합니다.
    -   **파일 (`file`)**: 로그를 파일로 저장하며, `daily`, `hourly` 등 시간 기반 또는 크기 기반의 로그 로테이션(`Rolling`) 기능을 지원합니다.
    -   **HTTP/HTTPS (`http`)**: 지정된 웹훅(Webhook) URL로 로그를 `POST` 방식으로 전송합니다.
-   **오프셋 관리**: 각 파일의 처리된 위치를 `offset.state` 파일에 저장하여, 재시작 시 중복 처리를 방지하고 데이터 정합성을 보장합니다.
-   **배치 처리**: 여러 로그를 모아 한 번에 전송하는 `BATCH_SIZE` 및 `FLUSH_INTERVAL` 옵션을 통해 네트워크 및 시스템 부하를 최소화합니다.
-   **내장 상태 서버**: `server` 모듈을 통해 애플리케이션의 상태(`healthCheck`) 및 성능 메트릭(`metrics`)을 HTTP 엔드포인트로 제공합니다.

## 프로젝트 구조

```
/
├───.gitignore
├───config.yml         # 애플리케이션 설정 파일
├───go.mod
├───go.sum
├───main.go            # 애플리케이션 진입점
├───README.md
├───confmanager/       # 설정 관리
│   └───config.go
├───filter/            # 로그 필터링
│   ├───filter.go
│   └───manager.go
├───generate/          # 테스트용 로그 생성기
│   └───genLog.go
├───input/             # 로그 입력 처리
│   ├───manager.go
│   ├───node_file.go
│   ├───node_input.go
│   └───parse/
│       └───json_parser.go
├───offset/            # 오프셋 관리
│   └───offset_manager.go
├───output/            # 결과 출력 처리
│   ├───console.go
│   ├───file.go
│   ├───httppost.go
│   └───manager.go
├───server/            # 내장 HTTP 서버
│   ├───healthCheck.go
│   ├───logReciever.go
│   ├───metrics.go
│   └───runServer.go
└───shared/            # 공용 데이터 및 함수
    ├───data.go
    └───stat.go
```

## 설정 (`config.yml`)

애플리케이션의 모든 동작은 `config.yml`을 통해 제어됩니다.

```yaml
INPUTS:
  - NAME: syslog      # 입력 소스 식별자
    TYPE: file         # 입력 타입 (현재 'file'만 지원)
    PATH: ./logs/      # 감시할 로그 폴더 경로
    PARSER: text       # 파서 타입 ('text' 또는 'json')
  - NAME: applog
    TYPE: file
    PATH: ./json_logs/
    PARSER: json

FILTERS:
  syslog:              # 'syslog' 입력에 적용할 필터
    MODE: OR           # 필터 논리 (AND 또는 OR)
    RULES:
      - TYPE: grep     # 일반 텍스트 필터
        OPTIONS:
          IGNORE_CASE: true
          PATTERN: "error"
  applog:              # 'applog' 입력에 적용할 필터
    MODE: AND
    RULES:
      - TYPE: json_grep # JSON 필드 필터
        OPTIONS:
          FIELD: "level"
          IGNORE_CASE: true
          PATTERN: "error"

OUTPUTS:
  - TYPE: stdout       # 콘솔 출력
    TARGETS: [syslog, applog] # 이 출력을 사용할 입력 소스
    OPTIONS:
      BATCH_SIZE: 5
      FLUSH_INTERVAL: 2s
  - TYPE: file         # 파일 출력
    TARGETS: [syslog]
    OPTIONS:
      PATH: ./output_logs/
      FILENAME: syslog.log
      ROLLING: daily    # 로그 로테이션 (daily, hourly, monthly)
      MAX_SIZE: 100MB   # 최대 파일 크기
      MAX_FILES: 7      # 최대 보존 파일 수
      BATCH_SIZE : 10
      FLUSH_INTERVAL: 5s
  - TYPE: http         # HTTP 출력
    TARGETS: [syslog, applog]
    OPTIONS:
      URL: http://localhost:8080/logs
      METHOD: POST
      HEADERS:
        Content-Type: application/json
      TIMEOUT: 5s
      BATCH_SIZE: 10
      FLUSH_INTERVAL: 5s
```

## 사용법

1.  **저장소 복제**:
    ```bash
    git clone https://github.com/your-username/go-log-pipe-mini.git
    cd go-log-pipe-mini
    ```

2.  **의존성 설치**:
    ```bash
    go mod tidy
    ```

3.  **설정 파일 수정**:
    `config.yml` 파일을 열어 `INPUTS`, `FILTERS`, `OUTPUTS`를 필요에 맞게 수정합니다.

4.  **애플리케이션 실행**:
    ```bash
    go run main.go
    ```

5.  **상태 확인 (선택 사항)**:
    애플리케이션 실행 중 다음 엔드포인트를 통해 상태를 확인할 수 있습니다.
    -   **Health Check**: `http://localhost:8080/health`
    -   **Metrics**: `http://localhost:8080/metrics`

6.  **테스트 로그 생성 (선택 사항)**:
    `generate` 패키지는 테스트용 로그 생성을 지원합니다. `main.go`에서 `generate.GenLogWithFolder` 및 `generate.GenerateJsonLog` 고루틴을 활성화하여 파이프라인을 테스트할 수 있습니다.
