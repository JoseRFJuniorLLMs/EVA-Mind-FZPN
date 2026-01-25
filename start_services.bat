@echo off
REM EVA-Mind-FZPN - Start All Services
REM SPRINT 7 - Integration Layer

echo ========================================
echo EVA-Mind Integration Services Startup
echo ========================================
echo.

echo [1/4] Instalando dependencias Python...
python -m pip install -r requirements.txt
if errorlevel 1 (
    echo ERRO: Falha ao instalar dependencias Python
    pause
    exit /b 1
)
echo.

echo [2/4] Compilando Go Integration Service...
cd cmd\integration_service
go build -o ..\..\eva_integration_service.exe main.go
if errorlevel 1 (
    echo ERRO: Falha ao compilar Go service
    cd ..\..
    pause
    exit /b 1
)
cd ..\..
echo.

echo [3/4] Iniciando Go Integration Service (porta 8081)...
start "EVA Go Service" cmd /k eva_integration_service.exe
timeout /t 3 > nul
echo.

echo [4/4] Iniciando Python API Server (porta 8000)...
start "EVA API Server" cmd /k python api_server.py
timeout /t 3 > nul
echo.

echo ========================================
echo SERVICOS INICIADOS COM SUCESSO!
echo ========================================
echo.
echo Go Integration Service: http://localhost:8081
echo Python API Server:      http://localhost:8000
echo API Documentation:      http://localhost:8000/docs
echo.
echo Pressione qualquer tecla para abrir o navegador...
pause > nul

start http://localhost:8000/docs

echo.
echo Para parar os servicos, feche as janelas cmd abertas.
echo.
pause
