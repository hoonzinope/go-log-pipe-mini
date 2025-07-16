# go-log-pipe-mini

## 개요

`go-log-pipe-mini`는 로그 파일을 실시간으로 감시하고, 지정된 조건에 따라 로그를 필터링하여 출력하는 간단한 로그 처리 파이프라인입니다. 설정 파일을 통해 입력 소스, 필터링 규칙 및 출력 대상을 쉽게 구성할 수 있습니다.

## 주요 기능

-   **폴더 단위 파일 감시**: 지정된 폴더의 모든 로그 파일을 실시간으로 감시하여 새로운 로그 라인을 읽어들입니다.
-   **로그 필터링**: `config.yml`에 정의된 규칙에 따라 로그를 필터링합니다.
    -   `grep`과 유사한 패턴 매칭 필터를 지원합니다.
    -   여러 필터 조건을 `AND` 또는 `OR` 논리로 결합할 수 있습니다.
    -   패턴 매칭 시 대소문자 구분을 무시하는 옵션을 제공합니다.
-   **다중 파일 오프셋 관리**: 처리된 각 로그 파일의 마지막 위치(오프셋)를 `offset.state` 파일에 저장하여, 재시작 시 중복 처리를 방지합니다.
-   **유연한 설정**: `config.yml` 파일을 통해 입력, 필터, 출력 동작을 쉽게 변경할 수 있습니다.

## 프로젝트 구조

```
/
├───.gitignore
├───config.yml         # 애플리케이션 설정 파일
├───go.mod
├───go.sum
├───main.go            # 애플리케이션 진입점
├───README.md
├───confmanager/       # 설정 관리 패키지
│   └───config.go
├───filter/            # 로그 필터링 패키지
│   └───filter.go
├───generate/          # 테스트용 로그 생성기
│   └───genLog.go
├───input/             # 로그 입력 처리 패키지
│   └───input_manager.go
├───offset/            # 오프셋 관리 패키지
│   ├───offset_manager.go
│   └───offsetData.go
└───output/            # 결과 출력 패키지
    └───out.go
```

## 설정 (`config.yml`)

애플리케이션의 동작은 `config.yml` 파일을 통해 제어됩니다.

```yaml
INPUT:
  NAME: testlog
  TYPE: file
  PATH: ./logs/

FILTER:
  MODE: OR
  FILTERS:
    - TYPE: grep
      OPTIONS:
        IGNORE_CASE: true
        PATTERN: "ERROR"
    - TYPE: grep
      OPTIONS:
        IGNORE_CASE: true
        PATTERN: "WARN"

OUTPUT:
  TYPE: stdout
```

-   **INPUT**: 입력 소스를 정의합니다.
    -   `NAME`: 로그에 추가될 접두사입니다.
    -   `TYPE`: 현재 `file`만 지원됩니다.
    -   `PATH`: 감시할 로그 파일이 있는 디렉토리의 경로를 지정합니다.
-   **FILTER**: 필터링 규칙을 정의합니다.
    -   `MODE`: `AND` 또는 `OR`를 지정하여 여러 필터의 논리적 관계를 설정합니다.
    -   `FILTERS`: 적용할 필터 목록입니다.
        -   `TYPE`: 현재 `grep`만 지원됩니다.
        -   `OPTIONS`:
            -   `IGNORE_CASE`: `true`로 설정하면 대소문자를 구분하지 않고 패턴을 비교합니다.
            -   `PATTERN`: 필터링할 키워드를 지정합니다.
-   **OUTPUT**: 출력 대상을 정의합니다.
    -   `TYPE`: 현재 `stdout`(표준 출력)만 지원됩니다.

## 사용법

1.  **저장소 복제**:
    ```bash
    git clone https://github.com/your-username/go-log-pipe-mini.git
    cd go-log-pipe-mini
    ```

2.  **설정 파일 수정**:
    `config.yml` 파일을 열어 `INPUT.PATH`에 감시할 로그 파일이 있는 디렉토리 경로를 지정하고, `FILTER` 규칙을 필요에 맞게 수정합니다.

3.  **애플리케이션 실행**:
    ```bash
    go run main.go
    ```

4.  **테스트 로그 생성 (선택 사항)**:
    `generate` 패키지는 테스트를 위해 무작위 로그를 생성하는 기능을 포함하고 있습니다. `main.go`에서 `generate.GenLogWithFolder` 고루틴을 활성화하여 테스트 로그를 생성하고 파이프라인을 테스트할 수 있습니다.