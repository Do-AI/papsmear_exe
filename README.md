# papsmear agent
golang을 사용한 Agent 프로그램

### godoc
```bash
godoc -http=localhost:8080
```
위의 명령어를 terminal에서 실행 한 후 아래 url로 internal package 문서를 볼 수 있다.  

[godoc internal package](http://localhost:8080/pkg/github.com/p829911/agent/papsmear/internal/)

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