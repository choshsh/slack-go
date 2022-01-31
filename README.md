# slack-go

AWS CodeBuild에서 Slack Incoming Webhook을 사용하여 알림을 보내기 위한 cli 도구입니다.

## 설치

### Dockerfile

Dockerfile에 바이너리 파일을 다운받아 설치합니다. (AWS CodeBuild에 사용하는 컨테이너 이미지)

```docker
...
RUN wget -q https://github.com/choshsh/slack-go/releases/latest/download/slack-go-linux-amd64 -O /usr/local/bin/slack-go && \
    chmod +x /usr/local/bin/slack-go
```

## 사용법

### 파라미터 설정

| Key | Data Type | Description | Example |
| --- | --- | --- | --- |
| msg | String | Slack 메시지 | -msg '빌드를 시작합니다’ |
| notifyStatus | Boolean | (Optional) 빌드 결과를 알림에 포함할지 여부. Default false | -notifyStatus true |

### buildspec.yml

aws codebuild 스크립트에 slack-go 명령어를 추가합니다.

```yaml
phases:
  pre_build:
    commands: slack-go -msg '빌드를 시작합니다'
	...
  post_build:
    commands: slack-go -msg '빌드가 종료됐습니다' -notifyStatus true
```

### 실행 예시

![https://user-images.githubusercontent.com/40452325/151831388-2272393e-4875-4f55-8db7-0fe40e22a996.png](https://user-images.githubusercontent.com/40452325/151831388-2272393e-4875-4f55-8db7-0fe40e22a996.png)

![https://user-images.githubusercontent.com/40452325/151831488-3de8a0b5-a8cf-40d6-b6ae-99e3949d5036.png](https://user-images.githubusercontent.com/40452325/151831488-3de8a0b5-a8cf-40d6-b6ae-99e3949d5036.png)

![https://user-images.githubusercontent.com/40452325/151831562-9108ad10-0933-40aa-a770-2dc5a8efcbdb.png](https://user-images.githubusercontent.com/40452325/151831562-9108ad10-0933-40aa-a770-2dc5a8efcbdb.png)