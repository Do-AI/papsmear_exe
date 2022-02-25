# papsmear agent
golang을 사용한 Agent 프로그램  

# 사용 방법
config/config.yaml 파일의 credential 설치할 의료기관 정보에 맞춰 바꾼다.  

# 빌드 방법
프로젝트 루트 경로에서 `go build cmd/main.go` 명령어를 실행하여 `exe` 파일을 만든다. 
이때 폴더가 위치하는 경로는 `$gopath/src/github.com/doai/papsmear` (`user_home_directory/go/src/github.com/doai/papsmear`) 이어야 한다.

# 배포 방법
papsmear_exe 레파지토리의 bin(opencv dll 파일들)을 `exe` 파일과 함께 압축하여 해당 의료기관 컴퓨터에 압축해제 후 `exe` 파일 실행하면 된다. 
이때 bin path를 해당 컴퓨터 system path에 등록해줘야 하는데 이 부분은 추가로 cmd script를 작성하여야 한다.  
