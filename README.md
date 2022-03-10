# papsmear agent
golang을 사용한 Agent 프로그램  

### Setting
config/config.yaml 파일의 credential 설치할 의료기관 정보에 맞춰 바꾼다.  

### Build
프로젝트 루트 경로에서 `go build cmd/main.go` 명령어를 실행하여 `exe` 파일을 만든다. 
이때 폴더가 위치하는 경로는 `$gopath/src/github.com/doai/papsmear` (`user_home_directory/go/src/github.com/doai/papsmear`) 이어야 한다.

### install env path
```powershell
./install.bat
```

- powershell 실행 권장
- opencv 의존성 라이브러리(.dll)를 환경변수로 등록하여 프로그램에서 불러올 수 있도록 세팅해주는 batch 파일이다.

### background process in window
```powershell
./start.vbs
```

- powershell 실행 권장
- agent 프로그램을 background process, no gui process로 실행시키기 위해서 작성한 script이다.
- window에서 powershell을 통해서 실행해야 정상적으로 백그라운드 프로세스에서 작동한다.

- 실행 파일을 못찾는다면 `run.bat`에서 경로를 확인할 것
### Deploy
papsmear_exe 레파지토리의 bin(opencv dll 파일들)을 `exe` 파일과 함께 압축하여 해당 의료기관 컴퓨터에 압축해제 후 `exe` 파일 실행하면 된다. 
이때 bin path를 해당 컴퓨터 system path에 등록해줘야 하는데 이 부분은 추가로 cmd script를 작성하여야 한다.  
