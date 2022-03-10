@ECHO OFF
:: This is batch file for Papsmear Agent Install Setting
:: - Set Env variable for OPENCV .dll files doai\papsmear\bin
:: - Excute Main Program for daemon process

Echo Please wait.. Checking Environment Information

:: Set Env variable

Echo %cd%
SET ENV_PATH=%cd%\bin\opencv

ECHO Add Environment opencv folder path 
ECHO %ENV_PATH%

setx PATH "%PATH%;%ENV_PATH%" -m